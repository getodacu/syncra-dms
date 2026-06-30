package webhooks

import (
	"net/netip"
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestWebhookBeforeCreateAssignsIDAndTableName(t *testing.T) {
	webhook := Webhook{}
	if err := webhook.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v", err)
	}
	if webhook.ID == uuid.Nil {
		t.Fatal("BeforeCreate() left nil ID")
	}
	if (Webhook{}).TableName() != "webhooks" {
		t.Fatalf("TableName() = %q, want webhooks", (Webhook{}).TableName())
	}
}

func TestValidateURLAcceptsAbsoluteHTTPAndHTTPSOnly(t *testing.T) {
	for _, raw := range []string{
		"http://example.com/webhook",
		"https://example.com/webhook",
		"http://93.184.216.34/webhook",
		" https://example.com/webhook ",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := ValidateURL(raw)
			if err != nil {
				t.Fatalf("ValidateURL() error = %v", err)
			}
			if got == "" || got[0] == ' ' || got[len(got)-1] == ' ' {
				t.Fatalf("ValidateURL() = %q, want trimmed non-empty URL", got)
			}
		})
	}

	for _, raw := range []string{
		"",
		"example.com/webhook",
		"/webhook",
		"http://:80/path",
		"ftp://example.com/webhook",
		"mailto:ops@example.com",
	} {
		t.Run(raw, func(t *testing.T) {
			if _, err := ValidateURL(raw); err == nil {
				t.Fatal("ValidateURL() error = nil, want validation error")
			}
		})
	}
}

func TestValidateURLRejectsUnsafeHosts(t *testing.T) {
	for _, raw := range []string{
		"http://localhost/webhook",
		"http://LOCALHOST/webhook",
		"http://127.0.0.1/webhook",
		"http://127.0.0.1./webhook",
		"http://2130706433/webhook",
		"http://0177.0.0.1/webhook",
		"http://127.1/webhook",
		"http://10.0.0.1/webhook",
		"http://169.254.169.254/latest/meta-data",
	} {
		t.Run(raw, func(t *testing.T) {
			_, err := ValidateURL(raw)
			if err == nil {
				t.Fatal("ValidateURL() error = nil, want validation error")
			}
			if err.Error() != "url must be an absolute http or https URL" {
				t.Fatalf("ValidateURL() error = %q", err.Error())
			}
		})
	}
}

func TestValidateResolvedAddressAllowsPublicAddresses(t *testing.T) {
	for _, raw := range []string{
		"8.8.8.8",
		"9.255.255.255",
		"11.0.0.0",
		"93.184.216.34",
		"100.63.255.255",
		"100.128.0.1",
		"172.15.255.255",
		"172.32.0.1",
		"192.0.1.1",
		"198.17.255.255",
		"198.20.0.1",
		"203.0.112.255",
		"203.0.114.1",
		"223.255.255.254",
		"2606:2800:220:1:248:1893:25c8:1946",
		"2001:200::1",
		"2003::1",
	} {
		t.Run(raw, func(t *testing.T) {
			if err := ValidateResolvedAddress(netip.MustParseAddr(raw)); err != nil {
				t.Fatalf("ValidateResolvedAddress() error = %v", err)
			}
		})
	}
}

func TestValidateResolvedAddressRejectsUnsafeAddresses(t *testing.T) {
	for _, raw := range []string{
		"0.0.0.0",
		"10.0.0.1",
		"100.64.0.1",
		"127.0.0.1",
		"169.254.169.254",
		"172.16.0.1",
		"192.0.0.1",
		"192.0.2.1",
		"192.88.99.1",
		"192.168.0.1",
		"198.18.0.1",
		"198.51.100.1",
		"203.0.113.1",
		"224.0.0.1",
		"240.0.0.1",
		"255.255.255.255",
		"::",
		"::1",
		"::ffff:127.0.0.1",
		"64:ff9b::1",
		"64:ff9b:1::1",
		"100::1",
		"2001::1",
		"2001:db8::1",
		"2002::1",
		"fc00::1",
		"fe80::1",
		"ff02::1",
	} {
		t.Run(raw, func(t *testing.T) {
			if err := ValidateResolvedAddress(netip.MustParseAddr(raw)); err == nil {
				t.Fatal("ValidateResolvedAddress() error = nil, want validation error")
			}
		})
	}
}

func TestNormalizeEventsAllowsEmptyAndRejectsUnknown(t *testing.T) {
	events, err := NormalizeEvents([]Event{EventJobStarted, EventJobSucceeded})
	if err != nil {
		t.Fatalf("NormalizeEvents() error = %v", err)
	}
	if !slices.Equal(events, []Event{EventJobStarted, EventJobSucceeded}) {
		t.Fatalf("events = %#v", events)
	}

	if _, err := NormalizeEvents([]Event{"job.unknown"}); err == nil {
		t.Fatal("NormalizeEvents() expected unknown event error")
	}

	empty, err := NormalizeEvents(nil)
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty events = %#v err=%v", empty, err)
	}
}

func TestEncodeDecodeEventsRoundTripAndNormalize(t *testing.T) {
	encoded, err := EncodeEvents([]Event{EventJobStarted, EventJobStarted, EventJobFailed})
	if err != nil {
		t.Fatalf("EncodeEvents() error = %v", err)
	}

	decoded := DecodeEvents(encoded)
	if !slices.Equal(decoded, []Event{EventJobStarted, EventJobFailed}) {
		t.Fatalf("DecodeEvents() = %#v", decoded)
	}
	if string(encoded) != `["job.started","job.failed"]` {
		t.Fatalf("EncodeEvents() = %s, want normalized unique events", encoded)
	}

	empty, err := EncodeEvents(nil)
	if err != nil {
		t.Fatalf("EncodeEvents(nil) error = %v", err)
	}
	if string(empty) != "[]" {
		t.Fatalf("EncodeEvents(nil) = %s, want []", empty)
	}

	if _, err := EncodeEvents([]Event{"job.unknown"}); err == nil {
		t.Fatal("EncodeEvents() error = nil, want unsupported event error")
	}
}
