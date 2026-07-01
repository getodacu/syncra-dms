package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStorageStoresAndOpensDocument(t *testing.T) {
	root := t.TempDir()
	store := NewLocalStorage(root, 1024, nil)

	result, err := store.Save(context.Background(), strings.NewReader("hello"), "../invoice.pdf")
	if err != nil {
		t.Fatalf("Save error = %v", err)
	}

	wantHash := sha256.Sum256([]byte("hello"))
	if result.SizeBytes != 5 || result.SHA256Hash != hex.EncodeToString(wantHash[:]) || result.StorageKey == "" {
		t.Fatalf("result = %#v", result)
	}
	if result.OriginalFileName != "invoice.pdf" {
		t.Fatalf("OriginalFileName = %q, want invoice.pdf", result.OriginalFileName)
	}
	if result.MimeType != "text/plain; charset=utf-8" {
		t.Fatalf("MimeType = %q, want text/plain; charset=utf-8", result.MimeType)
	}
	if result.Extension == nil || *result.Extension != ".pdf" {
		t.Fatalf("Extension = %v, want .pdf", result.Extension)
	}
	if filepath.IsAbs(result.StorageKey) {
		t.Fatalf("StorageKey = %q, want relative key", result.StorageKey)
	}
	if !strings.HasPrefix(result.StorageKey, "documents/") || !strings.HasSuffix(result.StorageKey, ".bin") {
		t.Fatalf("StorageKey = %q, want documents/.../*.bin", result.StorageKey)
	}
	if strings.Contains(result.StorageKey, string(os.PathSeparator)) && os.PathSeparator != '/' {
		t.Fatalf("StorageKey = %q, want forward slash separators", result.StorageKey)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(result.StorageKey))); err != nil {
		t.Fatalf("stored file stat error = %v", err)
	}

	reader, err := store.Open(result.StorageKey)
	if err != nil {
		t.Fatalf("Open error = %v", err)
	}
	defer reader.Close()
	body, _ := io.ReadAll(reader)
	if string(body) != "hello" {
		t.Fatalf("body = %q", body)
	}
}

func TestLocalStorageRejectsUploadsOverMaxSize(t *testing.T) {
	root := t.TempDir()
	store := NewLocalStorage(root, 4, nil)

	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); !errors.Is(err, ErrUploadTooLarge) {
		t.Fatalf("Save error = %v, want ErrUploadTooLarge", err)
	}
	assertNoTempFiles(t, root)
}

func TestLocalStorageRejectsDisallowedMIMEType(t *testing.T) {
	root := t.TempDir()
	store := NewLocalStorage(root, 1024, []string{"application/pdf"})

	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); !errors.Is(err, ErrUnsupportedMIME) {
		t.Fatalf("Save error = %v, want ErrUnsupportedMIME", err)
	}
	assertNoTempFiles(t, root)
}

func TestLocalStorageOpenRejectsTraversal(t *testing.T) {
	store := NewLocalStorage(t.TempDir(), 1024, nil)

	for _, key := range []string{"../outside.txt", "/outside.txt", "documents/../../outside.txt"} {
		t.Run(key, func(t *testing.T) {
			if reader, err := store.Open(key); !errors.Is(err, ErrInvalidStorageKey) {
				if reader != nil {
					reader.Close()
				}
				t.Fatalf("Open(%q) error = %v, want ErrInvalidStorageKey", key, err)
			} else if reader != nil {
				reader.Close()
				t.Fatalf("Open(%q) reader = %v, want nil reader", key, reader)
			}
		})
	}
}

func assertNoTempFiles(t *testing.T, root string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, "tmp"))
	if err != nil {
		t.Fatalf("read temp directory: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("temp directory entries = %d, want 0", len(entries))
	}
}
