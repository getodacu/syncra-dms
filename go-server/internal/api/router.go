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
	registerSwaggerRoutes(router)

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

	users := newUserHandler(options, auth)
	userAPI := router.Group("/api/users")
	userAPI.Use(auth.requireTrustedInternalRequest())
	userAPI.GET("", users.list)
	userAPI.GET("/:id", users.get)
	userAPI.POST("", users.create)
	userAPI.PATCH("/:id", users.update)
	userAPI.POST("/:id/activate", users.activate)
	userAPI.POST("/:id/deactivate", users.deactivate)
	userAPI.POST("/:id/suspend", users.suspend)
	userAPI.DELETE("/:id", users.softDelete)
	userAPI.POST("/:id/primary-organization-unit", users.setPrimaryOrganizationUnit)
	userAPI.POST("/:id/roles", users.assignRole)
	userAPI.DELETE("/:id/roles/:assignmentId", users.removeRole)
	userAPI.POST("/:id/groups", users.addGroup)
	userAPI.DELETE("/:id/groups/:groupId", users.removeGroup)

	permissions := newPermissionHandler(options, auth)
	permissionAPI := router.Group("/api/permissions")
	permissionAPI.Use(auth.requireTrustedInternalRequest())
	permissionAPI.GET("", permissions.list)
	permissionAPI.GET("/categories", permissions.categories)

	roles := newRoleHandler(options, auth)
	roleAPI := router.Group("/api/roles")
	roleAPI.Use(auth.requireTrustedInternalRequest())
	roleAPI.GET("", roles.list)
	roleAPI.GET("/:id", roles.get)
	roleAPI.POST("", roles.create)
	roleAPI.PATCH("/:id", roles.update)
	roleAPI.DELETE("/:id", roles.delete)
	roleAPI.GET("/:id/permissions", roles.listPermissions)
	roleAPI.POST("/:id/permissions", roles.assignPermission)
	roleAPI.DELETE("/:id/permissions/:permissionId", roles.removePermission)

	groups := newGroupHandler(options, auth)
	groupAPI := router.Group("/api/groups")
	groupAPI.Use(auth.requireTrustedInternalRequest())
	groupAPI.GET("", groups.list)
	groupAPI.GET("/:id", groups.get)
	groupAPI.POST("", groups.create)
	groupAPI.PATCH("/:id", groups.update)
	groupAPI.DELETE("/:id", groups.delete)
	groupAPI.POST("/:id/users", groups.addUser)
	groupAPI.DELETE("/:id/users/:userId", groups.removeUser)
	groupAPI.POST("/:id/roles", groups.assignRole)
	groupAPI.DELETE("/:id/roles/:assignmentId", groups.removeRole)

	return router
}
