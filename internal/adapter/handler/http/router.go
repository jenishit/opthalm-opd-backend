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
