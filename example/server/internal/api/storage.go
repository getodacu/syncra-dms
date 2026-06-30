package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) storageRootDir() (string, error) {
	dir := strings.TrimSpace(h.StorageDir)
	if filepath.IsAbs(dir) {
		return filepath.Clean(dir), nil
	}

	root, err := serverRootDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return filepath.Join(root, ".data"), nil
	}
	if filepath.IsAbs(dir) {
		return filepath.Clean(dir), nil
	}
	return filepath.Join(root, dir), nil
}

func (h *Handler) ocrJobFileDir() (string, error) {
	root, err := h.storageRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "ocr-files"), nil
}

func (h *Handler) billingInvoicePDFDir() (string, error) {
	root, err := h.storageRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "invoices"), nil
}

func serverRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := cwd; ; dir = filepath.Dir(dir) {
		if isServerRoot(dir) {
			return dir, nil
		}
		serverDir := filepath.Join(dir, "server")
		if isServerRoot(serverDir) {
			return serverDir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return "", errors.New("failed to locate server root")
}

func isServerRoot(dir string) bool {
	for _, path := range []string{
		filepath.Join(dir, "go.mod"),
		filepath.Join(dir, "internal", "api"),
		filepath.Join(dir, "cmd", "syncra"),
	} {
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}
	return true
}
