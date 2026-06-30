package ocr

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestDataURL(t *testing.T) {
	got := DataURL("image/png", []byte("png-data"))
	want := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte("png-data"))
	if got != want {
		t.Fatalf("DataURL() = %q, want %q", got, want)
	}
}

func TestJoinMarkdown(t *testing.T) {
	got := JoinMarkdown([]MistralPage{
		{Index: 0, Markdown: "# First"},
		{Index: 1},
		{Index: 2, Markdown: "Second"},
	})
	if got != "# First\n\nSecond" {
		t.Fatalf("JoinMarkdown() = %q", got)
	}
}

func TestJoinMarkdownReplacesMistralImagePlaceholders(t *testing.T) {
	images := json.RawMessage(`[{"id":"img-0.jpeg","image_base64":"data:image/jpeg;base64,ZmFrZQ=="}]`)
	got := JoinMarkdown([]MistralPage{
		{
			Index:    0,
			Markdown: "Before\n\n![img-0.jpeg](img-0.jpeg)\n\nAfter",
			Images:   images,
		},
	})
	want := "Before\n\n![img-0.jpeg](data:image/jpeg;base64,ZmFrZQ==)\n\nAfter"
	if got != want {
		t.Fatalf("JoinMarkdown() = %q, want %q", got, want)
	}
}

func TestJoinMarkdownAppendsUnreferencedImages(t *testing.T) {
	images := json.RawMessage(`[{"id":"chart.png","image_base64":"iVBORw0KGgo="}]`)
	got := JoinMarkdown([]MistralPage{
		{
			Index:    0,
			Markdown: "# Chart",
			Images:   images,
		},
	})
	want := "# Chart\n\n![chart.png](data:image/png;base64,iVBORw0KGgo=)"
	if got != want {
		t.Fatalf("JoinMarkdown() = %q, want %q", got, want)
	}
}

func TestJoinMarkdownIgnoresUnsupportedImages(t *testing.T) {
	got := JoinMarkdown([]MistralPage{
		{
			Index:    0,
			Markdown: "# Invoice",
			Images:   json.RawMessage(`{"not":"an array"}`),
		},
		{
			Index:    1,
			Markdown: "Total",
			Images:   json.RawMessage(`[{"id":"missing-payload.png"}]`),
		},
	})
	want := "# Invoice\n\nTotal"
	if got != want {
		t.Fatalf("JoinMarkdown() = %q, want %q", got, want)
	}
}

func TestJoinMarkdownEscapesImageAltText(t *testing.T) {
	injectedID := "bad](javascript:alert(1))\n![owned"
	raw, err := json.Marshal([]mistralPageImage{
		{ID: injectedID, ImageBase64: "data:image/png;base64,c2FmZQ=="},
	})
	if err != nil {
		t.Fatalf("marshal images: %v", err)
	}

	got := JoinMarkdown([]MistralPage{
		{
			Index:    0,
			Markdown: "# Image",
			Images:   raw,
		},
	})

	if strings.Contains(got, "](javascript:alert(1))") || strings.Contains(got, "\n![owned") {
		t.Fatalf("JoinMarkdown() contains raw injected markdown: %q", got)
	}
	want := "# Image\n\n![bad\\]\\(javascript:alert\\(1\\)\\)\\n!\\[owned](data:image/png;base64,c2FmZQ==)"
	if got != want {
		t.Fatalf("JoinMarkdown() = %q, want %q", got, want)
	}
}

func TestJoinMarkdownIgnoresMalformedDataImagePayloads(t *testing.T) {
	got := JoinMarkdown([]MistralPage{
		{
			Index:    0,
			Markdown: "# Image",
			Images: json.RawMessage(`[
				{"id":"bad-base64.png","image_base64":"data:image/png;base64,not-base64"},
				{"id":"unsupported.svg","image_base64":"data:image/svg+xml;base64,PHN2Zz48L3N2Zz4="}
			]`),
		},
	})
	want := "# Image"
	if got != want {
		t.Fatalf("JoinMarkdown() = %q, want %q", got, want)
	}
}

func TestJoinMarkdownIgnoresMalformedRawImagePayloads(t *testing.T) {
	got := JoinMarkdown([]MistralPage{
		{
			Index:    0,
			Markdown: "# Image",
			Images: json.RawMessage(`[
				{"id":"bad-close.png","image_base64":"abc)def"},
				{"id":"bad-newline.png","image_base64":"YWJj\ndef"},
				{"id":"unsupported.svg","image_base64":"PHN2Zz48L3N2Zz4="}
			]`),
		},
	})
	want := "# Image"
	if got != want {
		t.Fatalf("JoinMarkdown() = %q, want %q", got, want)
	}
}

func TestCountRawResponsePages(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		want    int
		wantErr bool
	}{
		{
			name: "array pages",
			raw:  []byte(`{"pages":[{"index":0},{"index":1}],"model":"mistral-ocr-latest"}`),
			want: 2,
		},
		{
			name: "missing pages",
			raw:  []byte(`{"model":"mistral-ocr-latest"}`),
			want: 0,
		},
		{
			name: "null pages",
			raw:  []byte(`{"pages":null}`),
			want: 0,
		},
		{
			name:    "invalid JSON",
			raw:     []byte(`{"pages":[`),
			wantErr: true,
		},
		{
			name:    "non-array pages",
			raw:     []byte(`{"pages":{"index":0}}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountRawResponsePages(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatal("CountRawResponsePages() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("CountRawResponsePages() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("CountRawResponsePages() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestParseAnnotationJSON(t *testing.T) {
	annotation := `{"total":42}`
	got, err := ParseAnnotationJSON(&annotation, true)
	if err != nil {
		t.Fatalf("ParseAnnotationJSON() error = %v", err)
	}
	assertRawJSONEqual(t, got, annotation)
}

func TestParseAnnotationJSONAllowsMissingWhenOptional(t *testing.T) {
	got, err := ParseAnnotationJSON(nil, false)
	if err != nil {
		t.Fatalf("ParseAnnotationJSON() error = %v", err)
	}
	if got != nil {
		t.Fatalf("annotation = %s, want nil", string(got))
	}
}

func TestParseAnnotationJSONRequiresAnnotation(t *testing.T) {
	_, err := ParseAnnotationJSON(nil, true)
	if err == nil || !strings.Contains(err.Error(), "missing Mistral document annotation JSON") {
		t.Fatalf("error = %v, want missing annotation", err)
	}
}

func assertRawJSONEqual(t *testing.T, got json.RawMessage, want string) {
	t.Helper()
	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("decode got JSON: %v", err)
	}
	var wantValue any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("decode want JSON: %v", err)
	}
	if gotJSON, wantJSON := string(got), want; gotJSON != wantJSON {
		t.Fatalf("json = %s, want %s", gotJSON, wantJSON)
	}
}
