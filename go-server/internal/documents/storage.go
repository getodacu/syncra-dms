package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrUploadTooLarge    = errors.New("document exceeds maximum upload size")
	ErrUnsupportedMIME   = errors.New("document MIME type is not allowed")
	ErrInvalidStorageKey = errors.New("invalid storage key")
)

type LocalStorage struct {
	root             string
	maxUploadBytes   int64
	allowedMIMETypes map[string]bool
}

type StoredFile struct {
	OriginalFileName string
	MimeType         string
	Extension        *string
	SizeBytes        int64
	SHA256Hash       string
	StorageKey       string
}

func NewLocalStorage(root string, maxUploadBytes int64, allowedMIMETypes []string) *LocalStorage {
	root = strings.TrimSpace(root)
	if root == "" {
		root = "."
	}
	if absRoot, err := filepath.Abs(root); err == nil {
		root = absRoot
	} else {
		root = filepath.Clean(root)
	}

	var allowed map[string]bool
	for _, mimeType := range allowedMIMETypes {
		mimeType = strings.TrimSpace(mimeType)
		if mimeType == "" {
			continue
		}
		if allowed == nil {
			allowed = make(map[string]bool, len(allowedMIMETypes))
		}
		allowed[mimeType] = true
	}

	return &LocalStorage{
		root:             root,
		maxUploadBytes:   maxUploadBytes,
		allowedMIMETypes: allowed,
	}
}

func (s *LocalStorage) Save(ctx context.Context, reader io.Reader, originalFileName string) (StoredFile, error) {
	if err := os.MkdirAll(filepath.Join(s.root, "tmp"), 0o700); err != nil {
		return StoredFile{}, fmt.Errorf("create storage temp directory: %w", err)
	}

	tmpFile, err := os.CreateTemp(filepath.Join(s.root, "tmp"), "upload-*")
	if err != nil {
		return StoredFile{}, fmt.Errorf("create storage temp file: %w", err)
	}

	tmpPath := tmpFile.Name()
	keepTempFile := false
	closed := false
	closeTemp := func() error {
		if closed {
			return nil
		}
		closed = true
		return tmpFile.Close()
	}
	defer func() {
		_ = closeTemp()
		if !keepTempFile {
			_ = os.Remove(tmpPath)
		}
	}()

	hasher := sha256.New()
	buffer := make([]byte, 32*1024)
	sniffBytes := make([]byte, 0, 512)
	var sizeBytes int64

	for {
		if err := ctx.Err(); err != nil {
			return StoredFile{}, err
		}

		n, readErr := reader.Read(buffer)
		if n > 0 {
			chunk := buffer[:n]
			sizeBytes += int64(n)
			if s.maxUploadBytes > 0 && sizeBytes > s.maxUploadBytes {
				return StoredFile{}, fmt.Errorf("%w: limit is %d bytes", ErrUploadTooLarge, s.maxUploadBytes)
			}

			if remaining := 512 - len(sniffBytes); remaining > 0 {
				if n < remaining {
					remaining = n
				}
				sniffBytes = append(sniffBytes, chunk[:remaining]...)
			}
			_, _ = hasher.Write(chunk)
			if _, err := tmpFile.Write(chunk); err != nil {
				return StoredFile{}, fmt.Errorf("write storage temp file: %w", err)
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return StoredFile{}, fmt.Errorf("read upload: %w", readErr)
		}
	}

	mimeType := http.DetectContentType(sniffBytes)
	if len(s.allowedMIMETypes) > 0 && !s.allowedMIMETypes[mimeType] {
		return StoredFile{}, fmt.Errorf("%w: %s", ErrUnsupportedMIME, mimeType)
	}

	sha256Hash := hex.EncodeToString(hasher.Sum(nil))
	hashPrefix := sha256Hash[:2]
	documentsDir := filepath.Join(s.root, "documents")
	if err := os.MkdirAll(documentsDir, 0o700); err != nil {
		return StoredFile{}, fmt.Errorf("create storage documents directory: %w", err)
	}
	if err := ensureRealDirectory(documentsDir); err != nil {
		return StoredFile{}, err
	}

	targetDir := filepath.Join(documentsDir, hashPrefix)
	if err := os.MkdirAll(targetDir, 0o700); err != nil {
		return StoredFile{}, fmt.Errorf("create storage document directory: %w", err)
	}
	if err := ensureRealDirectory(targetDir); err != nil {
		return StoredFile{}, err
	}

	targetName := uuid.NewString() + ".bin"
	targetPath := filepath.Join(targetDir, targetName)
	storageKey := "documents/" + hashPrefix + "/" + targetName

	if err := closeTemp(); err != nil {
		return StoredFile{}, fmt.Errorf("close storage temp file: %w", err)
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		return StoredFile{}, fmt.Errorf("store document file: %w", err)
	}
	keepTempFile = true

	safeName := SafeOriginalFileName(originalFileName)
	return StoredFile{
		OriginalFileName: safeName,
		MimeType:         mimeType,
		Extension:        fileExtension(safeName),
		SizeBytes:        sizeBytes,
		SHA256Hash:       sha256Hash,
		StorageKey:       storageKey,
	}, nil
}

func (s *LocalStorage) Open(storageKey string) (*os.File, error) {
	parts, err := validateStorageKey(storageKey)
	if err != nil {
		return nil, err
	}

	documentsDir := filepath.Join(s.root, parts[0])
	if err := ensureRealDirectory(documentsDir); err != nil {
		return nil, err
	}
	prefixDir := filepath.Join(documentsDir, parts[1])
	if err := ensureRealDirectory(prefixDir); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(prefixDir, parts[2])
	info, err := os.Lstat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open stored document: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: unsafe document file", ErrInvalidStorageKey)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open stored document: %w", err)
	}
	return file, nil
}

func validateStorageKey(storageKey string) ([]string, error) {
	if storageKey == "" || strings.HasPrefix(storageKey, "/") || filepath.IsAbs(storageKey) || strings.Contains(storageKey, "\\") {
		return nil, ErrInvalidStorageKey
	}

	parts := strings.Split(storageKey, "/")
	if len(parts) != 3 {
		return nil, ErrInvalidStorageKey
	}
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return nil, ErrInvalidStorageKey
		}
	}
	if parts[0] != "documents" || !isLowerHexPair(parts[1]) {
		return nil, ErrInvalidStorageKey
	}
	if !strings.HasSuffix(parts[2], ".bin") {
		return nil, ErrInvalidStorageKey
	}
	id := strings.TrimSuffix(parts[2], ".bin")
	parsed, err := uuid.Parse(id)
	if err != nil || parsed.String() != id {
		return nil, ErrInvalidStorageKey
	}
	return parts, nil
}

func isLowerHexPair(value string) bool {
	if len(value) != 2 {
		return false
	}
	for _, char := range value {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}

func ensureRealDirectory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect storage directory: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fmt.Errorf("%w: unsafe document directory", ErrInvalidStorageKey)
	}
	return nil
}

func fileExtension(fileName string) *string {
	extension := strings.ToLower(filepath.Ext(fileName))
	if extension == "" {
		return nil
	}
	if len([]rune(extension)) > 32 {
		extension = string([]rune(extension)[:32])
	}
	return &extension
}
