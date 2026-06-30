package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai.ro/syncra/dms/docs"

	"github.com/gin-gonic/gin"
)

const swaggerIndexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Syncra DMS API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html,
    body {
      margin: 0;
      min-height: 100%;
      background: #ffffff;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
        deepLinking: true
      });
    };
  </script>
</body>
</html>
`

func registerSwaggerRoutes(router *gin.Engine) {
	router.GET("/swagger", redirectSwaggerIndex)
	router.GET("/swagger/", redirectSwaggerIndex)
	router.GET("/swagger/index.html", serveSwaggerIndex)
	router.GET("/swagger/doc.json", serveSwaggerDoc)
}

func redirectSwaggerIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}

func serveSwaggerIndex(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerIndexHTML))
}

func serveSwaggerDoc(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", swaggerDocForRequest(c.Request))
}

func swaggerDocForRequest(request *http.Request) []byte {
	var doc map[string]any
	if err := json.Unmarshal(docs.SwaggerJSON, &doc); err != nil {
		return docs.SwaggerJSON
	}

	if request.Host != "" {
		doc["host"] = request.Host
	}
	doc["schemes"] = []string{requestScheme(request)}

	rendered, err := json.Marshal(doc)
	if err != nil {
		return docs.SwaggerJSON
	}
	return rendered
}

func requestScheme(request *http.Request) string {
	if forwardedProto := strings.TrimSpace(request.Header.Get("X-Forwarded-Proto")); forwardedProto != "" {
		scheme, _, _ := strings.Cut(forwardedProto, ",")
		return strings.TrimSpace(scheme)
	}
	if request.TLS != nil {
		return "https"
	}
	return "http"
}
