package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

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
) (*Router, error) {

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	role := api.Group("/role")
	{
		role.POST("/create", roleHandler.CreateRole)
	}

	user := api.Group("/user")
	{
		user.POST("/create", userHandler.CreateUser)
	}

	profile := api.Group("/profile").Use(authMiddleware(token))
	{
		profile.GET("/getme", profileHandler.GetProfileByID)
		profile.GET("/profile-details", profileHandler.GetProfiles)
		profile.PATCH("/update-profile/:id", profileHandler.UpdateProfileByUserID)
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
