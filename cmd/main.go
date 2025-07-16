// Package main is the Opthalmic Management System API server.
//
//	@title			Opthalmic Management System API
//	@version		1.0
//	@description	Backend API for a multi-tenant (SaaS) ophthalmology clinic
//	@description	management system: patients, visits, billing/POS, inventory,
//	@description	lab jobs, reports, and optical calculators.
//	@description
//	@description	Every route under /api (other than /api/auth/*, /api/signup,
//	@description	and the /api/auth/password-reset|email endpoints) requires a
//	@description	Bearer access token, and most additionally require the
//	@description	caller's clinic to have an active subscription.
//
//	@contact.name	API Support
//
//	@license.name	Proprietary
//
//	@host		localhost:8082
//	@BasePath	/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and the access token, e.g. "Bearer eyJhbGci...".
package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"os"

	auth "github.com/jenish-brainztechs/go-backend/internal/adapter/auth/jwt"
	redisadapter "github.com/jenish-brainztechs/go-backend/internal/adapter/cache/redis"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/email/logsender"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/email/smtpsender"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres/repository"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
	"github.com/jenish-brainztechs/go-backend/internal/core/services"
)

func main() {
	config, err := config.New() //Creating a new configuration for the application
	if err != nil {             //if there is some error then print and log the error and exit from the application
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	initLogger(config.App.Env)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := redisadapter.New(ctx, config.Redis)
	if err != nil {
		slog.Error("Error initializing redis connection", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	slog.Info("Database has been initializerd and connected successfully", "db", config.DB.Connection)

	tokenService, err := auth.New(config.Token)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	refreshDuration, err := time.ParseDuration(config.Refresh.Duration)
	if err != nil {
		slog.Error("Error parsing refresh token duration", "error", err)
		os.Exit(1)
	}

	roleRepo := repository.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := http.NewRoleHandler(roleService)

	profileRepo := repository.NewProfileRepository(db)
	profileService := services.NewProfileService(profileRepo)
	profileHandler := http.NewProfileHandler(profileService)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	signupRepo := repository.NewSignupRepository(db)
	verificationRepo := repository.NewVerificationTokenRepository(db)

	var emailSender port.EmailSender
	if config.SMTP.Host != "" {
		emailSender = smtpsender.New(smtpsender.Config{
			Host:     config.SMTP.Host,
			Port:     config.SMTP.Port,
			Username: config.SMTP.Username,
			Password: config.SMTP.Password,
			From:     config.SMTP.From,
		})
		slog.Info("email sender configured", "driver", "smtp", "host", config.SMTP.Host)
	} else {
		emailSender = logsender.New()
		slog.Info("email sender configured", "driver", "log", "note", "set SMTP_HOST to enable real delivery")
	}

	authService := services.NewAuthService(userRepo, sessionRepo, signupRepo, verificationRepo, emailSender, tokenService, refreshDuration)
	authHandler := http.NewAuthHandler(authService)
	userService := services.NewUserService(userRepo, roleService, profileService)
	userHandler := http.NewUsersHandler(userService)

	clinicRepo := repository.NewClinicRepository(db)
	clinicService := services.NewClinicService(clinicRepo)
	clinicHandler := http.NewClinicHandler(clinicService)

	patientRepo := repository.NewPatientRepository(db)
	patientService := services.NewPatientService(patientRepo)
	patientHandler := http.NewPatientHandler(patientService)

	visitRepo := repository.NewVisitsRepository(db)
	visitService := services.NewVisitsService(visitRepo)
	visitHandler := http.NewVisitHandler(visitService)

	medicineRepo := repository.NewMedicineRepository(db)
	medicineService := services.NewMedicineService(medicineRepo)

	diagnosisRepo := repository.NewDiagnosisCatalogRepository(db)
	diagnosisService := services.NewDiagnosisCatalogService(diagnosisRepo)

	conditionRepo := repository.NewHistoryConditionRepository(db)
	conditionService := services.NewHistoryConditionService(conditionRepo)

	catalogHandler := http.NewCatalogHandler(medicineService, diagnosisService, conditionService)

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := services.NewInventoryService(inventoryRepo)
	inventoryHandler := http.NewInventoryHandler(inventoryService)

	vendorRepo := repository.NewVendorRepository(db)
	vendorService := services.NewVendorService(vendorRepo)
	vendorHandler := http.NewVendorHandler(vendorService)

	stockPurchaseRepo := repository.NewStockPurchaseRepository(db, inventoryRepo)
	stockPurchaseService := services.NewStockPurchaseService(stockPurchaseRepo)
	stockPurchaseHandler := http.NewStockPurchaseHandler(stockPurchaseService)

	invoiceRepo := repository.NewInvoiceRepository(db, inventoryRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)
	invoiceHandler := http.NewInvoiceHandler(invoiceService, clinicService)

	labJobRepo := repository.NewLabJobRepository(db)
	labJobService := services.NewLabJobService(labJobRepo)
	labJobHandler := http.NewLabJobHandler(labJobService)

	reportsRepo := repository.NewReportsRepository(db)
	reportsService := services.NewReportsService(reportsRepo, inventoryService)
	reportsHandler := http.NewReportsHandler(reportsService)

	calculatorService := services.NewCalculatorService()
	calculatorHandler := http.NewCalculatorHandler(calculatorService)

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, redisClient)
	subscriptionHandler := http.NewSubscriptionHandler(subscriptionService)

	router, err := http.NewRouter(
		config,
		tokenService,
		*roleHandler,
		*userHandler,
		*profileHandler,
		*authHandler,
		*clinicHandler,
		*patientHandler,
		*visitHandler,
		*catalogHandler,
		*invoiceHandler,
		*inventoryHandler,
		*vendorHandler,
		*stockPurchaseHandler,
		*labJobHandler,
		*reportsHandler,
		*calculatorHandler,
		*subscriptionHandler,
		subscriptionService,
		redisClient,
	)

	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)

