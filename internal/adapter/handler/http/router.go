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
	patientHandler PatientHandler,
	visitHandler VisitHandler,
	catalogHandler CatalogHandler,
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

	patient := api.Group("/patient").Use(authMiddleware(token))
	{
		patient.POST("", patientHandler.CreatePatient)
		patient.GET("", patientHandler.GetPatients)
		patient.GET("/search", patientHandler.SearchPatients)
		patient.GET("/:id", patientHandler.GetPatientByID)
		patient.PATCH("/:id", patientHandler.UpdatePatientByID)
		patient.PATCH("/:id/delete", patientHandler.DeletePatientByID)
	}

	visit := api.Group("/visit").Use(authMiddleware(token))
	{
		visit.POST("", visitHandler.CreateVisit)
		visit.GET("/:id", visitHandler.GetVisitByVisitID)
		visit.GET("/patient/:patientId", visitHandler.GetVisitsByPatientID)
		visit.PATCH("/:id", visitHandler.UpdateVisitByVisitID)
	}

	catalog := admin.Group("/catalog")
	{
		medicines := catalog.Group("/medicines")
		{
			medicines.GET("", catalogHandler.ListMedicines)
			medicines.GET("/search", catalogHandler.SearchMedicines)
			medicines.GET("/:id", catalogHandler.GetMedicineByID)
			medicines.PATCH("/:id", catalogHandler.UpdateMedicine)
			medicines.PATCH("/:id/delete", catalogHandler.DeleteMedicine)
		}

		diagnoses := catalog.Group("/diagnoses")
		{
			diagnoses.GET("", catalogHandler.ListDiagnoses)
			diagnoses.GET("/search", catalogHandler.SearchDiagnoses)
			diagnoses.GET("/:id", catalogHandler.GetDiagnosisByID)
			diagnoses.PATCH("/:id", catalogHandler.UpdateDiagnosis)
			diagnoses.PATCH("/:id/delete", catalogHandler.DeleteDiagnosis)
		}

		conditions := catalog.Group("/conditions")
		{
			conditions.GET("", catalogHandler.ListConditions)
			conditions.GET("/search", catalogHandler.SearchConditions)
			conditions.GET("/:id", catalogHandler.GetConditionByID)
			conditions.PATCH("/:id", catalogHandler.UpdateCondition)
			conditions.PATCH("/:id/delete", catalogHandler.DeleteCondition)
		}
	}

	return &Router{
		router,
	}, nil
}

func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
