package api

import (
	"bytes"
	"errors"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

var errPDFPageCount = errors.New("PDF page count could not be determined")

func init() {
	pdfcpuapi.DisableConfigDir()
}

func countUploadPages(mimeType string, data []byte) (int, error) {
	switch mimeType {
	case "image/png", "image/jpeg":
		return 1, nil
	case "application/pdf":
		return countPDFPages(data)
	default:
		return 0, errors.New("unsupported file type")
	}
}

func countPDFPages(data []byte) (count int, err error) {
	defer func() {
		if recover() != nil {
			count = 0
			err = errPDFPageCount
		}
	}()

	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	count, err = pdfcpuapi.PageCount(bytes.NewReader(data), conf)
	if err != nil {
		return 0, errPDFPageCount
	}
	if count <= 0 {
		return 0, errPDFPageCount
	}
	return count, nil
}
