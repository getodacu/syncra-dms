package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestAuthClientIntegrationSignupVerifySigninSessionSignout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	router := NewRouter(&Handler{
		DB:                  db,
		BetterAuthSecret:    testAuthSecret,
		AuthDeliveryToken:   testDeliveryToken,
		InternalAPIToken:    testInternalAPIToken,
		AuthSessionTTL:      7 * 24 * time.Hour,
		AuthVerificationTTL: 5 * time.Minute,
		AuthCookieSecure:    false,
	})

	apiServer := httptest.NewServer(router)
	defer apiServer.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("create cookie jar: %v", err)
	}
	client := apiServer.Client()
	client.Jar = jar
	client.Timeout = 5 * time.Second

	signupRes := postJSON(t, client, apiServer.URL, "/api/auth/sign-up/email", map[string]any{
		"name":     "Ada Lovelace",
		"email":    "ada@example.com",
		"password": "password1234",
	}, map[string]string{authDeliveryHeader: testDeliveryToken})
	requireHTTPStatus(t, signupRes, http.StatusOK)
	signup := decodeHTTPJSON[signUpEmailResponse](t, signupRes)
	if signup.VerificationCode == "" {
		t.Fatalf("verification code is empty: %#v", signup)
	}

	verifyRes := postJSON(t, client, apiServer.URL, "/api/auth/email-otp/verify-email", map[string]any{
		"email": "ada@example.com",
		"otp":   signup.VerificationCode,
	}, nil)
	requireHTTPStatus(t, verifyRes, http.StatusOK)
	_ = decodeHTTPJSON[verifyEmailOTPResponse](t, verifyRes)

	signinRes := postJSON(t, client, apiServer.URL, "/api/auth/sign-in/email", map[string]any{
		"email":    "ada@example.com",
		"password": "password1234",
	}, nil)
	requireHTTPStatus(t, signinRes, http.StatusOK)
	signin := decodeHTTPJSON[authSessionPayload](t, signinRes)
	if signin.Session.Token == "" {
		t.Fatalf("sign-in session token is empty: %#v", signin.Session)
	}

	serverURL, err := url.Parse(apiServer.URL)
	if err != nil {
		t.Fatalf("parse test server URL: %v", err)
	}
	sessionCookie := cookieFromJar(jar, serverURL, authSessionCookieName)
	if sessionCookie == nil {
		t.Fatalf("cookie jar missing %s after sign-in", authSessionCookieName)
	}
	if sessionCookie.Value != signin.Session.Token {
		t.Fatalf("session cookie = %q, want %q", sessionCookie.Value, signin.Session.Token)
	}

	sessionReq, err := http.NewRequest(http.MethodGet, apiServer.URL+"/api/auth/get-session", nil)
	if err != nil {
		t.Fatalf("create get-session request: %v", err)
	}
	sessionReq.Header.Set(internalAPIHeader, testInternalAPIToken)
	sessionRes, err := client.Do(sessionReq)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	requireHTTPStatus(t, sessionRes, http.StatusOK)
	session := decodeHTTPJSON[authSessionPayload](t, sessionRes)
	if session.User.Email != "ada@example.com" {
		t.Fatalf("session user email = %q, want ada@example.com", session.User.Email)
	}
	if session.Session.Token == "" {
		t.Fatalf("session token is empty: %#v", session.Session)
	}

	signoutRes := postJSON(t, client, apiServer.URL, "/api/auth/sign-out", map[string]any{}, nil)
	requireHTTPStatus(t, signoutRes, http.StatusOK)
	signout := decodeHTTPJSON[struct {
		Success bool `json:"success"`
	}](t, signoutRes)
	if !signout.Success {
		t.Fatalf("sign-out success = false")
	}
	if cookie := cookieFromJar(jar, serverURL, authSessionCookieName); cookie != nil {
		t.Fatalf("session cookie still present after sign-out: %#v", cookie)
	}

	signedOutReq, err := http.NewRequest(http.MethodGet, apiServer.URL+"/api/auth/get-session", nil)
	if err != nil {
		t.Fatalf("create signed-out get-session request: %v", err)
	}
	signedOutReq.Header.Set(internalAPIHeader, testInternalAPIToken)
	signedOutRes, err := client.Do(signedOutReq)
	if err != nil {
		t.Fatalf("get session after sign-out: %v", err)
	}
	requireHTTPStatus(t, signedOutRes, http.StatusOK)
	signedOut := decodeHTTPJSON[*authSessionPayload](t, signedOutRes)
	if signedOut != nil {
		t.Fatalf("signed-out session = %#v, want nil", signedOut)
	}
}

func postJSON(t *testing.T, client *http.Client, baseURL string, path string, body any, headers map[string]string) *http.Response {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal JSON body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create POST %s request: %v", path, err)
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.HasPrefix(path, "/api/") {
		req.Header.Set(internalAPIHeader, testInternalAPIToken)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return res
}

func decodeHTTPJSON[T any](t *testing.T, res *http.Response) T {
	t.Helper()
	defer res.Body.Close()
	var out T
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatalf("decode %s JSON response: %v", res.Request.URL.Path, err)
	}
	return out
}

func requireHTTPStatus(t *testing.T, res *http.Response, want int) {
	t.Helper()
	if res.StatusCode == want {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s status = %d, want %d; read body: %v", res.Request.URL.Path, res.StatusCode, want, err)
	}
	_ = res.Body.Close()
	t.Fatalf("%s status = %d, want %d; body=%s", res.Request.URL.Path, res.StatusCode, want, string(body))
}

func cookieFromJar(jar *cookiejar.Jar, u *url.URL, name string) *http.Cookie {
	for _, cookie := range jar.Cookies(u) {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
