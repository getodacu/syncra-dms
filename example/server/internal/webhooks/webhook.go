package webhooks

import (
	"encoding/json"
	"errors"
	"net/netip"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
)

type Event string

const (
	EventJobStarted   Event = "job.started"
	EventJobFailed    Event = "job.failed"
	EventJobSucceeded Event = "job.succeeded"
)

var supportedEvents = []Event{EventJobStarted, EventJobFailed, EventJobSucceeded}

var (
	errInvalidWebhookURL    = errors.New("url must be an absolute http or https URL")
	errUnsafeWebhookAddress = errors.New("webhook target address is not allowed")
)

var nonPublicWebhookAddressPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("255.255.255.255/32"),
	netip.MustParsePrefix("::/128"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("64:ff9b::/96"),
	netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("ff00::/8"),
}

type Webhook struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       string         `gorm:"column:user_id;type:uuid;not null;uniqueIndex" json:"user_id"`
	User         auth.User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	URL          string         `gorm:"column:url;type:text;not null" json:"url"`
	SecretKey    string         `gorm:"column:secret_key;type:text;not null" json:"-"`
	EventsActive datatypes.JSON `gorm:"column:events_active;type:jsonb;not null;default:'[]'" json:"events_active"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (w *Webhook) BeforeCreate(_ *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (Webhook) TableName() string {
	return "webhooks"
}

func ValidateURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("url is required")
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errInvalidWebhookURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errInvalidWebhookURL
	}
	hostname := strings.TrimSuffix(parsed.Hostname(), ".")
	if hostname == "" {
		return "", errInvalidWebhookURL
	}
	if strings.EqualFold(hostname, "localhost") {
		return "", errInvalidWebhookURL
	}
	if addr, err := netip.ParseAddr(hostname); err == nil {
		if err := ValidateResolvedAddress(addr); err != nil {
			return "", errInvalidWebhookURL
		}
	} else if isLegacyIPv4NumericHostname(hostname) {
		return "", errInvalidWebhookURL
	}
	return value, nil
}

func ValidateResolvedAddress(addr netip.Addr) error {
	addr = addr.Unmap()
	if !addr.IsValid() ||
		addr.IsLoopback() ||
		addr.IsPrivate() ||
		addr.IsLinkLocalUnicast() ||
		addr.IsMulticast() ||
		addr.IsUnspecified() {
		return errUnsafeWebhookAddress
	}
	for _, prefix := range nonPublicWebhookAddressPrefixes {
		if prefix.Contains(addr) {
			return errUnsafeWebhookAddress
		}
	}
	return nil
}

func isLegacyIPv4NumericHostname(hostname string) bool {
	parts := strings.Split(hostname, ".")
	if len(parts) == 0 || len(parts) > 4 {
		return false
	}
	for _, part := range parts {
		if !isIPv4NumericPart(part) {
			return false
		}
	}
	return true
}

func isIPv4NumericPart(part string) bool {
	if part == "" {
		return false
	}
	if len(part) > 2 && part[0] == '0' && (part[1] == 'x' || part[1] == 'X') {
		for _, r := range part[2:] {
			if !isASCIIHexDigit(r) {
				return false
			}
		}
		return true
	}
	for _, r := range part {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isASCIIHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

func NormalizeEvents(raw []Event) ([]Event, error) {
	out := make([]Event, 0, len(raw))
	for _, event := range raw {
		if !slices.Contains(supportedEvents, event) {
			return nil, errors.New("events_active contains unsupported event")
		}
		if !slices.Contains(out, event) {
			out = append(out, event)
		}
	}
	return out, nil
}

func EncodeEvents(events []Event) (datatypes.JSON, error) {
	normalized, err := NormalizeEvents(events)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(normalized)
	return datatypes.JSON(raw), err
}

func DecodeEvents(raw datatypes.JSON) []Event {
	var events []Event
	if err := json.Unmarshal(raw, &events); err != nil {
		return []Event{}
	}
	normalized, err := NormalizeEvents(events)
	if err != nil {
		return []Event{}
	}
	return normalized
}
