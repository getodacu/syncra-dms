package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/webhooks"
)

const testWebhookPrivateKey = "test-webhook-private-key-32-bytes"

type webhookTestResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       string    `json:"user_id"`
	URL          string    `json:"url"`
	EventsActive []string  `json:"events_active"`
	HasSecret    bool      `json:"has_secret"`
	SecretKey    string    `json:"secret_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type webhookEnvelopeTestResponse struct {
	Webhook *webhookTestResponse `json:"webhook"`
}

type deleteWebhookTestResponse struct {
	DeletedID    uuid.UUID `json:"deleted_id"`
	DeletedCount int       `json:"deleted_count"`
}

func testWebhookAuthRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	_, db := testAuthRouter(t)
	h := &Handler{
		DB:                  db,
		BetterAuthSecret:    testAuthSecret,
		AuthDeliveryToken:   testDeliveryToken,
		InternalAPIToken:    testInternalAPIToken,
		AuthSessionTTL:      7 * 24 * time.Hour,
		AuthVerificationTTL: 5 * time.Minute,
		AuthCookieSecure:    false,
		AppPrivateKey:       testWebhookPrivateKey,
	}
	return NewRouter(h), db
}

func TestWebhookGetReturnsNullWhenAbsent(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-absent@example.com")

	w := authJSON(t, router, http.MethodGet, "/api/auth/webhook/"+user.ID, "", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[webhookEnvelopeTestResponse](t, w)
	if got.Webhook != nil {
		t.Fatalf("webhook = %#v, want nil", got.Webhook)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode raw response: %v", err)
	}
	if string(raw["webhook"]) != "null" {
		t.Fatalf("raw webhook = %s, want null", raw["webhook"])
	}
}

func TestWebhookCreateReturnsSecretOnceAndStoresEncryptedSecret(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-create@example.com")

	w := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/syncra",
		"events_active":["job.started","job.succeeded"]
	}`, nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[webhookTestResponse](t, w)
	if got.ID == uuid.Nil || got.UserID != user.ID || got.URL != "https://hooks.example.com/syncra" {
		t.Fatalf("unexpected webhook response: %#v", got)
	}
	if got.SecretKey == "" {
		t.Fatal("secret_key is empty, want one-time plaintext secret")
	}
	if !got.HasSecret {
		t.Fatal("has_secret = false, want true")
	}
	if !equalStringSlices(got.EventsActive, []string{"job.started", "job.succeeded"}) {
		t.Fatalf("events_active = %#v, want started/succeeded", got.EventsActive)
	}

	var stored webhooks.Webhook
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load webhook: %v", err)
	}
	if stored.SecretKey == "" || stored.SecretKey == got.SecretKey {
		t.Fatalf("stored secret_key = %q, want encrypted value different from plaintext", stored.SecretKey)
	}
	decrypted, err := webhooks.DecryptSecret(testWebhookPrivateKey, stored.SecretKey)
	if err != nil {
		t.Fatalf("decrypt stored secret: %v", err)
	}
	if decrypted != got.SecretKey {
		t.Fatalf("decrypted secret = %q, want returned secret", decrypted)
	}

	get := authJSON(t, router, http.MethodGet, "/api/auth/webhook/"+user.ID, "", nil)
	if get.Code != http.StatusOK {
		t.Fatalf("get status = %d body=%s", get.Code, get.Body.String())
	}
	existing := decodeAuthResponse[webhookEnvelopeTestResponse](t, get)
	if existing.Webhook == nil || existing.Webhook.SecretKey != "" || !existing.Webhook.HasSecret {
		t.Fatalf("get existing webhook exposed secret or missed has_secret: %#v", existing.Webhook)
	}
	assertWebhookRawOmitsSecretKey(t, get.Body.Bytes())
}

func TestWebhookSecondPostUpdatesWithoutReturningSecret(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-update@example.com")

	first := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/old",
		"events_active":["job.started"]
	}`, nil)
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d body=%s", first.Code, first.Body.String())
	}
	created := decodeAuthResponse[webhookTestResponse](t, first)
	var before webhooks.Webhook
	if err := db.First(&before, "id = ?", created.ID).Error; err != nil {
		t.Fatalf("load created webhook: %v", err)
	}

	second := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/new",
		"events_active":["job.failed","job.succeeded"]
	}`, nil)

	if second.Code != http.StatusOK {
		t.Fatalf("second status = %d body=%s", second.Code, second.Body.String())
	}
	got := decodeAuthResponse[webhookTestResponse](t, second)
	if got.ID != created.ID || got.URL != "https://hooks.example.com/new" || got.SecretKey != "" {
		t.Fatalf("unexpected update response: %#v", got)
	}
	if !equalStringSlices(got.EventsActive, []string{"job.failed", "job.succeeded"}) {
		t.Fatalf("events_active = %#v, want failed/succeeded", got.EventsActive)
	}
	assertWebhookRawOmitsSecretKey(t, second.Body.Bytes())

	var after webhooks.Webhook
	if err := db.First(&after, "id = ?", created.ID).Error; err != nil {
		t.Fatalf("load updated webhook: %v", err)
	}
	if after.SecretKey != before.SecretKey {
		t.Fatal("stored encrypted secret changed on update")
	}
	if after.URL != "https://hooks.example.com/new" {
		t.Fatalf("stored url = %q, want updated url", after.URL)
	}
}

func TestWebhookSecondPostDoesNotResurrectDeletedWebhook(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-update-delete-race@example.com")
	hook := createStoredWebhook(t, db, user.ID, "https://hooks.example.com/old", []webhooks.Event{webhooks.EventJobStarted})
	deleteWebhookAfterNextFind(t, db, hook, "syncra:test_webhook_update_delete_after_find")

	w := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/new",
		"events_active":["job.failed"]
	}`, nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s, want not found after concurrent delete", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[ErrorResponse](t, w)
	if got.Error != "webhook not found" {
		t.Fatalf("error = %q, want webhook not found", got.Error)
	}
	assertWebhookRowCount(t, db, hook.ID, 0)
}

func TestWebhookCreateUniqueConflictFallsBackToUpdateWithoutSecret(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-create-race@example.com")
	existingPlaintextSecret := "race-winner-secret"
	existingEncryptedSecret, err := webhooks.EncryptSecret(testWebhookPrivateKey, existingPlaintextSecret)
	if err != nil {
		t.Fatalf("encrypt existing secret: %v", err)
	}
	existingEvents, err := webhooks.EncodeEvents([]webhooks.Event{webhooks.EventJobStarted})
	if err != nil {
		t.Fatalf("encode existing events: %v", err)
	}
	existingID := uuid.New()
	callbackName := "syncra:test_webhook_unique_conflict"
	injected := false
	if err := db.Callback().Create().Before("gorm:create").Register(callbackName, func(tx *gorm.DB) {
		hook, ok := tx.Statement.Dest.(*webhooks.Webhook)
		if !ok || hook.UserID != user.ID || injected {
			return
		}
		injected = true
		now := time.Now().UTC()
		if err := tx.Exec(
			`INSERT INTO webhooks (id, user_id, url, secret_key, events_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?::jsonb, ?, ?)`,
			existingID,
			user.ID,
			"https://hooks.example.com/race-winner",
			existingEncryptedSecret,
			string(existingEvents),
			now,
			now,
		).Error; err != nil {
			tx.AddError(err)
			return
		}
		tx.AddError(gorm.ErrDuplicatedKey)
	}); err != nil {
		t.Fatalf("register create callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Create().Remove(callbackName); err != nil {
			t.Errorf("remove create callback: %v", err)
		}
	})

	w := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/race-loser",
		"events_active":["job.failed"]
	}`, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s, want conflict fallback update", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[webhookTestResponse](t, w)
	if got.ID != existingID || got.UserID != user.ID || got.URL != "https://hooks.example.com/race-loser" {
		t.Fatalf("unexpected conflict fallback response: %#v", got)
	}
	if got.SecretKey != "" {
		t.Fatalf("conflict fallback exposed secret_key: %#v", got)
	}
	if !got.HasSecret {
		t.Fatal("has_secret = false, want true")
	}
	if !equalStringSlices(got.EventsActive, []string{"job.failed"}) {
		t.Fatalf("events_active = %#v, want failed", got.EventsActive)
	}
	assertWebhookRawOmitsSecretKey(t, w.Body.Bytes())

	var count int64
	if err := db.Model(&webhooks.Webhook{}).Where("user_id = ?", user.ID).Count(&count).Error; err != nil {
		t.Fatalf("count webhooks: %v", err)
	}
	if count != 1 {
		t.Fatalf("webhook count = %d, want 1", count)
	}
	var stored webhooks.Webhook
	if err := db.First(&stored, "user_id = ?", user.ID).Error; err != nil {
		t.Fatalf("load webhook: %v", err)
	}
	if stored.ID != existingID || stored.URL != "https://hooks.example.com/race-loser" {
		t.Fatalf("unexpected stored webhook after fallback: %#v", stored)
	}
	if stored.SecretKey != existingEncryptedSecret {
		t.Fatal("conflict fallback changed stored encrypted secret")
	}
	storedEvents := webhooks.DecodeEvents(stored.EventsActive)
	if len(storedEvents) != 1 || storedEvents[0] != webhooks.EventJobFailed {
		t.Fatalf("stored events = %#v, want job.failed", storedEvents)
	}
}

func TestWebhookRegenerateSecretReturnsDifferentOneTimeSecret(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-regenerate@example.com")

	create := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/regenerate",
		"events_active":["job.failed"]
	}`, nil)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", create.Code, create.Body.String())
	}
	created := decodeAuthResponse[webhookTestResponse](t, create)
	var before webhooks.Webhook
	if err := db.First(&before, "id = ?", created.ID).Error; err != nil {
		t.Fatalf("load created webhook: %v", err)
	}

	regen := authJSON(t, router, http.MethodPatch, "/api/auth/webhook/"+user.ID+"/secret", "", nil)

	if regen.Code != http.StatusOK {
		t.Fatalf("regenerate status = %d body=%s", regen.Code, regen.Body.String())
	}
	got := decodeAuthResponse[webhookTestResponse](t, regen)
	if got.ID != created.ID || got.SecretKey == "" || got.SecretKey == created.SecretKey {
		t.Fatalf("unexpected regenerate response: %#v created_secret=%q", got, created.SecretKey)
	}
	var after webhooks.Webhook
	if err := db.First(&after, "id = ?", created.ID).Error; err != nil {
		t.Fatalf("load regenerated webhook: %v", err)
	}
	if after.SecretKey == before.SecretKey || after.SecretKey == got.SecretKey {
		t.Fatal("stored secret was not regenerated as encrypted material")
	}
	decrypted, err := webhooks.DecryptSecret(testWebhookPrivateKey, after.SecretKey)
	if err != nil {
		t.Fatalf("decrypt regenerated secret: %v", err)
	}
	if decrypted != got.SecretKey {
		t.Fatalf("decrypted secret = %q, want regenerated secret", decrypted)
	}
}

func TestWebhookRegenerateSecretDoesNotResurrectDeletedWebhook(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-regenerate-delete-race@example.com")
	hook := createStoredWebhook(t, db, user.ID, "https://hooks.example.com/regenerate", []webhooks.Event{webhooks.EventJobFailed})
	deleteWebhookAfterNextFind(t, db, hook, "syncra:test_webhook_regenerate_delete_after_find")

	regen := authJSON(t, router, http.MethodPatch, "/api/auth/webhook/"+user.ID+"/secret", "", nil)

	if regen.Code != http.StatusNotFound {
		t.Fatalf("regenerate status = %d body=%s, want not found after concurrent delete", regen.Code, regen.Body.String())
	}
	got := decodeAuthResponse[ErrorResponse](t, regen)
	if got.Error != "webhook not found" {
		t.Fatalf("error = %q, want webhook not found", got.Error)
	}
	assertWebhookRowCount(t, db, hook.ID, 0)
}

func TestWebhookDeleteScopesToUser(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-delete@example.com")
	other := createTestUser(t, db, "webhook-delete-other@example.com")
	userWebhook := createStoredWebhook(t, db, user.ID, "https://hooks.example.com/delete-user", []webhooks.Event{webhooks.EventJobStarted})
	otherWebhook := createStoredWebhook(t, db, other.ID, "https://hooks.example.com/delete-other", []webhooks.Event{webhooks.EventJobFailed})

	w := authJSON(t, router, http.MethodDelete, "/api/auth/webhook?user_id="+url.QueryEscape(user.ID), "", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[deleteWebhookTestResponse](t, w)
	if got.DeletedID != userWebhook.ID || got.DeletedCount != 1 {
		t.Fatalf("delete response = %#v, want user webhook id and count 1", got)
	}
	var deletedCount int64
	if err := db.Model(&webhooks.Webhook{}).Where("id = ?", userWebhook.ID).Count(&deletedCount).Error; err != nil {
		t.Fatalf("count deleted webhook: %v", err)
	}
	if deletedCount != 0 {
		t.Fatalf("deleted webhook count = %d, want 0", deletedCount)
	}
	var otherCount int64
	if err := db.Model(&webhooks.Webhook{}).Where("id = ?", otherWebhook.ID).Count(&otherCount).Error; err != nil {
		t.Fatalf("count other webhook: %v", err)
	}
	if otherCount != 1 {
		t.Fatalf("other webhook count = %d, want 1", otherCount)
	}
}

func TestWebhookValidationRejectsInvalidURLAndUnsupportedEvents(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-validation@example.com")

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "invalid url",
			body: `{
				"user_id":"` + user.ID + `",
				"url":"http://localhost/webhook",
				"events_active":["job.started"]
			}`,
			want: "url must be an absolute http or https URL",
		},
		{
			name: "unsupported event",
			body: `{
				"user_id":"` + user.ID + `",
				"url":"https://hooks.example.com/syncra",
				"events_active":["job.deleted"]
			}`,
			want: "events_active contains unsupported event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := authJSON(t, router, http.MethodPost, "/api/auth/webhook", tt.body, nil)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			got := decodeAuthResponse[ErrorResponse](t, w)
			if got.Error != tt.want {
				t.Fatalf("error = %q, want %q", got.Error, tt.want)
			}
		})
	}
}

func TestWebhookEmptyEventsActiveSucceeds(t *testing.T) {
	router, db := testWebhookAuthRouter(t)
	user := createTestUser(t, db, "webhook-empty-events@example.com")

	w := authJSON(t, router, http.MethodPost, "/api/auth/webhook", `{
		"user_id":"`+user.ID+`",
		"url":"https://hooks.example.com/paused",
		"events_active":[]
	}`, nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	got := decodeAuthResponse[webhookTestResponse](t, w)
	if got.SecretKey == "" || !got.HasSecret {
		t.Fatalf("secret response fields = secret:%q has_secret:%v", got.SecretKey, got.HasSecret)
	}
	if len(got.EventsActive) != 0 {
		t.Fatalf("events_active = %#v, want empty", got.EventsActive)
	}
	var stored webhooks.Webhook
	if err := db.First(&stored, "id = ?", got.ID).Error; err != nil {
		t.Fatalf("load webhook: %v", err)
	}
	if string(stored.EventsActive) != "[]" {
		t.Fatalf("stored events_active = %s, want []", string(stored.EventsActive))
	}
}

func createStoredWebhook(t *testing.T, db *gorm.DB, userID string, rawURL string, events []webhooks.Event) webhooks.Webhook {
	t.Helper()
	encryptedSecret, err := webhooks.EncryptSecret(testWebhookPrivateKey, "stored-"+uuid.NewString())
	if err != nil {
		t.Fatalf("encrypt test webhook secret: %v", err)
	}
	encodedEvents, err := webhooks.EncodeEvents(events)
	if err != nil {
		t.Fatalf("encode test webhook events: %v", err)
	}
	hook := webhooks.Webhook{
		UserID:       userID,
		URL:          rawURL,
		SecretKey:    encryptedSecret,
		EventsActive: encodedEvents,
	}
	if err := db.Create(&hook).Error; err != nil {
		t.Fatalf("create stored webhook: %v", err)
	}
	return hook
}

func deleteWebhookAfterNextFind(t *testing.T, db *gorm.DB, hook webhooks.Webhook, callbackName string) {
	t.Helper()
	deleted := false
	if err := db.Callback().Query().After("gorm:query").Register(callbackName, func(tx *gorm.DB) {
		if deleted {
			return
		}
		loaded, ok := tx.Statement.Dest.(*webhooks.Webhook)
		if !ok || loaded.ID != hook.ID || loaded.UserID != hook.UserID {
			return
		}
		deleted = true
		if err := tx.Session(&gorm.Session{NewDB: true}).
			Exec("DELETE FROM webhooks WHERE id = ? AND user_id = ?", hook.ID, hook.UserID).
			Error; err != nil {
			tx.AddError(err)
		}
	}); err != nil {
		t.Fatalf("register query callback: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Callback().Query().Remove(callbackName); err != nil {
			t.Errorf("remove query callback: %v", err)
		}
	})
}

func assertWebhookRowCount(t *testing.T, db *gorm.DB, id uuid.UUID, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&webhooks.Webhook{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("count webhook rows: %v", err)
	}
	if count != want {
		t.Fatalf("webhook row count = %d, want %d", count, want)
	}
}

func assertWebhookRawOmitsSecretKey(t *testing.T, raw []byte) {
	t.Helper()
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		t.Fatalf("decode raw webhook JSON: %v", err)
	}
	if containsJSONKey(value, "secret_key") {
		t.Fatalf("raw response includes secret_key: %s", string(raw))
	}
}

func containsJSONKey(value any, key string) bool {
	switch typed := value.(type) {
	case map[string]any:
		if _, ok := typed[key]; ok {
			return true
		}
		for _, child := range typed {
			if containsJSONKey(child, key) {
				return true
			}
		}
	case []any:
		for _, child := range typed {
			if containsJSONKey(child, key) {
				return true
			}
		}
	}
	return false
}

func equalStringSlices(got []string, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
