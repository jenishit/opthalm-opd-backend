package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jenish-brainztechs/go-backend/docs" // registers the generated swagger spec
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// shutdownGracePeriod is how long Serve waits for in-flight requests to
// finish draining after a shutdown signal before forcing the listener closed.
const shutdownGracePeriod = 15 * time.Second

type Router struct {
	*gin.Engine
}

func NewRouter(
	config *config.Container,
	token port.TokenService,
	roleHandler RoleHandler,
	userHandler UserHandler,
	profileHandler ProfileHandler,
	authHandler AuthHandler,
	clinicHandler ClinicHandler,
	patientHandler PatientHandler,
	visitHandler VisitHandler,
	catalogHandler CatalogHandler,
	invoiceHandler InvoiceHandler,
	inventoryHandler InventoryHandler,
	vendorHandler VendorHandler,
	stockPurchaseHandler StockPurchaseHandler,
	labJobHandler LabJobHandler,
	reportsHandler ReportsHandler,
	calculatorHandler CalculatorHandler,
	subscriptionHandler SubscriptionHandler,
	subSvc port.SubscriptionService,
	redisClient *redis.Client,
) (*Router, error) {

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(requestIDMiddleware())
	router.Use(requestLoggerMiddleware())
	router.Use(recoveryMiddleware())
	router.Use(CORSMiddleware(config.HTTP.AllowedOrigins))

	// Swagger UI + spec, for exploring/testing the API and for frontend
	// codegen. Left off in production — it documents the full request/
	// response shape of the whole API and is dev/staging tooling, not
	// something to leave publicly reachable once actually deployed.
	if config.App.Env != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", rateLimitMiddleware(redisClient, "login", 10, time.Minute), authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/password-reset/request", rateLimitMiddleware(redisClient, "pwreset", 5, time.Hour), authHandler.RequestPasswordReset)
		auth.POST("/password-reset/confirm", rateLimitMiddleware(redisClient, "pwreset-confirm", 10, time.Hour), authHandler.ConfirmPasswordReset)
		auth.POST("/email/verify/confirm", authHandler.ConfirmEmailVerification)
		auth.POST("/email/verify/resend", authMiddleware(token), authHandler.RequestEmailVerification)
	}

	api.POST("/signup", rateLimitMiddleware(redisClient, "signup", 5, time.Hour), authHandler.Signup)

	platform := api.Group("/platform")
	platform.Use(authMiddleware(token), superadminMiddleware())
	{
		subscriptions := platform.Group("/subscriptions")
		{
			subscriptions.GET("/:clinicId", subscriptionHandler.GetByClinicID)
			subscriptions.PUT("/:clinicId", subscriptionHandler.Upsert)
		}
	}

	admin := api.Group("/admin")
	admin.Use(authMiddleware(token), adminMiddleware(), subscriptionMiddleware(subSvc))

	role := api.Group("/role")
	role.Use(authMiddleware(token), adminMiddleware(), subscriptionMiddleware(subSvc))
	{
		role.POST("/create", roleHandler.CreateRole)
	}

	user := api.Group("/user")
	user.Use(authMiddleware(token), adminMiddleware(), subscriptionMiddleware(subSvc))
	{
		user.POST("/create", userHandler.CreateUser)
	}

	profile := api.Group("/profile")
	profile.Use(authMiddleware(token), subscriptionMiddleware(subSvc))
	{
		profile.GET("/search", profileHandler.SearchProfiles)
		profile.GET("/getme", profileHandler.GetProfileByID)
		profile.GET("/profile-details", profileHandler.GetProfiles)
		profile.PATCH("/update-profile/:id", profileHandler.UpdateProfileByUserID)
	}

	clinics := admin.Group("/clinic")
	{
		clinics.POST("", clinicHandler.InsertClinic)
		clinics.GET("", clinicHandler.GetAllClinics)
		clinics.GET("/:id", clinicHandler.GetClinicByID)
		clinics.PATCH("/:id", clinicHandler.UpdateClinic)
	}

