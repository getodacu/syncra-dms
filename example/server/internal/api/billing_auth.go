package api

import "github.com/gin-gonic/gin"

const internalAPIHeader = "X-Syncra-Internal-Token"

func (h *Handler) trustedInternalRequest(c *gin.Context) bool {
	return h.InternalAPIToken != "" && c.GetHeader(internalAPIHeader) == h.InternalAPIToken
}

func (h *Handler) requireTrustedInternalRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.trustedInternalRequest(c) {
			writeError(c, 401, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
