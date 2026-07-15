package main

import (
	"context"
	"fmt"
	"log/slog"

	"os"

	auth "github.com/jenish-brainztechs/go-backend/internal/adapter/auth/jwt"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres/repository"
	"github.com/jenish-brainztechs/go-backend/internal/core/services"
)

func main() {
	config, err := config.New() //Creating a new configuration for the application
	if err != nil {             //if there is some error then print and log the error and exit from the application
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Database has been initializerd and connected successfully", "db", config.DB.Connection)

	tokenService, err := auth.New(config.Token)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	roleRepo := repository.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := http.NewRoleHandler(roleService)

	profileRepo := repository.NewProfileRepository(db)
	profileService := services.NewProfileService(profileRepo)
	profileHandler := http.NewProfileHandler(profileService)

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, tokenService)
	authHandler := http.NewAuthHandler(authService)
	userService := services.NewUserService(userRepo, roleService, profileService)
	userHandler := http.NewUsersHandler(userService)

	router, err := http.NewRouter(
		config,
		tokenService,
		*roleHandler,
		*userHandler,
		*profileHandler,
		*authHandler,
	)

	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)

	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
