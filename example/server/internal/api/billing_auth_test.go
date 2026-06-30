package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTrustedInternalRequestRequiresConfiguredExactHeader(t *testing.T) {
	tests := []struct {
		name        string
		configToken string
		headerToken string
		want        bool
	}{
		{
			name:        "configured exact match",
			configToken: "trusted-internal-token",
			headerToken: "trusted-internal-token",
			want:        true,
		},
		{
			name:        "configured missing header",
			configToken: "trusted-internal-token",
			want:        false,
		},
		{
			name:        "configured mismatched header",
			configToken: "trusted-internal-token",
			headerToken: "other-token",
			want:        false,
		},
		{
			name:        "unconfigured rejects matching empty header",
			configToken: "",
			headerToken: "",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("POST", "/api/billing/orders", nil)
			if tt.headerToken != "" {
				c.Request.Header.Set(internalAPIHeader, tt.headerToken)
			}

			got := (&Handler{InternalAPIToken: tt.configToken}).trustedInternalRequest(c)
			if got != tt.want {
				t.Fatalf("trustedInternalRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
