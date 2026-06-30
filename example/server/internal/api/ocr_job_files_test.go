package api

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestWriteOCRJobFileUsesJobIDAndDetectedExtension(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	data := validPNGBytes()

	path, err := writeOCRJobFile(dir, id, "image/png", data)
	if err != nil {
		t.Fatalf("write OCR job file: %v", err)
	}

	wantPath := filepath.Join(dir, id.String()+".png")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}
	if strings.Contains(path, "scan.png") {
		t.Fatalf("path contains uploaded filename: %q", path)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read stored file: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("stored file = %v, want %v", got, data)
	}
}

func TestWriteOCRJobFileRejectsUnsupportedMIME(t *testing.T) {
	_, err := writeOCRJobFile(t.TempDir(), uuid.New(), "text/plain", []byte("hello"))
	if err == nil {
		t.Fatal("err = nil")
	}
}

func TestWriteOCRJobFileRejectsNilJobID(t *testing.T) {
	dir := t.TempDir()

	_, err := writeOCRJobFile(dir, uuid.Nil, "image/png", validPNGBytes())
	if err == nil {
		t.Fatal("err = nil")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("stored files = %d, want 0", len(entries))
	}
}

func TestWriteOCRJobFileRejectsDuplicateJobIDWithoutReplacingFile(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	original := validPNGBytes()
	replacement := append([]byte{}, original...)
	replacement = append(replacement, 1, 2, 3)

	path, err := writeOCRJobFile(dir, id, "image/png", original)
	if err != nil {
		t.Fatalf("write original OCR job file: %v", err)
	}
	if _, err := writeOCRJobFile(dir, id, "image/png", replacement); err == nil {
		t.Fatal("duplicate write err = nil")
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read stored file: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("stored file changed to %v, want %v", got, original)
	}
	assertDirEntryNames(t, dir, id.String()+".claim", id.String()+".png")
}

func TestWriteOCRJobFileUsesSupportedMIMEExtensions(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		data     []byte
		ext      string
	}{
		{
			name:     "PDF",
			mimeType: "application/pdf",
			data:     validPDFBytes(),
			ext:      ".pdf",
		},
		{
			name:     "JPEG",
			mimeType: "image/jpeg",
			data:     []byte{0xff, 0xd8, 0xff, 0xdb},
			ext:      ".jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			id := uuid.New()

			path, err := writeOCRJobFile(dir, id, tt.mimeType, tt.data)
			if err != nil {
				t.Fatalf("write OCR job file: %v", err)
			}

			wantPath := filepath.Join(dir, id.String()+tt.ext)
			if path != wantPath {
				t.Fatalf("path = %q, want %q", path, wantPath)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read stored file: %v", err)
			}
			if !bytes.Equal(got, tt.data) {
				t.Fatalf("stored file = %v, want %v", got, tt.data)
			}
		})
	}
}

func TestWriteOCRJobFileRejectsSameJobIDAcrossMIMETypes(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	original := validPNGBytes()

	path, err := writeOCRJobFile(dir, id, "image/png", original)
	if err != nil {
		t.Fatalf("write original OCR job file: %v", err)
	}

	for _, tt := range []struct {
		name     string
		mimeType string
		data     []byte
	}{
		{
			name:     "PDF",
			mimeType: "application/pdf",
			data:     validPDFBytes(),
		},
		{
			name:     "JPEG",
			mimeType: "image/jpeg",
			data:     []byte{0xff, 0xd8, 0xff, 0xdb},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := writeOCRJobFile(dir, id, tt.mimeType, tt.data); err == nil {
				t.Fatal("duplicate write err = nil")
			}
		})
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read stored file: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Fatalf("stored file changed to %v, want %v", got, original)
	}
	assertDirEntryNames(t, dir, id.String()+".claim", id.String()+".png")
}

func TestWriteOCRJobFileRejectsExistingLegacyFinalFileAcrossExtensions(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	legacyPath := filepath.Join(dir, id.String()+".pdf")
	if err := os.WriteFile(legacyPath, validPDFBytes(), 0o600); err != nil {
		t.Fatalf("write legacy OCR job file: %v", err)
	}

	if _, err := writeOCRJobFile(dir, id, "image/png", validPNGBytes()); err == nil {
		t.Fatal("duplicate write err = nil")
	}

	assertDirEntryNames(t, dir, id.String()+".pdf")
}

func TestWriteOCRJobFileRejectsClaimOnlyLeftoverWithoutFinalFile(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	claimPath := ocrJobClaimPath(dir, id)
	claimBytes := []byte("orphan claim")
	if err := os.WriteFile(claimPath, claimBytes, 0o600); err != nil {
		t.Fatalf("write orphan claim: %v", err)
	}
	oldTime := time.Now().Add(-10 * time.Minute)
	if err := os.Chtimes(claimPath, oldTime, oldTime); err != nil {
		t.Fatalf("age orphan claim: %v", err)
	}

	if _, err := writeOCRJobFile(dir, id, "image/png", validPNGBytes()); err == nil {
		t.Fatal("err = nil")
	}
	gotClaim, err := os.ReadFile(claimPath)
	if err != nil {
		t.Fatalf("read orphan claim: %v", err)
	}
	if !bytes.Equal(gotClaim, claimBytes) {
		t.Fatalf("claim = %q, want %q", gotClaim, claimBytes)
	}
	assertDirEntryNames(t, dir, id.String()+".claim")
}

func TestWriteOCRJobFileConcurrentSameIDAllowsOneWinner(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	type writeCase struct {
		mimeType string
		data     []byte
	}
	cases := []writeCase{
		{mimeType: "image/png", data: validPNGBytes()},
		{mimeType: "application/pdf", data: validPDFBytes()},
	}
	type result struct {
		path string
		data []byte
		err  error
	}
	start := make(chan struct{})
	results := make(chan result, len(cases))
	var wg sync.WaitGroup
	for _, tc := range cases {
		tc := tc
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			path, err := writeOCRJobFile(dir, id, tc.mimeType, tc.data)
			results <- result{path: path, data: tc.data, err: err}
		}()
	}
	close(start)
	wg.Wait()
	close(results)

	var successes []result
	var failures []error
	for result := range results {
		if result.err != nil {
			failures = append(failures, result.err)
			continue
		}
		successes = append(successes, result)
	}
	if len(successes) != 1 || len(failures) != 1 {
		t.Fatalf("successes = %d failures = %d, want 1 success and 1 failure", len(successes), len(failures))
	}
	got, err := os.ReadFile(successes[0].path)
	if err != nil {
		t.Fatalf("read winning file: %v", err)
	}
	if !bytes.Equal(got, successes[0].data) {
		t.Fatalf("winning file = %v, want %v", got, successes[0].data)
	}
	assertOneFinalFileForJobID(t, dir, id)
}

func TestWriteOCRJobFileCleansUpClaimAfterClaimDirSyncFailure(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	syncErr := errors.New("claim directory sync failed")
	cleanupErr := errors.New("cleanup sync failed")
	oldSyncPath := syncPathFunc
	dirSyncs := 0
	syncPathFunc = func(path string) error {
		if path == dir {
			dirSyncs++
			if dirSyncs == 1 {
				return syncErr
			}
			return cleanupErr
		}
		return oldSyncPath(path)
	}
	t.Cleanup(func() {
		syncPathFunc = oldSyncPath
	})

	_, err := writeOCRJobFile(dir, id, "image/png", validPNGBytes())
	if !errors.Is(err, syncErr) {
		t.Fatalf("error = %v, want %v", err, syncErr)
	}
	if !errors.Is(err, cleanupErr) {
		t.Fatalf("error = %v, want cleanup error %v", err, cleanupErr)
	}
	assertNoJobIDEntries(t, dir, id)
}

func TestWriteOCRJobFileCleansUpAfterPostLinkSyncFailure(t *testing.T) {
	dir := t.TempDir()
	id := uuid.New()
	syncErr := errors.New("sync failed")
	cleanupErr := errors.New("cleanup sync failed")
	oldSyncPath := syncPathFunc
	dirSyncs := 0
	syncPathFunc = func(path string) error {
		if path == dir {
			dirSyncs++
			if dirSyncs == 2 {
				return syncErr
			}
			if dirSyncs > 2 {
				return cleanupErr
			}
		}
		return oldSyncPath(path)
	}
	t.Cleanup(func() {
		syncPathFunc = oldSyncPath
	})

	_, err := writeOCRJobFile(dir, id, "image/png", validPNGBytes())
	if !errors.Is(err, syncErr) {
		t.Fatalf("error = %v, want %v", err, syncErr)
	}
	if !errors.Is(err, cleanupErr) {
		t.Fatalf("error = %v, want cleanup error %v", err, cleanupErr)
	}
	assertNoJobIDEntries(t, dir, id)
}

func TestWriteOCRJobFileCreatesNestedStorageDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), ".data", "ocr-files")
	id := uuid.New()

	path, err := writeOCRJobFile(dir, id, "image/png", validPNGBytes())
	if err != nil {
		t.Fatalf("write OCR job file: %v", err)
	}

	wantPath := filepath.Join(dir, id.String()+".png")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat stored file: %v", err)
	}
}

func TestWriteOCRJobFileReturnsAbsolutePathForRelativeStorageDir(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	dir := filepath.Join("relative", "ocr-jobs")
	id := uuid.New()

	path, err := writeOCRJobFile(dir, id, "image/png", validPNGBytes())
	if err != nil {
		t.Fatalf("write OCR job file: %v", err)
	}

	wantPath := filepath.Join(root, dir, id.String()+".png")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}
	if !filepath.IsAbs(path) {
		t.Fatalf("path = %q, want absolute path", path)
	}
}

func TestDefaultStorageDirsFromServerRoot(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	serverRoot := filepath.Clean(filepath.Join(cwd, "..", ".."))
	t.Chdir(serverRoot)

	h := &Handler{}
	storageRoot, err := h.storageRootDir()
	if err != nil {
		t.Fatalf("default storage root dir: %v", err)
	}
	if want := filepath.Join(serverRoot, ".data"); storageRoot != want {
		t.Fatalf("storage root = %q, want %q", storageRoot, want)
	}
	ocrDir, err := h.ocrJobFileDir()
	if err != nil {
		t.Fatalf("default OCR job file dir: %v", err)
	}
	if want := filepath.Join(serverRoot, ".data", "ocr-files"); ocrDir != want {
		t.Fatalf("OCR dir = %q, want %q", ocrDir, want)
	}
	invoiceDir, err := h.billingInvoicePDFDir()
	if err != nil {
		t.Fatalf("default billing invoice PDF dir: %v", err)
	}
	if want := filepath.Join(serverRoot, ".data", "invoices"); invoiceDir != want {
		t.Fatalf("invoice dir = %q, want %q", invoiceDir, want)
	}
	if !strings.HasSuffix(ocrDir, filepath.Join("server", ".data", "ocr-files")) {
		t.Fatalf("OCR dir = %q, want suffix %q", ocrDir, filepath.Join("server", ".data", "ocr-files"))
	}
}

func TestRelativeStorageDirFromServerRoot(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	serverRoot := filepath.Clean(filepath.Join(cwd, "..", ".."))
	t.Chdir(filepath.Dir(serverRoot))

	h := &Handler{StorageDir: filepath.Join("custom", "storage")}
	storageRoot, err := h.storageRootDir()
	if err != nil {
		t.Fatalf("relative storage root dir: %v", err)
	}
	if want := filepath.Join(serverRoot, "custom", "storage"); storageRoot != want {
		t.Fatalf("storage root = %q, want %q", storageRoot, want)
	}
	ocrDir, err := h.ocrJobFileDir()
	if err != nil {
		t.Fatalf("relative OCR job file dir: %v", err)
	}
	if want := filepath.Join(serverRoot, "custom", "storage", "ocr-files"); ocrDir != want {
		t.Fatalf("OCR dir = %q, want %q", ocrDir, want)
	}
}

func TestAbsoluteStorageDir(t *testing.T) {
	storageDir := t.TempDir()
	t.Chdir(t.TempDir())

	h := &Handler{StorageDir: storageDir}

	storageRoot, err := h.storageRootDir()
	if err != nil {
		t.Fatalf("absolute storage root dir: %v", err)
	}
	if storageRoot != filepath.Clean(storageDir) {
		t.Fatalf("storage root = %q, want %q", storageRoot, filepath.Clean(storageDir))
	}
	ocrDir, err := h.ocrJobFileDir()
	if err != nil {
		t.Fatalf("absolute OCR job file dir: %v", err)
	}
	if want := filepath.Join(storageDir, "ocr-files"); ocrDir != want {
		t.Fatalf("OCR dir = %q, want %q", ocrDir, want)
	}
	invoiceDir, err := h.billingInvoicePDFDir()
	if err != nil {
		t.Fatalf("absolute billing invoice PDF dir: %v", err)
	}
	if want := filepath.Join(storageDir, "invoices"); invoiceDir != want {
		t.Fatalf("invoice dir = %q, want %q", invoiceDir, want)
	}
}

func TestEnsureDirDurableCreatesNestedMissingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), ".data", "ocr-files")

	if err := ensureDirDurable(dir); err != nil {
		t.Fatalf("ensure durable dir: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat dir: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("path is not a directory: %q", dir)
	}
}

func TestEnsureDirDurableAcceptsExistingDirectory(t *testing.T) {
	dir := t.TempDir()

	if err := ensureDirDurable(dir); err != nil {
		t.Fatalf("ensure durable dir: %v", err)
	}
}

func TestEnsureDirDurableRejectsExistingFileInPath(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, ".data")
	if err := os.WriteFile(filePath, []byte("not a dir"), 0o600); err != nil {
		t.Fatalf("write file path component: %v", err)
	}

	err := ensureDirDurable(filepath.Join(filePath, "ocr-files"))
	if err == nil {
		t.Fatal("err = nil")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("error = %q, want not a directory", err.Error())
	}
}

func TestEnsureDirDurableRejectsExistingFileAtTarget(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "ocr-jobs")
	if err := os.WriteFile(filePath, []byte("not a dir"), 0o600); err != nil {
		t.Fatalf("write target file: %v", err)
	}

	err := ensureDirDurable(filePath)
	if err == nil {
		t.Fatal("err = nil")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("error = %q, want not a directory", err.Error())
	}
}

func TestSyncPathSyncsExistingDirectory(t *testing.T) {
	if err := syncPath(t.TempDir()); err != nil {
		t.Fatalf("sync path: %v", err)
	}
}

func TestSyncPathReturnsErrorForMissingPath(t *testing.T) {
	err := syncPath(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("err = nil")
	}
}

func entryNames(entries []os.DirEntry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func assertDirEntryNames(t *testing.T, dir string, want ...string) {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	got := entryNames(entries)
	if len(got) != len(want) {
		t.Fatalf("stored entries = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("stored entries = %v, want %v", got, want)
		}
	}
}

func assertNoJobIDEntries(t *testing.T, dir string, id uuid.UUID) {
	t.Helper()

	prefix := id.String()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prefix) {
			t.Fatalf("unexpected job entry after failed write: %q", entry.Name())
		}
	}
}

func assertOneFinalFileForJobID(t *testing.T, dir string, id uuid.UUID) {
	t.Helper()

	prefix := id.String()
	finalFiles := 0
	claimFiles := 0
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		switch {
		case name == prefix+".claim":
			claimFiles++
		case strings.HasSuffix(name, ".pdf") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg"):
			finalFiles++
		default:
			t.Fatalf("unexpected job entry: %q", name)
		}
	}
	if claimFiles != 1 || finalFiles != 1 {
		t.Fatalf("claim files = %d final files = %d, want 1 claim and 1 final", claimFiles, finalFiles)
	}
}
