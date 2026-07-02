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
	ErrUploadTooLarge          = errors.New("document exceeds maximum upload size")
	ErrUnsupportedMIME         = errors.New("document MIME type is not allowed")
	ErrInvalidStorageKey       = errors.New("invalid storage key")
	ErrUnsafeStorageDirectory  = errors.New("unsafe storage directory")
	privateStorageDirMode      = os.FileMode(0o700)
	groupOrOtherPermissionMask = os.FileMode(0o077)
)

// LocalStorage expects root to be private to the service user. It makes managed
// subdirectories private and rejects symlinks below root, but it is not a
// same-UID filesystem race defense.
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
	tmpDir := filepath.Join(s.root, "tmp")
	if err := os.MkdirAll(tmpDir, privateStorageDirMode); err != nil {
		return StoredFile{}, fmt.Errorf("create storage temp directory: %w", err)
	}
	if err := ensurePrivateDirectory(tmpDir); err != nil {
		return StoredFile{}, err
	}

	tmpFile, err := os.CreateTemp(tmpDir, "upload-*")
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
	if err := os.MkdirAll(documentsDir, privateStorageDirMode); err != nil {
		return StoredFile{}, fmt.Errorf("create storage documents directory: %w", err)
	}
	if err := ensurePrivateDirectory(documentsDir); err != nil {
		return StoredFile{}, err
	}

	targetDir := filepath.Join(documentsDir, hashPrefix)
	if err := os.MkdirAll(targetDir, privateStorageDirMode); err != nil {
		return StoredFile{}, fmt.Errorf("create storage document directory: %w", err)
	}
	if err := ensurePrivateDirectory(targetDir); err != nil {
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
	if err := ensurePrivateDirectory(documentsDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("open stored document: %w", err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidStorageKey, err)
	}
	prefixDir := filepath.Join(documentsDir, parts[1])
	if err := ensurePrivateDirectory(prefixDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("open stored document: %w", err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidStorageKey, err)
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

func (s *LocalStorage) Delete(storageKey string) error {
	parts, err := validateStorageKey(storageKey)
	if err != nil {
		return err
	}

	documentsDir := filepath.Join(s.root, parts[0])
	if err := ensurePrivateDirectory(documentsDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("%w: %v", ErrInvalidStorageKey, err)
	}
	prefixDir := filepath.Join(documentsDir, parts[1])
	if err := ensurePrivateDirectory(prefixDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("%w: %v", ErrInvalidStorageKey, err)
	}

	fullPath := filepath.Join(prefixDir, parts[2])
	info, err := os.Lstat(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("delete stored document: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return fmt.Errorf("%w: unsafe document file", ErrInvalidStorageKey)
	}
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("delete stored document: %w", err)
	}
	return nil
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

func ensurePrivateDirectory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("%w: inspect %s: %w", ErrUnsafeStorageDirectory, path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fmt.Errorf("%w: %s is not a real directory", ErrUnsafeStorageDirectory, path)
	}

	if info.Mode().Perm() != privateStorageDirMode {
		if err := os.Chmod(path, privateStorageDirMode); err != nil {
			return fmt.Errorf("%w: make %s private: %w", ErrUnsafeStorageDirectory, path, err)
		}
		info, err = os.Lstat(path)
		if err != nil {
			return fmt.Errorf("%w: inspect %s after chmod: %w", ErrUnsafeStorageDirectory, path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("%w: %s is not a real directory after chmod", ErrUnsafeStorageDirectory, path)
		}
		if info.Mode().Perm()&groupOrOtherPermissionMask != 0 {
			return fmt.Errorf("%w: %s remains group/other accessible after chmod", ErrUnsafeStorageDirectory, path)
		}
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
