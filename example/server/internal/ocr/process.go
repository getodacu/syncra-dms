package ocr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

const maxMistralResponseBytes int64 = 10 << 20

type Processor func(context.Context, ProcessInput) (*MistralResponse, []byte, error)

type ProcessInput struct {
	Filename string
	MimeType string
	DataURL  string
	Schema   json.RawMessage
	Strict   bool
}

type MistralConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

type MistralResponse struct {
	Pages              []MistralPage  `json:"pages"`
	Model              string         `json:"model"`
	DocumentAnnotation *string        `json:"document_annotation"`
	UsageInfo          map[string]any `json:"usage_info"`
}

type MistralPage struct {
	Index    int             `json:"index"`
	Markdown string          `json:"markdown"`
	Images   json.RawMessage `json:"images,omitempty"`
}

type mistralPageImage struct {
	ID          string `json:"id"`
	ImageBase64 string `json:"image_base64"`
}

type UpstreamError string

func (e UpstreamError) Error() string {
	return string(e)
}

func NewMistralProcessor(cfg MistralConfig) Processor {
	return func(ctx context.Context, input ProcessInput) (*MistralResponse, []byte, error) {
		if cfg.APIKey == "" {
			return nil, nil, errors.New("MISTRAL_API_KEY is required")
		}

		document := map[string]any{"type": "document_url", "document_url": input.DataURL}
		if strings.HasPrefix(input.MimeType, "image/") {
			document = map[string]any{"type": "image_url", "image_url": input.DataURL}
		}

		payload := map[string]any{
			"model":                cfg.Model,
			"document":             document,
			"include_image_base64": true,
		}
		if len(input.Schema) > 0 {
			payload["document_annotation_format"] = map[string]any{
				"type": "json_schema",
				"json_schema": map[string]any{
					"name":   "response_schema",
					"schema": json.RawMessage(input.Schema),
					"strict": input.Strict,
				},
			}
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return nil, nil, err
		}

		baseURL := strings.TrimRight(cfg.BaseURL, "/")
		if baseURL == "" {
			baseURL = "https://api.mistral.ai"
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/ocr", bytes.NewReader(body))
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		req.Header.Set("Content-Type", "application/json")

		client := cfg.Client
		if client == nil {
			client = &http.Client{Timeout: 120 * time.Second}
		}
		res, err := client.Do(req)
		if err != nil {
			return nil, nil, UpstreamError("mistral OCR request failed")
		}
		defer res.Body.Close()

		raw, err := ReadBounded(res.Body, maxMistralResponseBytes)
		if err != nil {
			if errors.Is(err, ErrReaderTooLarge) {
				return nil, nil, UpstreamError("Mistral OCR response too large")
			}
			return nil, nil, err
		}
		if res.StatusCode >= 400 {
			return nil, raw, UpstreamError(fmt.Sprintf("mistral OCR failed with status %d", res.StatusCode))
		}

		var out MistralResponse
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, raw, UpstreamError("invalid Mistral OCR response")
		}
		return &out, raw, nil
	}
}

var ErrReaderTooLarge = errors.New("reader too large")

func ReadBounded(r io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, ErrReaderTooLarge
	}
	return data, nil
}

func DataURL(mimeType string, data []byte) string {
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
}

func JoinMarkdown(pages []MistralPage) string {
	parts := make([]string, 0, len(pages))
	for _, page := range pages {
		markdown := markdownWithPageImages(page.Markdown, parsePageImages(page.Images))
		if markdown != "" {
			parts = append(parts, markdown)
		}
	}
	return strings.Join(parts, "\n\n")
}

func parsePageImages(raw json.RawMessage) []mistralPageImage {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var images []mistralPageImage
	if err := json.Unmarshal(raw, &images); err != nil {
		return nil
	}
	return images
}

func markdownWithPageImages(markdown string, images []mistralPageImage) string {
	for _, image := range images {
		id := strings.TrimSpace(image.ID)
		payload := strings.TrimSpace(image.ImageBase64)
		if id == "" || payload == "" {
			continue
		}

		dataURL := imageDataURL(id, payload)
		if dataURL == "" {
			continue
		}

		imageMarkdown := fmt.Sprintf("![%s](%s)", escapeMarkdownImageAlt(id), dataURL)
		placeholder := fmt.Sprintf("![%s](%s)", id, id)
		if strings.Contains(markdown, placeholder) {
			markdown = strings.ReplaceAll(markdown, placeholder, imageMarkdown)
			continue
		}
		if strings.Contains(markdown, dataURL) {
			continue
		}
		if markdown == "" {
			markdown = imageMarkdown
			continue
		}
		markdown += "\n\n" + imageMarkdown
	}
	return markdown
}

func imageDataURL(id, payload string) string {
	if strings.HasPrefix(payload, "data:") {
		mimeType, base64Payload, ok := parseImageDataURL(payload)
		if !ok {
			return ""
		}
		normalized, ok := normalizeImageBase64(base64Payload)
		if !ok {
			return ""
		}
		return "data:" + mimeType + ";base64," + normalized
	}

	mimeType := inferredImageMIME(id)
	if mimeType == "" {
		return ""
	}
	normalized, ok := normalizeImageBase64(payload)
	if !ok {
		return ""
	}
	return "data:" + mimeType + ";base64," + normalized
}

func inferredImageMIME(id string) string {
	mimeType := mime.TypeByExtension(path.Ext(id))
	if mimeType == "" || !strings.HasPrefix(mimeType, "image/") {
		return "image/png"
	}
	if !isAllowedImageMIME(mimeType) {
		return ""
	}
	return mimeType
}

func parseImageDataURL(payload string) (string, string, bool) {
	withoutPrefix, ok := strings.CutPrefix(payload, "data:")
	if !ok {
		return "", "", false
	}
	mimeType, base64Payload, ok := strings.Cut(withoutPrefix, ";base64,")
	if !ok || mimeType == "" || base64Payload == "" {
		return "", "", false
	}
	if strings.Contains(mimeType, ";") || !isAllowedImageMIME(mimeType) {
		return "", "", false
	}
	return mimeType, base64Payload, true
}

func normalizeImageBase64(payload string) (string, bool) {
	if payload == "" || strings.ContainsAny(payload, "\r\n") {
		return "", false
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", false
	}
	return base64.StdEncoding.EncodeToString(data), true
}

func isAllowedImageMIME(mimeType string) bool {
	switch mimeType {
	case "image/png", "image/jpeg", "image/gif", "image/webp", "image/avif":
		return true
	default:
		return false
	}
}

func escapeMarkdownImageAlt(alt string) string {
	var escaped strings.Builder
	for _, r := range alt {
		switch r {
		case '\\', '[', ']', '(', ')':
			escaped.WriteByte('\\')
			escaped.WriteRune(r)
		case '\r':
			escaped.WriteString(`\r`)
		case '\n':
			escaped.WriteString(`\n`)
		default:
			escaped.WriteRune(r)
		}
	}
	return escaped.String()
}

func CountRawResponsePages(raw []byte) (int, error) {
	var response struct {
		Pages json.RawMessage `json:"pages"`
	}
	if err := json.Unmarshal(raw, &response); err != nil {
		return 0, fmt.Errorf("invalid OCR raw response JSON: %w", err)
	}
	if len(response.Pages) == 0 || string(response.Pages) == "null" {
		return 0, nil
	}

	var pages []json.RawMessage
	if err := json.Unmarshal(response.Pages, &pages); err != nil {
		return 0, fmt.Errorf("invalid OCR raw response pages: %w", err)
	}
	return len(pages), nil
}

func ParseAnnotationJSON(annotation *string, required bool) (json.RawMessage, error) {
	if annotation == nil || strings.TrimSpace(*annotation) == "" {
		if required {
			return nil, errors.New("missing Mistral document annotation JSON")
		}
		return nil, nil
	}
	raw := json.RawMessage(*annotation)
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return raw, nil
}
