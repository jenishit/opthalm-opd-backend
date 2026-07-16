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
	clinicHandler ClinicHandler,
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

	admin := api.Group("/admin")
	admin.Use(authMiddleware(token), adminMiddleware())

	profile := api.Group("/profile").Use(authMiddleware(token))
	{
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

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
