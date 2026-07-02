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

	"github.com/google/uuid"
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

func TestLocalStorageAllowsConfiguredSniffedMIMEType(t *testing.T) {
	store := NewLocalStorage(t.TempDir(), 1024, []string{"text/plain; charset=utf-8"})

	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); err != nil {
		t.Fatalf("Save error = %v", err)
	}
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

func TestLocalStorageOpenRejectsNonCanonicalStorageKeys(t *testing.T) {
	store := NewLocalStorage(t.TempDir(), 1024, nil)
	validID := uuid.NewString()
	upperID := strings.ToUpper(uuid.NewString())

	for _, key := range []string{
		"",
		"tmp/file",
		"documents/../tmp/file",
		"documents/./xx/file",
		`documents\ab\` + validID + `.bin`,
		"documents/ab",
		"documents/ab/",
		"documents/a/" + validID + ".bin",
		"documents/abc/" + validID + ".bin",
		"documents/zz/" + validID + ".bin",
		"documents/ab/not-a-uuid.bin",
		"documents/ab/" + upperID + ".bin",
		"documents/ab/" + validID,
		"documents/ab/" + validID + ".txt",
		"documents/ab/" + validID + ".bin/extra",
	} {
		t.Run(key, func(t *testing.T) {
			assertInvalidOpenKey(t, store, key)
		})
	}
}

func TestLocalStorageOpenRejectsRootContainedInternalFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "tmp"), 0o700); err != nil {
		t.Fatalf("create temp directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "tmp", "some-file"), []byte("secret"), 0o600); err != nil {
		t.Fatalf("write internal file: %v", err)
	}
	store := NewLocalStorage(root, 1024, nil)

	assertInvalidOpenKey(t, store, "tmp/some-file")
	assertInvalidOpenKey(t, store, "documents/../tmp/some-file")
}

func TestLocalStorageOpenRejectsSymlinkPathComponents(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	id := uuid.NewString()
	if err := os.MkdirAll(filepath.Join(root, "documents"), 0o700); err != nil {
		t.Fatalf("create documents directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outside, id+".bin"), []byte("outside"), 0o600); err != nil {
		t.Fatalf("write outside file: %v", err)
	}
	createSymlink(t, outside, filepath.Join(root, "documents", "ab"))
	store := NewLocalStorage(root, 1024, nil)

	assertInvalidOpenKey(t, store, "documents/ab/"+id+".bin")
}

func TestLocalStorageOpenRejectsSymlinkFinalFile(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	id := uuid.NewString()
	prefixDir := filepath.Join(root, "documents", "ab")
	if err := os.MkdirAll(prefixDir, 0o700); err != nil {
		t.Fatalf("create prefix directory: %v", err)
	}
	outsideFile := filepath.Join(outside, "outside.bin")
	if err := os.WriteFile(outsideFile, []byte("outside"), 0o600); err != nil {
		t.Fatalf("write outside file: %v", err)
	}
	createSymlink(t, outsideFile, filepath.Join(prefixDir, id+".bin"))
	store := NewLocalStorage(root, 1024, nil)

	assertInvalidOpenKey(t, store, "documents/ab/"+id+".bin")
}

func TestLocalStorageSaveRejectsSymlinkPrefixDirectory(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "documents"), 0o700); err != nil {
		t.Fatalf("create documents directory: %v", err)
	}
	createSymlink(t, outside, filepath.Join(root, "documents", "2c"))
	store := NewLocalStorage(root, 1024, nil)

	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); err == nil {
		t.Fatal("Save error = nil, want symlink prefix rejection")
	}
	assertNoTempFiles(t, root)
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatalf("read outside directory: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("outside directory entries = %d, want 0", len(entries))
	}
}

func TestLocalStorageSaveRejectsSymlinkDocumentsDirectory(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	createSymlink(t, outside, filepath.Join(root, "documents"))
	store := NewLocalStorage(root, 1024, nil)

	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); err == nil {
		t.Fatal("Save error = nil, want symlink documents rejection")
	}
	assertNoTempFiles(t, root)
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatalf("read outside directory: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("outside directory entries = %d, want 0", len(entries))
	}
}

func TestLocalStorageSaveRejectsSymlinkTempDirectory(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	createSymlink(t, outside, filepath.Join(root, "tmp"))
	store := NewLocalStorage(root, 1024, nil)

	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); err == nil {
		t.Fatal("Save error = nil, want symlink temp rejection")
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatalf("read outside directory: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("outside directory entries = %d, want 0", len(entries))
	}
}

func TestLocalStorageSaveMakesManagedDirectoriesPrivate(t *testing.T) {
	root := t.TempDir()
	tmpDir := filepath.Join(root, "tmp")
	documentsDir := filepath.Join(root, "documents")
	prefixDir := filepath.Join(documentsDir, "2c")
	for _, dir := range []string{tmpDir, prefixDir} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatalf("create directory %s: %v", dir, err)
		}
	}
	for _, dir := range []string{tmpDir, documentsDir, prefixDir} {
		chmodForPermissionTest(t, dir, 0o777)
	}

	store := NewLocalStorage(root, 1024, nil)
	if _, err := store.Save(context.Background(), strings.NewReader("hello"), "invoice.txt"); err != nil {
		t.Fatalf("Save error = %v", err)
	}

	for _, dir := range []string{tmpDir, documentsDir, prefixDir} {
		assertPrivateDirectoryMode(t, dir)
	}
}

func assertInvalidOpenKey(t *testing.T, store *LocalStorage, key string) {
	t.Helper()
	reader, err := store.Open(key)
	if !errors.Is(err, ErrInvalidStorageKey) {
		if reader != nil {
			reader.Close()
		}
		t.Fatalf("Open(%q) error = %v, want ErrInvalidStorageKey", key, err)
	}
	if reader != nil {
		reader.Close()
		t.Fatalf("Open(%q) reader = %v, want nil reader", key, reader)
	}
}

func createSymlink(t *testing.T, target string, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
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

func chmodForPermissionTest(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	if err := os.Chmod(path, mode); err != nil {
		t.Skipf("chmod unavailable for permission test: %v", err)
	}
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("inspect %s: %v", path, err)
	}
	if info.Mode().Perm()&0o777 != mode {
		t.Skipf("filesystem does not preserve chmod mode %o on %s; got %o", mode, path, info.Mode().Perm())
	}
}

func assertPrivateDirectoryMode(t *testing.T, path string) {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("inspect %s: %v", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		t.Fatalf("%s mode = %v, want real directory", path, info.Mode())
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("%s mode = %o, want no group/other permissions", path, info.Mode().Perm())
	}
}
