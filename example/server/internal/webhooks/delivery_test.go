package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

const webhookDeliverySecret = "webhook-delivery-secret"

var _ interface {
	DispatchJobEvent(context.Context, JobEventInput) error
} = (*Dispatcher)(nil)

func TestDispatcherSkipsInactiveEvent(t *testing.T) {
	db := newDeliveryTestDB(t)
	var hits atomic.Int32
	server, webhookURL, config := newDeliveryTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	config.DB = db

	userID := uuid.NewString()
	createDeliveryTestWebhook(t, db, userID, webhookURL, []Event{EventJobSucceeded})

	dispatcher := NewDispatcher(config)
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobFailed,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-1",
			Status:           "failed",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
	if hits.Load() != 0 {
		t.Fatalf("server hits = %d, want 0", hits.Load())
	}
}

func TestDispatcherSkipsNilUserIDAndMissingWebhook(t *testing.T) {
	db := newDeliveryTestDB(t)
	var hits atomic.Int32
	server, _, config := newDeliveryTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	config.DB = db

	dispatcher := NewDispatcher(config)
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event: EventJobStarted,
		Job: JobPayload{
			ID:               "job-0",
			Status:           "queued",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() with nil UserID error = %v", err)
	}

	userID := uuid.NewString()
	err = dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobStarted,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-0",
			Status:           "queued",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() with missing webhook error = %v", err)
	}
	if hits.Load() != 0 {
		t.Fatalf("server hits = %d, want 0", hits.Load())
	}
}

func TestDispatcherSendsSucceededEventWithSignatureAndJobIDOnly(t *testing.T) {
	db := newDeliveryTestDB(t)

	requests := make(chan capturedWebhookRequest, 1)
	server, webhookURL, config := newDeliveryTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		requests <- capturedWebhookRequest{
			Method:         r.Method,
			ContentType:    r.Header.Get("Content-Type"),
			Timestamp:      r.Header.Get("X-Syncra-Timestamp"),
			Signature:      r.Header.Get("X-Syncra-Signature"),
			Body:           body,
			RequestURIPath: r.URL.Path,
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	config.DB = db

	userID := uuid.NewString()
	createDeliveryTestWebhook(t, db, userID, webhookURL, []Event{EventJobSucceeded})

	dispatcher := NewDispatcher(config)
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobSucceeded,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-1",
			Status:           "completed",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}

	var captured capturedWebhookRequest
	select {
	case captured = <-requests:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for webhook request")
	}
	if captured.Method != http.MethodPost {
		t.Fatalf("method = %s, want POST", captured.Method)
	}
	if captured.RequestURIPath != "/webhook" {
		t.Fatalf("path = %s, want /webhook", captured.RequestURIPath)
	}
	if captured.ContentType != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", captured.ContentType)
	}
	if captured.Timestamp != "1710000000" {
		t.Fatalf("X-Syncra-Timestamp = %q", captured.Timestamp)
	}
	wantSignature := expectedWebhookSignature(webhookDeliverySecret, captured.Timestamp, captured.Body)
	if captured.Signature != wantSignature {
		t.Fatalf("X-Syncra-Signature = %q, want %q", captured.Signature, wantSignature)
	}
	if !hmac.Equal([]byte(captured.Signature), []byte(SignPayload(webhookDeliverySecret, captured.Timestamp, captured.Body))) {
		t.Fatal("SignPayload() did not reproduce request signature")
	}
	wantBody := `{"event":"job.succeeded","data":{"job":{"id":"job-1"}}}`
	if string(captured.Body) != wantBody {
		t.Fatalf("body = %s, want %s", captured.Body, wantBody)
	}

	var payload map[string]any
	if err := json.Unmarshal(captured.Body, &payload); err != nil {
		t.Fatalf("payload is not valid JSON: %v\n%s", err, captured.Body)
	}
	if payload["event"] != string(EventJobSucceeded) {
		t.Fatalf("event = %#v, want %q", payload["event"], EventJobSucceeded)
	}
	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("data = %#v, want object", payload["data"])
	}
	if _, ok := data["document"]; ok {
		t.Fatalf("document was present in success payload: %s", captured.Body)
	}
	job, ok := data["job"].(map[string]any)
	if !ok {
		t.Fatalf("job = %#v, want object", data["job"])
	}
	if job["id"] != "job-1" {
		t.Fatalf("job.id = %#v, want job-1", job["id"])
	}
	for _, forbidden := range []string{"status", "original_filename", "document_id"} {
		if _, ok := job[forbidden]; ok {
			t.Fatalf("job.%s was present in success payload: %s", forbidden, captured.Body)
		}
	}
}

func TestDispatcherIncludesFailedJobErrorMessageAndOmitsNilDocument(t *testing.T) {
	db := newDeliveryTestDB(t)
	requests := make(chan []byte, 1)
	server, webhookURL, config := newDeliveryTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		requests <- body
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	config.DB = db

	userID := uuid.NewString()
	createDeliveryTestWebhook(t, db, userID, webhookURL, []Event{EventJobFailed})

	errorMessage := "OCR failed"
	dispatcher := NewDispatcher(config)
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobFailed,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-2",
			Status:           "failed",
			OriginalFilename: "invoice.pdf",
			ErrorMessage:     &errorMessage,
		},
	})
	if err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}

	var body []byte
	select {
	case body = <-requests:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for webhook request")
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}
	data := payload["data"].(map[string]any)
	job := data["job"].(map[string]any)
	if job["error_message"] != errorMessage {
		t.Fatalf("error_message = %#v", job["error_message"])
	}
	if _, ok := data["document"]; ok {
		t.Fatalf("document was present for failed event with nil document: %s", body)
	}
}

func TestDispatcherReturnsErrorForNon2xxResponse(t *testing.T) {
	db := newDeliveryTestDB(t)
	server, webhookURL, config := newDeliveryTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	config.DB = db

	userID := uuid.NewString()
	createDeliveryTestWebhook(t, db, userID, webhookURL, []Event{EventJobStarted})

	dispatcher := NewDispatcher(config)
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobStarted,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-3",
			Status:           "queued",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err == nil {
		t.Fatal("Dispatch() error = nil, want non-2xx error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("Dispatch() error = %v, want status code", err)
	}
}

func TestDispatcherReturnsErrorForRedirectAndDoesNotFollow(t *testing.T) {
	for _, status := range []int{http.StatusFound, http.StatusTemporaryRedirect} {
		t.Run(fmt.Sprintf("status_%d", status), func(t *testing.T) {
			db := newDeliveryTestDB(t)
			var redirectedHits atomic.Int32
			redirectTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				redirectedHits.Add(1)
				w.WriteHeader(http.StatusNoContent)
			}))
			defer redirectTarget.Close()
			server, webhookURL, config := newDeliveryTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, redirectTarget.URL+"/capture", status)
			}))
			defer server.Close()
			config.DB = db

			userID := uuid.NewString()
			createDeliveryTestWebhook(t, db, userID, webhookURL, []Event{EventJobStarted})

			dispatcher := NewDispatcher(config)
			err := dispatcher.Dispatch(context.Background(), JobEventInput{
				Event:  EventJobStarted,
				UserID: &userID,
				Job: JobPayload{
					ID:               "job-redirect",
					Status:           "queued",
					OriginalFilename: "invoice.pdf",
				},
			})
			if err == nil {
				t.Fatal("Dispatch() error = nil, want redirect status error")
			}
			if !strings.Contains(err.Error(), fmt.Sprintf("%d", status)) {
				t.Fatalf("Dispatch() error = %v, want redirect status code", err)
			}
			if redirectedHits.Load() != 0 {
				t.Fatalf("redirect target hits = %d, want 0", redirectedHits.Load())
			}
		})
	}
}

func TestDispatcherSanitizesRequestErrors(t *testing.T) {
	db := newDeliveryTestDB(t)
	userID := uuid.NewString()
	webhookURL := "http://delivery-user:delivery-password@delivery-error.test:443/webhook?token=secret-token"
	createDeliveryTestWebhook(t, db, userID, webhookURL, []Event{EventJobStarted})

	dispatcher := NewDispatcher(DispatcherConfig{
		DB:         db,
		PrivateKey: testPrivateKey,
		Timeout:    time.Second,
		Now: func() time.Time {
			return time.Unix(1710000000, 0).UTC()
		},
		LookupIPAddr: func(_ context.Context, host string) ([]net.IPAddr, error) {
			if host != "delivery-error.test" {
				return nil, fmt.Errorf("unexpected host lookup %q", host)
			}
			return []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}, nil
		},
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			return nil, fmt.Errorf("connection refused")
		},
	})
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobStarted,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-request-error",
			Status:           "queued",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err == nil {
		t.Fatal("Dispatch() error = nil, want request failure")
	}
	got := err.Error()
	if !strings.Contains(got, "deliver webhook: request failed") {
		t.Fatalf("Dispatch() error = %v, want generic request failure message", err)
	}
	if !strings.Contains(got, "connection refused") {
		t.Fatalf("Dispatch() error = %v, want sanitized underlying cause", err)
	}
	for _, leaked := range []string{
		webhookURL,
		"delivery-user",
		"delivery-password",
		"secret-token",
		"token=secret-token",
	} {
		if strings.Contains(got, leaked) {
			t.Fatalf("Dispatch() error leaked %q: %v", leaked, err)
		}
	}
}

func TestDispatcherBlocksPrivateResolvedAddressAtDeliveryTime(t *testing.T) {
	db := newDeliveryTestDB(t)
	userID := uuid.NewString()
	createDeliveryTestWebhook(t, db, userID, "http://delivery-blocked.test/webhook", []Event{EventJobStarted})

	var dialed atomic.Bool
	dispatcher := NewDispatcher(DispatcherConfig{
		DB:         db,
		PrivateKey: testPrivateKey,
		Timeout:    time.Second,
		Now: func() time.Time {
			return time.Unix(1710000000, 0).UTC()
		},
		LookupIPAddr: func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil
		},
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			dialed.Store(true)
			return nil, fmt.Errorf("dial should not be reached")
		},
	})
	err := dispatcher.Dispatch(context.Background(), JobEventInput{
		Event:  EventJobStarted,
		UserID: &userID,
		Job: JobPayload{
			ID:               "job-4",
			Status:           "queued",
			OriginalFilename: "invoice.pdf",
		},
	})
	if err == nil {
		t.Fatal("Dispatch() error = nil, want blocked address error")
	}
	if !strings.Contains(err.Error(), "not allowed") {
		t.Fatalf("Dispatch() error = %v, want blocked address error", err)
	}
	if dialed.Load() {
		t.Fatal("DialContext was called for blocked private address")
	}
}

type capturedWebhookRequest struct {
	Method         string
	ContentType    string
	Timestamp      string
	Signature      string
	Body           []byte
	RequestURIPath string
}

func newDeliveryTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + uuid.NewString() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&auth.User{}, &Webhook{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func createDeliveryTestWebhook(t *testing.T, db *gorm.DB, userID string, url string, events []Event) {
	t.Helper()
	user := auth.User{
		ID:            userID,
		Name:          "Webhook User",
		Email:         strings.ReplaceAll(userID, "-", "") + "@example.com",
		EmailVerified: true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	encryptedSecret, err := EncryptSecret(testPrivateKey, webhookDeliverySecret)
	if err != nil {
		t.Fatalf("EncryptSecret() error = %v", err)
	}
	encodedEvents, err := EncodeEvents(events)
	if err != nil {
		t.Fatalf("EncodeEvents() error = %v", err)
	}
	hook := Webhook{
		UserID:       userID,
		URL:          url,
		SecretKey:    encryptedSecret,
		EventsActive: encodedEvents,
	}
	if err := db.Create(&hook).Error; err != nil {
		t.Fatalf("create webhook: %v", err)
	}
}

func newDeliveryTestServer(t *testing.T, handler http.Handler) (*httptest.Server, string, DispatcherConfig) {
	t.Helper()
	server := httptest.NewServer(handler)
	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		server.Close()
		t.Fatalf("split listener address: %v", err)
	}
	publicHost := "webhook-delivery.test"
	webhookURL := "http://" + net.JoinHostPort(publicHost, port) + "/webhook"
	config := DispatcherConfig{
		PrivateKey: testPrivateKey,
		Timeout:    time.Second,
		Now: func() time.Time {
			return time.Unix(1710000000, 0).UTC()
		},
		LookupIPAddr: func(_ context.Context, host string) ([]net.IPAddr, error) {
			if host != publicHost {
				return nil, fmt.Errorf("unexpected host lookup %q", host)
			}
			return []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}, nil
		},
		DialContext: func(ctx context.Context, network string, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, network, server.Listener.Addr().String())
		},
	}
	return server, webhookURL, config
}

func expectedWebhookSignature(secret string, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}
