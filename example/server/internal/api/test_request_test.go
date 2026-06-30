package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func newTestRequest(method string, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)
	if strings.HasPrefix(req.URL.Path, "/api/") {
		req.Header.Set(internalAPIHeader, testInternalAPIToken)
	}
	return req
}
