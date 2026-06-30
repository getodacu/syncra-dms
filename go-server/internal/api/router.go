package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionInfo struct {
	AppName string
	Module  string
	Version string
}

type RouterOptions struct {
	Version VersionInfo
	Ready   func(context.Context) error
}

func NewRouter(options RouterOptions) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())

	version := options.Version
	if version.AppName == "" {
		version.AppName = "Syncra DMS"
	}
	if version.Module == "" {
		version.Module = "ai.ro/syncra/dms"
	}
	if version.Version == "" {
		version.Version = "dev"
	}

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if options.Ready != nil {
			if err := options.Ready(c.Request.Context()); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "not_ready",
					"error":  err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     version.AppName,
			"module":  version.Module,
			"version": version.Version,
		})
	})

	return router
}
