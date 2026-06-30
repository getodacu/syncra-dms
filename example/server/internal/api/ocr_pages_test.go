package api

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCountUploadPages(t *testing.T) {
	t.Run("PNG is one page", func(t *testing.T) {
		got, err := countUploadPages("image/png", validPNGBytes())
		if err != nil {
			t.Fatalf("countUploadPages() error = %v", err)
		}
		if got != 1 {
			t.Fatalf("page count = %d, want 1", got)
		}
	})

	t.Run("JPEG is one page", func(t *testing.T) {
		got, err := countUploadPages("image/jpeg", []byte{0xff, 0xd8, 0xff, 0xdb})
		if err != nil {
			t.Fatalf("countUploadPages() error = %v", err)
		}
		if got != 1 {
			t.Fatalf("page count = %d, want 1", got)
		}
	})

	t.Run("PDF counts page objects", func(t *testing.T) {
		got, err := countUploadPages("application/pdf", twoPagePDFBytes())
		if err != nil {
			t.Fatalf("countUploadPages() error = %v", err)
		}
		if got != 2 {
			t.Fatalf("page count = %d, want 2", got)
		}
	})

	t.Run("PDF counts pages with relaxed whitespace", func(t *testing.T) {
		got, err := countUploadPages("application/pdf", relaxedWhitespacePDFBytes())
		if err != nil {
			t.Fatalf("countUploadPages() error = %v", err)
		}
		if got != 1 {
			t.Fatalf("page count = %d, want 1", got)
		}
	})

	t.Run("PDF without pages is invalid", func(t *testing.T) {
		if _, err := countUploadPages("application/pdf", validPDFBytes()); err == nil {
			t.Fatal("countUploadPages() error = nil, want invalid PDF page count error")
		}
	})

	t.Run("PDF page-looking text is not counted as a page", func(t *testing.T) {
		if _, err := countUploadPages("application/pdf", zeroPagePDFWithPageTextBytes()); err == nil {
			t.Fatal("countUploadPages() error = nil, want invalid PDF page count error")
		}
	})

	t.Run("malformed xref-backed PDF returns an error", func(t *testing.T) {
		if _, err := countUploadPages("application/pdf", malformedXrefPDFBytes()); err == nil {
			t.Fatal("countUploadPages() error = nil, want invalid PDF page count error")
		}
	})

	t.Run("unsupported MIME errors", func(t *testing.T) {
		if _, err := countUploadPages("text/plain", []byte("x")); err == nil {
			t.Fatal("countUploadPages() error = nil, want unsupported MIME error")
		}
	})
}

func twoPagePDFBytes() []byte {
	return testPDFBytes(
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Count 2 /Kids [3 0 R 4 0 R] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
	)
}

func threePagePDFBytes() []byte {
	return testPDFBytes(
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Count 3 /Kids [3 0 R 4 0 R 5 0 R] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
	)
}

func pageCountPDFBytes(pageCount int) []byte {
	kids := make([]byte, 0, pageCount*6)
	objects := make([]string, 0, pageCount+2)
	objects = append(objects, "<< /Type /Catalog /Pages 2 0 R >>")
	for i := 0; i < pageCount; i++ {
		kids = fmt.Appendf(kids, "%d 0 R ", i+3)
	}
	objects = append(objects, fmt.Sprintf("<< /Type /Pages /Count %d /Kids [%s] >>", pageCount, kids))
	for i := 0; i < pageCount; i++ {
		objects = append(objects, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>")
	}
	return testPDFBytes(objects...)
}

func relaxedWhitespacePDFBytes() []byte {
	return testPDFBytesWithWhitespace(
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Count 1 /Kids [3 0 R] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
	)
}

func zeroPagePDFWithPageTextBytes() []byte {
	return testPDFBytes(
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Count 0 /Kids [] >>",
		"<< /Length 10 >>\nstream\n/Type /Page\nendstream",
	)
}

func malformedXrefPDFBytes() []byte {
	return testPDFBytes(
		"<< /Type /Catalog /Pages 2 0 R",
		"<< /Type /Pages /Count 1 /Kids [3 0 R] >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>",
	)
}

func testPDFBytes(objects ...string) []byte {
	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects)+1)
	for i, object := range objects {
		offsets[i+1] = out.Len()
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", i+1, object)
	}
	xrefOffset := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n", len(objects)+1)
	out.WriteString("0000000000 65535 f \n")
	for i := 1; i <= len(objects); i++ {
		fmt.Fprintf(&out, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&out, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, xrefOffset)
	return out.Bytes()
}

func testPDFBytesWithWhitespace(objects ...string) []byte {
	var out bytes.Buffer
	out.WriteString("%PDF-1.7 \n")
	offsets := make([]int, len(objects)+1)
	for i, object := range objects {
		offsets[i+1] = out.Len()
		fmt.Fprintf(&out, "%d 0 obj \n%s\nendobj\n", i+1, object)
	}
	xrefOffset := out.Len()
	fmt.Fprintf(&out, "xref \n0 %d \n", len(objects)+1)
	out.WriteString("0000000000 65535 f \n")
	for i := 1; i <= len(objects); i++ {
		fmt.Fprintf(&out, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&out, "trailer \n<< /Size %d /Root 1 0 R >>\nstartxref \n%d \n%%%%EOF\r\n", len(objects)+1, xrefOffset)
	return out.Bytes()
}
