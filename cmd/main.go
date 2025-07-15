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
