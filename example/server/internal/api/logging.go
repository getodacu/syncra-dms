package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"unicode"

	"ai.ro/syncra/internal/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
	ginLoggerKey    = "syncra.logger"
	ginRequestIDKey = "syncra.request_id"
)

func (h *Handler) logger() *slog.Logger {
	if h != nil && h.Logger != nil {
		return h.Logger
	}
	return logging.Nop()
}

func loggerFromGin(c *gin.Context) *slog.Logger {
	if c != nil {
		if value, ok := c.Get(ginLoggerKey); ok {
			if logger, ok := value.(*slog.Logger); ok && logger != nil {
				return logger
			}
		}
		return logging.FromContext(c.Request.Context())
	}
	return slog.Default()
}

func requestIDFromGin(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if value, ok := c.Get(ginRequestIDKey); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}

func slogRequestMiddleware(root *slog.Logger) gin.HandlerFunc {
	if root == nil {
		root = logging.Nop()
	}
	return func(c *gin.Context) {
		start := time.Now()
		requestID := cleanRequestID(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Writer.Header().Set(requestIDHeader, requestID)

		requestLogger := root.With("component", "http", "request_id", requestID)
		c.Set(ginLoggerKey, requestLogger)
		c.Set(ginRequestIDKey, requestID)
		c.Request = c.Request.WithContext(logging.WithContext(c.Request.Context(), requestLogger))

		c.Next()

		status := c.Writer.Status()
		attrs := requestLogAttrs(c, start, status)
		switch {
		case status >= http.StatusInternalServerError:
			requestLogger.Error("http.request_completed", attrs...)
		case status >= http.StatusBadRequest:
			requestLogger.Warn("http.request_completed", attrs...)
		default:
			requestLogger.Info("http.request_completed", attrs...)
		}
	}
}

func slogRecoveryMiddleware(root *slog.Logger) gin.HandlerFunc {
	if root == nil {
		root = logging.Nop()
	}
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger := loggerFromGin(c)
				if logger == nil {
					logger = root
				}
				logger.Error("http.request_panic",
					"panic", fmt.Sprint(recovered),
					"method", c.Request.Method,
					"route", routeForLog(c),
					"path", c.Request.URL.Path,
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func requestLogAttrs(c *gin.Context, start time.Time, status int) []any {
	attrs := []any{
		"method", c.Request.Method,
		"route", routeForLog(c),
		"path", c.Request.URL.Path,
		"status", status,
		"duration_ms", time.Since(start).Milliseconds(),
		"client_ip", c.ClientIP(),
		"user_agent_class", userAgentClass(c.Request.UserAgent()),
	}
	if c.Request.ContentLength >= 0 {
		attrs = append(attrs, "request_bytes", c.Request.ContentLength)
	}
	if responseBytes := c.Writer.Size(); responseBytes >= 0 {
		attrs = append(attrs, "response_bytes", responseBytes)
	}
	return attrs
}

func routeForLog(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if route := c.FullPath(); route != "" {
		return route
	}
	return "unmatched"
}

func cleanRequestID(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" || len(value) > 128 {
		return ""
	}
	for _, r := range value {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return ""
		}
	}
	return value
}

func userAgentClass(raw string) string {
	ua := strings.ToLower(strings.TrimSpace(raw))
	if ua == "" {
		return "absent"
	}
	switch {
	case strings.Contains(ua, "mozilla/"):
		return "browser"
	case strings.Contains(ua, "curl/"):
		return "curl"
	case strings.Contains(ua, "go-http-client/"):
		return "go-http-client"
	default:
		return "api-client"
	}
}

func resolvedSchemaLogAttrs(schema resolvedSchema) []any {
	source := "none"
	switch {
	case schema.SchemaID != nil:
		source = "saved"
	case schema.Inline:
		source = "inline"
	}
	attrs := []any{
		"schema_source", source,
		"has_schema", len(schema.Schema) > 0,
		"strict", schema.Strict,
	}
	if schema.SchemaID != nil {
		attrs = append(attrs, "schema_id", schema.SchemaID.String())
	}
	return attrs
}

func safeLogError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	switch {
	case strings.Contains(message, "write OCR job file"):
		return "write OCR job file"
	case strings.Contains(message, "remove OCR job file"):
		return "remove OCR job file"
	case strings.Contains(message, "read OCR job file"):
		return "read OCR job file"
	default:
		return message
	}
}
