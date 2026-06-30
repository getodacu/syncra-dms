package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VersionInfo struct {
	AppName string
	Module  string
	Version string
}

type RouterOptions struct {
	Version             VersionInfo
	Ready               func(context.Context) error
	DB                  *gorm.DB
	BetterAuthSecret    string
	AuthDeliveryToken   string
	InternalAPIToken    string
	AuthSessionTTL      time.Duration
	AuthVerificationTTL time.Duration
	AuthCookieSecure    bool
	GoogleClientID      string
	GoogleClientSecret  string
	GitHubClientID      string
	GitHubClientSecret  string
	OAuthProfileFetcher func(context.Context, string, string, string) (OAuthProfile, error)
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

	auth := newAuthHandler(options)
	authAPI := router.Group("/api/auth")
	authAPI.Use(auth.requireTrustedInternalRequest())
	authAPI.POST("/sign-up/email", auth.signUpEmail)
	authAPI.POST("/sign-in/email", auth.signInEmail)
	authAPI.GET("/get-session", auth.getSession)
	authAPI.POST("/sign-out", auth.signOut)
	authAPI.POST("/email-otp/send-verification-otp", auth.sendVerificationOTP)
	authAPI.POST("/email-otp/verify-email", auth.verifyEmailOTP)
	authAPI.POST("/password-reset/request", auth.requestPasswordReset)
	authAPI.POST("/password-reset/confirm", auth.confirmPasswordReset)
	authAPI.POST("/oauth/google/start", auth.startGoogleOAuth)
	authAPI.POST("/oauth/google/callback", auth.signInGoogleOAuth)
	authAPI.POST("/oauth/github/start", auth.startGitHubOAuth)
	authAPI.POST("/oauth/github/callback", auth.signInGitHubOAuth)

	orgUnits := newOrganizationUnitHandler(options, auth)
	orgUnitAPI := router.Group("/api/organization-units")
	orgUnitAPI.Use(auth.requireTrustedInternalRequest())
	orgUnitAPI.GET("/tree", orgUnits.listTree)
	orgUnitAPI.GET("/archived", orgUnits.listArchived)
	orgUnitAPI.POST("", orgUnits.create)
	orgUnitAPI.PATCH("/:id", orgUnits.update)
	orgUnitAPI.PATCH("/:id/parent", orgUnits.move)
	orgUnitAPI.POST("/:id/archive", orgUnits.archive)

	return router
}
