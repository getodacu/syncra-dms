package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"ai.ro/syncra/internal/ocr"
	"github.com/gin-gonic/gin"
)

func TestAPIClientIntegrationCreatesSchemaAndOCRDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := apiPostgresTx(t)
	file := validPNGBytes()
	const fakeMistralResponse = `{
		"pages":[{"index":0,"markdown":"# Integrated OCR","images":[]}],
		"model":"mistral-ocr-latest",
		"document_annotation":"{\"Furnizor\":\"Acme\"}",
		"usage_info":{"pages_processed":1}
	}`

	upstreamErrs := make(chan error, 1)
	var upstreamCalls atomic.Int32
	recordUpstreamErr := func(format string, args ...any) {
		select {
		case upstreamErrs <- fmt.Errorf(format, args...):
		default:
		}
	}
	writeMistralResponse := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fakeMistralResponse))
	}
	t.Cleanup(func() {
		select {
		case err := <-upstreamErrs:
			t.Errorf("Mistral upstream assertion: %v", err)
		default:
		}
	})

	mistralUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamCalls.Add(1)
		if r.Method != http.MethodPost || r.URL.Path != "/v1/ocr" {
			recordUpstreamErr("unexpected Mistral request: %s %s", r.Method, r.URL.Path)
			writeMistralResponse(w)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			recordUpstreamErr("Authorization = %q", got)
			writeMistralResponse(w)
			return
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			recordUpstreamErr("Content-Type = %q", got)
			writeMistralResponse(w)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			recordUpstreamErr("read Mistral request: %v", err)
			writeMistralResponse(w)
			return
		}
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			recordUpstreamErr("decode Mistral payload: %v", err)
			writeMistralResponse(w)
			return
		}
		if payload["model"] != "mistral-ocr-latest" {
			recordUpstreamErr("model = %v", payload["model"])
			writeMistralResponse(w)
			return
		}
		if payload["include_image_base64"] != true {
			recordUpstreamErr("include_image_base64 = %v", payload["include_image_base64"])
			writeMistralResponse(w)
			return
		}

		document, ok := payload["document"].(map[string]any)
		if !ok {
			recordUpstreamErr("document = %#v", payload["document"])
			writeMistralResponse(w)
			return
		}
		if document["type"] != "image_url" {
			recordUpstreamErr("document type = %v", document["type"])
			writeMistralResponse(w)
			return
		}
		imageURL, ok := document["image_url"].(string)
		if !ok || !strings.HasPrefix(imageURL, "data:image/png;base64,") {
			recordUpstreamErr("image_url = %#v", document["image_url"])
			writeMistralResponse(w)
			return
		}
		encodedImage := strings.TrimPrefix(imageURL, "data:image/png;base64,")
		imageData, err := base64.StdEncoding.DecodeString(encodedImage)
		if err != nil {
			recordUpstreamErr("decode image_url base64: %v", err)
			writeMistralResponse(w)
			return
		}
		if !bytes.Equal(imageData, file) {
			recordUpstreamErr("image data = %v, want %v", imageData, file)
			writeMistralResponse(w)
			return
		}

		format, ok := payload["document_annotation_format"].(map[string]any)
		if !ok {
			recordUpstreamErr("document_annotation_format = %#v", payload["document_annotation_format"])
			writeMistralResponse(w)
			return
		}
		if format["type"] != "json_schema" {
			recordUpstreamErr("document_annotation_format type = %v", format["type"])
			writeMistralResponse(w)
			return
		}
		schemaPayload, ok := format["json_schema"].(map[string]any)
		if !ok {
			recordUpstreamErr("json_schema = %#v", format["json_schema"])
			writeMistralResponse(w)
			return
		}
		if schemaPayload["name"] != "response_schema" {
			recordUpstreamErr("schema name = %v", schemaPayload["name"])
			writeMistralResponse(w)
			return
		}
		if schemaPayload["strict"] != true {
			recordUpstreamErr("strict = %v", schemaPayload["strict"])
			writeMistralResponse(w)
			return
		}
		schemaBytes, err := json.Marshal(schemaPayload["schema"])
		if err != nil {
			recordUpstreamErr("marshal schema payload: %v", err)
			writeMistralResponse(w)
			return
		}
		var gotSchema any
		if err := json.Unmarshal(schemaBytes, &gotSchema); err != nil {
			recordUpstreamErr("decode schema payload: %v", err)
			writeMistralResponse(w)
			return
		}
		var wantSchema any
		if err := json.Unmarshal([]byte(`{"type":"object","properties":{"Furnizor":{"type":"string"}},"required":["Furnizor"]}`), &wantSchema); err != nil {
			recordUpstreamErr("decode expected schema payload: %v", err)
			writeMistralResponse(w)
			return
		}
		if !reflect.DeepEqual(gotSchema, wantSchema) {
			recordUpstreamErr("schema payload = %s", string(schemaBytes))
			writeMistralResponse(w)
			return
		}

		writeMistralResponse(w)
	}))
	defer mistralUpstream.Close()

	apiServer := httptest.NewServer(NewRouter(&Handler{
		DB:               db,
		MaxUploadBytes:   20 << 20,
		MistralAPIKey:    "test-key",
		MistralBaseURL:   mistralUpstream.URL,
		MistralModel:     "mistral-ocr-latest",
		InternalAPIToken: testInternalAPIToken,
	}))
	defer apiServer.Close()

	apiHTTPClient := apiServer.Client()
	apiHTTPClient.Timeout = 5 * time.Second
	client := newTestAPIClient(apiServer.URL, apiHTTPClient)
	schema := client.createSchema(t, `{
		"name":"invoice",
		"description":"Supplier extraction",
		"schema":{"type":"object","properties":{"Furnizor":{"type":"string"}},"required":["Furnizor"]},
		"strict":true
	}`)

	doc := client.ocr(t, map[string]string{
		"schema_id": schema.ID.String(),
	}, "scan.png", file)
	if got := upstreamCalls.Load(); got != 1 {
		t.Fatalf("Mistral upstream calls = %d, want 1", got)
	}
	select {
	case err := <-upstreamErrs:
		t.Fatalf("Mistral upstream assertion: %v", err)
	default:
	}

	if doc.Markdown != "# Integrated OCR" {
		t.Fatalf("markdown = %q", doc.Markdown)
	}
	if doc.SchemaID == nil || *doc.SchemaID != schema.ID {
		t.Fatalf("schema id = %#v, want %s", doc.SchemaID, schema.ID)
	}
	assertJSONEqual(t, doc.AnnotationJSON, `{"Furnizor":"Acme"}`)

	var stored ocr.OCRDocument
	if err := db.First(&stored, "id = ?", doc.ID).Error; err != nil {
		t.Fatalf("load stored OCR document: %v", err)
	}
	if stored.SchemaID == nil || *stored.SchemaID != schema.ID {
		t.Fatalf("stored schema id = %#v, want %s", stored.SchemaID, schema.ID)
	}
	if stored.OriginalFilename != "scan.png" {
		t.Fatalf("stored original filename = %q", stored.OriginalFilename)
	}
	if stored.MimeType != "image/png" {
		t.Fatalf("stored mime type = %q", stored.MimeType)
	}
	if stored.FileSize != int64(len(file)) {
		t.Fatalf("stored file size = %d, want %d", stored.FileSize, len(file))
	}
	if stored.DocumentHash == "" || doc.DocumentHash != stored.DocumentHash {
		t.Fatalf("document hash response=%q stored=%q", doc.DocumentHash, stored.DocumentHash)
	}
	if doc.Cached {
		t.Fatal("cached = true, want false")
	}
	if stored.Markdown != "# Integrated OCR" {
		t.Fatalf("stored markdown = %q", stored.Markdown)
	}
	if len(stored.InlineSchemaJSON) != 0 {
		t.Fatalf("stored inline schema = %s, want empty", string(stored.InlineSchemaJSON))
	}
	assertJSONEqual(t, json.RawMessage(stored.AnnotationJSON), `{"Furnizor":"Acme"}`)
	assertJSONEqual(t, json.RawMessage(stored.RawResponseJSON), fakeMistralResponse)
}

type testAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func newTestAPIClient(baseURL string, httpClient *http.Client) *testAPIClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &testAPIClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *testAPIClient) createSchema(t *testing.T, body string) SchemaResponse {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/ocr/schemas", strings.NewReader(body))
	if err != nil {
		t.Fatalf("create schema request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(internalAPIHeader, testInternalAPIToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}
	defer res.Body.Close()
	requireAPIClientStatus(t, res, http.StatusCreated)

	var schema SchemaResponse
	if err := json.NewDecoder(res.Body).Decode(&schema); err != nil {
		t.Fatalf("decode schema response: %v", err)
	}
	return schema
}

func (c *testAPIClient) ocr(t *testing.T, fields map[string]string, filename string, data []byte) OCRDocumentResponse {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write multipart field: %v", err)
		}
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/ocr", &body)
	if err != nil {
		t.Fatalf("ocr request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(internalAPIHeader, testInternalAPIToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		t.Fatalf("ocr: %v", err)
	}
	defer res.Body.Close()
	requireAPIClientStatus(t, res, http.StatusCreated)

	var doc OCRDocumentResponse
	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		t.Fatalf("decode OCR response: %v", err)
	}
	return doc
}

func requireAPIClientStatus(t *testing.T, res *http.Response, want int) {
	t.Helper()
	if res.StatusCode == want {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("status = %d, want %d; read body: %v", res.StatusCode, want, err)
	}
	t.Fatalf("status = %d, want %d; body=%s", res.StatusCode, want, string(body))
}
