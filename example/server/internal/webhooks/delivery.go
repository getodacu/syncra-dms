package webhooks

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	signatureHeader = "X-Syncra-Signature"
	timestampHeader = "X-Syncra-Timestamp"

	defaultDeliveryTimeout = 10 * time.Second
)

// JobPayload is the webhook representation of a processed job.
type JobPayload struct {
	ID               string  `json:"id"`
	Status           string  `json:"status,omitempty"`
	OriginalFilename string  `json:"original_filename,omitempty"`
	ErrorMessage     *string `json:"error_message,omitempty"`
}

// JobEventInput describes a user webhook event ready for delivery.
type JobEventInput struct {
	Event  Event
	UserID *string
	Job    JobPayload
}

// DispatcherConfig configures webhook delivery.
type DispatcherConfig struct {
	DB         *gorm.DB
	PrivateKey string
	// HTTPClient bypasses the secure default transport and redirect policy.
	// Callers that set it must provide equivalent protections.
	HTTPClient *http.Client
	Timeout    time.Duration

	Now          func() time.Time
	LookupIPAddr func(context.Context, string) ([]net.IPAddr, error)
	DialContext  func(context.Context, string, string) (net.Conn, error)
}

// Dispatcher delivers active user webhook events.
type Dispatcher struct {
	db         *gorm.DB
	privateKey string
	httpClient *http.Client
	timeout    time.Duration
	now        func() time.Time
}

type deliveryPayload struct {
	Event Event        `json:"event"`
	Data  deliveryData `json:"data"`
}

type deliveryData struct {
	Job JobPayload `json:"job"`
}

// NewDispatcher creates a webhook dispatcher. If HTTPClient is omitted, the
// dispatcher uses a transport that validates the resolved address before dialing.
func NewDispatcher(config DispatcherConfig) *Dispatcher {
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultDeliveryTimeout
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	client := config.HTTPClient
	if client == nil {
		client = newSecureHTTPClient(timeout, config.LookupIPAddr, config.DialContext)
	}
	return &Dispatcher{
		db:         config.DB,
		privateKey: config.PrivateKey,
		httpClient: client,
		timeout:    timeout,
		now:        now,
	}
}

// Dispatch sends the webhook event if the user has an active webhook for it.
func (d *Dispatcher) Dispatch(ctx context.Context, input JobEventInput) error {
	if input.UserID == nil {
		return nil
	}
	if d == nil {
		return errors.New("webhook dispatcher is nil")
	}
	if d.db == nil {
		return errors.New("webhook dispatcher database is required")
	}
	if d.httpClient == nil {
		return errors.New("webhook dispatcher HTTP client is required")
	}

	var hook Webhook
	err := d.db.WithContext(ctx).Where("user_id = ?", *input.UserID).First(&hook).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("load webhook: %w", err)
	}
	if !slices.Contains(DecodeEvents(hook.EventsActive), input.Event) {
		return nil
	}

	secret, err := DecryptSecret(d.privateKey, hook.SecretKey)
	if err != nil {
		return fmt.Errorf("decrypt webhook secret: %w", err)
	}

	payload := deliveryPayload{
		Event: input.Event,
		Data: deliveryData{
			Job: webhookJobPayload(input.Event, input.Job),
		},
	}
	rawBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}

	return d.deliver(ctx, hook.URL, secret, rawBody)
}

func (d *Dispatcher) DispatchJobEvent(ctx context.Context, input JobEventInput) error {
	return d.Dispatch(ctx, input)
}

// SignPayload signs raw webhook body bytes with the Syncra v1 HMAC contract.
func SignPayload(secret string, timestamp string, rawBody []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp))
	_, _ = mac.Write([]byte("."))
	_, _ = mac.Write(rawBody)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}

func (d *Dispatcher) deliver(ctx context.Context, targetURL string, secret string, rawBody []byte) error {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(rawBody))
	if err != nil {
		return fmt.Errorf("create webhook request: %w", err)
	}

	timestamp := strconv.FormatInt(d.now().UTC().Unix(), 10)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(timestampHeader, timestamp)
	req.Header.Set(signatureHeader, SignPayload(secret, timestamp, rawBody))

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("deliver webhook: request failed: %w", sanitizeRequestError(err))
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("webhook delivery returned status %d", resp.StatusCode)
	}
	return nil
}

type sanitizedURLError struct {
	op  string
	err error
}

func (e sanitizedURLError) Error() string {
	if e.err == nil {
		if e.op == "" {
			return "request failed"
		}
		return e.op
	}
	if e.op == "" {
		return e.err.Error()
	}
	return e.op + ": " + e.err.Error()
}

func (e sanitizedURLError) Unwrap() error {
	return e.err
}

func sanitizeRequestError(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return sanitizedURLError{
			op:  urlErr.Op,
			err: sanitizeRequestError(urlErr.Err),
		}
	}
	return err
}

func webhookJobPayload(event Event, job JobPayload) JobPayload {
	if event == EventJobSucceeded {
		return JobPayload{ID: job.ID}
	}
	return job
}

func newSecureHTTPClient(timeout time.Duration, lookupIPAddr func(context.Context, string) ([]net.IPAddr, error), dialContext func(context.Context, string, string) (net.Conn, error)) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil

	dialer := secureDialer{
		lookupIPAddr: lookupIPAddr,
		dialContext:  dialContext,
	}
	if dialer.lookupIPAddr == nil {
		dialer.lookupIPAddr = net.DefaultResolver.LookupIPAddr
	}
	if dialer.dialContext == nil {
		netDialer := &net.Dialer{}
		dialer.dialContext = netDialer.DialContext
	}
	transport.DialContext = dialer.DialContext

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

type secureDialer struct {
	lookupIPAddr func(context.Context, string) ([]net.IPAddr, error)
	dialContext  func(context.Context, string, string) (net.Conn, error)
}

func (d secureDialer) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("parse webhook target address: %w", err)
	}

	addrs, err := d.resolve(ctx, host)
	if err != nil {
		return nil, err
	}

	var lastDialErr error
	for _, resolved := range addrs {
		addr, ok := netipAddrFromIP(resolved.IP)
		if !ok {
			lastDialErr = fmt.Errorf("invalid resolved webhook address %q", resolved.IP.String())
			continue
		}
		if err := ValidateResolvedAddress(addr); err != nil {
			return nil, err
		}
		conn, err := d.dialContext(ctx, network, net.JoinHostPort(addr.String(), port))
		if err == nil {
			return conn, nil
		}
		lastDialErr = err
	}
	if lastDialErr != nil {
		return nil, fmt.Errorf("dial webhook target: %w", lastDialErr)
	}
	return nil, errors.New("webhook target resolved no addresses")
}

func (d secureDialer) resolve(ctx context.Context, host string) ([]net.IPAddr, error) {
	hostname := strings.TrimSuffix(host, ".")
	if addr, err := netip.ParseAddr(hostname); err == nil {
		return []net.IPAddr{{IP: net.ParseIP(addr.String())}}, nil
	}
	addrs, err := d.lookupIPAddr(ctx, hostname)
	if err != nil {
		return nil, fmt.Errorf("resolve webhook target: %w", err)
	}
	if len(addrs) == 0 {
		return nil, errors.New("webhook target resolved no addresses")
	}
	return addrs, nil
}

func netipAddrFromIP(ip net.IP) (netip.Addr, bool) {
	if ip4 := ip.To4(); ip4 != nil {
		return netip.AddrFrom4([4]byte{ip4[0], ip4[1], ip4[2], ip4[3]}), true
	}
	if ip16 := ip.To16(); ip16 != nil {
		var raw [16]byte
		copy(raw[:], ip16)
		return netip.AddrFrom16(raw), true
	}
	return netip.Addr{}, false
}
