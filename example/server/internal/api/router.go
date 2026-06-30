package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"ai.ro/syncra/docs"
	"ai.ro/syncra/internal/ocr"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Handler struct {
	DB                    *gorm.DB
	OCR                   OCRProcessor
	OCRJobNotifier        ocr.JobNotifier
	Logger                *slog.Logger
	StorageDir            string
	MaxUploadBytes        int64
	GotenbergAPIURL       string
	MistralAPIKey         string
	MistralBaseURL        string
	MistralModel          string
	BetterAuthSecret      string
	AppPrivateKey         string
	AuthDeliveryToken     string
	InternalAPIToken      string
	AuthSessionTTL        time.Duration
	AuthVerificationTTL   time.Duration
	AuthCookieSecure      bool
	OnboardingCredits     int
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleOAuthExchange   func(context.Context, string, string) (googleOAuthToken, error)
	GoogleIDTokenValidate func(context.Context, string, string) (googleIDTokenPayload, error)
	GitHubClientID        string
	GitHubClientSecret    string
	GitHubOAuthExchange   func(context.Context, string, string) (githubOAuthToken, error)
	GitHubProfileFetch    func(context.Context, string) (githubOAuthProfile, error)
	SwaggerHost           string
	SwaggerSchemes        []string
	APIKeyGenerator       func() (string, error)
	Now                   func() time.Time
}

func NewRouter(h *Handler) *gin.Engine {
	if h == nil {
		h = &Handler{}
	}
	logger := h.logger()
	router := gin.New()
	router.Use(slogRecoveryMiddleware(logger), slogRequestMiddleware(logger))
	router.MaxMultipartMemory = 8 << 20

	internalAPI := router.Group("/api")
	internalAPI.Use(h.requireTrustedInternalRequest())

	adminAPI := internalAPI.Group("/admin")
	adminAPI.Use(h.requireAdminSession())
	adminAPI.POST("/impersonation/stop", h.StopAdminImpersonation)
	adminBilling := adminAPI.Group("/billing")
	adminBilling.GET("/invoices", h.ListAdminBillingInvoices)
	adminBilling.GET("/orders", h.ListAdminBillingOrders)
	adminBilling.POST("/orders/:id/invoice", h.CreateAdminBillingOrderInvoice)
	adminUsers := adminAPI.Group("/users")
	adminUsers.GET("", h.ListAdminUsers)
	adminUsers.GET("/:id", h.GetAdminUser)
	adminUsers.PATCH("/:id", h.PatchAdminUser)
	adminUsers.POST("/:id/impersonation", h.StartAdminUserImpersonation)
	adminUsers.POST("/:id/password", h.SetAdminUserPassword)
	adminUsers.POST("/:id/balance-adjustment", h.AdjustAdminUserBalance)
	adminUsers.PUT("/:id/billing-profile", h.UpsertAdminUserBillingProfile)
	adminJSONRecipes := adminAPI.Group("/json-recipes")
	adminJSONRecipes.GET("", h.ListJSONRecipes)
	adminJSONRecipes.POST("", h.CreateJSONRecipe)
	adminJSONRecipes.GET("/:id", h.GetJSONRecipe)
	adminJSONRecipes.PUT("/:id", h.UpdateJSONRecipe)
	adminJSONRecipes.DELETE("/:id", h.DeleteJSONRecipe)
	adminJSONRecipeCategories := adminAPI.Group("/json-recipe-categories")
	adminJSONRecipeCategories.GET("", h.ListJSONRecipeCategories)
	adminJSONRecipeCategories.POST("", h.CreateJSONRecipeCategory)
	adminJSONRecipeCategories.GET("/:id", h.GetJSONRecipeCategory)
	adminJSONRecipeCategories.PUT("/:id", h.UpdateJSONRecipeCategory)
	adminJSONRecipeCategories.DELETE("/:id", h.DeleteJSONRecipeCategory)

	internalAPI.GET("/json-recipes", h.ListJSONRecipes)
	internalAPI.POST("/json-recipes/:id/deploy", h.DeployJSONRecipe)
	internalAPI.GET("/dashboard/summary", h.GetDashboardSummary)

	ocrAPI := internalAPI.Group("/ocr")
	ocrAPI.POST("", h.CreateOCRDocument)

	documents := ocrAPI.Group("/documents")
	documents.GET("", h.ListOCRDocuments)
	documents.DELETE("", h.DeleteOCRDocuments)
	documents.PUT("/collections", h.MoveOCRDocumentsToCollections)
	documents.PATCH("/:id", h.UpdateOCRDocument)
	documents.DELETE("/:id", h.DeleteOCRDocument)
	ocrAPI.GET("/document/:id", h.GetOCRDocument)

	jobs := ocrAPI.Group("/jobs")
	jobs.POST("", h.CreateOCRJob)
	jobs.GET("", h.ListOCRJobs)
	jobs.DELETE("", h.DeleteOCRJobs)
	jobs.GET("/:id", h.GetOCRJob)

	publicAPI := router.Group("/v1")
	publicAPI.Use(publicAPICORS())
	publicAPI.OPTIONS("/*path", handlePublicAPIPreflight)
	publicAPI.Use(h.publicAPIAuth())
	publicAPI.GET("/get-balance", h.GetPublicCreditBalance)
	publicOCRJobs := publicAPI.Group("/ocr/jobs")
	publicOCRJobs.POST("", h.CreatePublicOCRJob)
	publicOCRJobs.GET("/:id", h.GetPublicOCRJob)

	invoiceStatic := router.Group("/static/invoice")
	invoiceStatic.Use(h.requireTrustedInternalRequest())
	invoiceStatic.GET("/:invoice_pdf", h.ServeBillingInvoicePDF)

	schemas := ocrAPI.Group("/schemas")
	schemas.POST("", h.CreateSchema)
	schemas.GET("", h.ListSchemas)
	schemas.DELETE("", h.DeleteSchemas)
	schemas.GET("/:id", h.GetSchema)
	schemas.PUT("/:id", h.UpdateSchema)
	schemas.DELETE("/:id", h.DeleteSchema)

	collections := internalAPI.Group("")
	collections.POST("/collection", h.CreateCollection)
	collections.GET("/collections", h.ListCollections)
	collections.GET("/collections/:id", h.GetCollection)
	collections.PUT("/collection/:id", h.UpdateCollection)
	collections.DELETE("/collection/:id", h.DeleteCollection)

	datasets := internalAPI.Group("/datasets")
	datasets.POST("", h.CreateDataset)
	datasets.GET("", h.ListDatasets)
	datasets.GET("/:id", h.GetDataset)
	datasets.GET("/:id/rows", h.GetDatasetRows)
	datasets.GET("/:id/export", h.ExportDataset)
	datasets.PUT("/:id", h.UpdateDataset)
	datasets.DELETE("/:id", h.DeleteDataset)

	authAPI := internalAPI.Group("/auth")
	authAPI.POST("/sign-up/email", h.SignUpEmail)
	authAPI.POST("/sign-in/email", h.SignInEmail)
	authAPI.POST("/oauth/google/start", h.StartGoogleOAuth)
	authAPI.POST("/oauth/google/callback", h.SignInGoogleOAuth)
	authAPI.POST("/oauth/github/start", h.StartGitHubOAuth)
	authAPI.POST("/oauth/github/callback", h.SignInGitHubOAuth)
	authAPI.GET("/sessions", h.ListAuthSessions)
	authAPI.DELETE("/sessions/:id", h.RevokeAuthSession)
	authAPI.GET("/accounts", h.ListAuthAccounts)
	authAPI.DELETE("/accounts/:provider_id", h.UnlinkAuthAccount)
	authAPI.POST("/accounts/google/start", h.StartGoogleAccountLink)
	authAPI.POST("/accounts/google/callback", h.LinkGoogleAccount)
	authAPI.POST("/accounts/github/start", h.StartGitHubAccountLink)
	authAPI.POST("/accounts/github/callback", h.LinkGitHubAccount)
	authAPI.GET("/get-session", h.GetSession)
	authAPI.PATCH("/user", h.PatchAuthUser)
	authAPI.POST("/sign-out", h.SignOut)
	authAPI.POST("/email-otp/send-verification-otp", h.SendVerificationOTP)
	authAPI.POST("/email-otp/verify-email", h.VerifyEmailOTP)
	authAPI.POST("/password-reset/request", h.RequestPasswordReset)
	authAPI.POST("/password-reset/confirm", h.ConfirmPasswordReset)
	apiKeys := authAPI.Group("/apikeys")
	apiKeys.POST("", h.CreateAPIKey)
	apiKeys.GET("/:user_id", h.ListAPIKeys)
	apiKeys.DELETE("", h.DeleteAPIKey)
	webhook := authAPI.Group("/webhook")
	webhook.GET("/:user_id", h.GetWebhook)
	webhook.POST("", h.UpsertWebhook)
	webhook.PATCH("/:user_id/secret", h.RegenerateWebhookSecret)
	webhook.DELETE("", h.DeleteWebhook)

	billingAPI := internalAPI.Group("/billing")
	billingAPI.GET("/balance", h.GetCreditBalance)
	billingAPI.GET("/profile", h.GetBillingProfile)
	billingAPI.PUT("/profile", h.UpsertBillingProfile)
	billingAPI.POST("/invoices", h.CreateBillingInvoice)
	billingAPI.GET("/invoices/:id/pdf", h.ServeUserBillingInvoicePDF)
	billingAPI.POST("/invoices/:id/email-delivery/claim", h.ClaimBillingInvoiceEmailDelivery)
	billingAPI.POST("/invoices/:id/email-delivery/sent", h.MarkBillingInvoiceEmailSent)
	billingAPI.POST("/generate-invoice-pdf", h.GenerateBillingInvoicePDF)
	billingAPI.GET("/credit-usage-history", h.ListCreditUsageHistory)
	billingAPI.GET("/orders", h.ListBillingOrders)
	billingAPI.POST("/orders", h.CreateBillingOrder)
	billingAPI.POST("/orders/:id/checkout-session", h.AttachBillingOrderCheckoutSession)
	billingAPI.POST("/orders/:id/paid", h.MarkBillingOrderPaid)
	billingAPI.POST("/orders/:id/failed", h.MarkBillingOrderFailed)

	router.GET("/swagger/doc.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", h.swaggerJSON())
	})
	swaggerUI := ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	)
	router.GET("/swagger-public/doc.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", h.publicSwaggerJSON())
	})
	publicSwaggerUI := ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger-public/doc.json"),
	)
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/swagger/") && swaggerUIPath(path, "/swagger/") {
			c.Status(http.StatusOK)
			swaggerUI(c)
			return
		}
		if strings.HasPrefix(path, "/swagger-public/") && swaggerUIPath(path, "/swagger-public/") {
			c.Status(http.StatusOK)
			publicSwaggerUI(c)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})

	return router
}

func publicAPICORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
		c.Header("Access-Control-Max-Age", "600")
		c.Next()
	}
}

func handlePublicAPIPreflight(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func swaggerUIPath(path string, prefix string) bool {
	name := strings.TrimPrefix(path, prefix)
	return !strings.Contains(name, "/") && name != "doc.json"
}

func (h *Handler) swaggerJSON() []byte {
	return h.filteredSwaggerJSON("Syncra Internal API", internalSwaggerPath)
}

func (h *Handler) publicSwaggerJSON() []byte {
	return h.filteredSwaggerJSON("Syncra Public API", publicSwaggerPath)
}

type swaggerPathPredicate func(path string) bool

func internalSwaggerPath(path string) bool {
	return !strings.HasPrefix(path, "/v1/")
}

func publicSwaggerPath(path string) bool {
	return strings.HasPrefix(path, "/v1/")
}

func (h *Handler) filteredSwaggerJSON(title string, include swaggerPathPredicate) []byte {
	var spec map[string]any
	if err := json.Unmarshal(docs.SwaggerJSON, &spec); err != nil {
		return docs.SwaggerJSON
	}
	setSwaggerInfoTitle(spec, title)
	if h.SwaggerHost != "" {
		spec["host"] = h.SwaggerHost
	}
	if len(h.SwaggerSchemes) > 0 {
		schemes := make([]string, len(h.SwaggerSchemes))
		copy(schemes, h.SwaggerSchemes)
		spec["schemes"] = schemes
	}
	filterSwaggerPaths(spec, include)
	pruneSwaggerDefinitions(spec)
	out, err := json.Marshal(spec)
	if err != nil {
		return docs.SwaggerJSON
	}
	return out
}

func setSwaggerInfoTitle(spec map[string]any, title string) {
	info, ok := spec["info"].(map[string]any)
	if !ok {
		info = make(map[string]any)
		spec["info"] = info
	}
	info["title"] = title
}

func filterSwaggerPaths(spec map[string]any, include swaggerPathPredicate) {
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		return
	}
	filtered := make(map[string]any, len(paths))
	for path, operations := range paths {
		if include(path) {
			filtered[path] = operations
		}
	}
	spec["paths"] = filtered
}

func pruneSwaggerDefinitions(spec map[string]any) {
	definitions, ok := spec["definitions"].(map[string]any)
	if !ok {
		return
	}

	needed := make(map[string]struct{})
	var visitValue func(value any)
	var visitDefinitionRef func(ref string)

	visitDefinitionRef = func(ref string) {
		const prefix = "#/definitions/"
		if !strings.HasPrefix(ref, prefix) {
			return
		}
		name := strings.TrimPrefix(ref, prefix)
		if _, ok := needed[name]; ok {
			return
		}
		definition, ok := definitions[name]
		if !ok {
			return
		}
		needed[name] = struct{}{}
		visitValue(definition)
	}

	visitValue = func(value any) {
		switch typed := value.(type) {
		case map[string]any:
			if ref, ok := typed["$ref"].(string); ok {
				visitDefinitionRef(ref)
			}
			for _, child := range typed {
				visitValue(child)
			}
		case []any:
			for _, child := range typed {
				visitValue(child)
			}
		}
	}

	visitValue(spec["paths"])

	filtered := make(map[string]any, len(needed))
	for name := range needed {
		if definition, ok := definitions[name]; ok {
			filtered[name] = definition
		}
	}
	spec["definitions"] = filtered
}
