package domain

type SignupRequest struct {
	ClinicName     string `json:"clinic_name" binding:"required"`
	RegistrationNo string `json:"registration_no" binding:"required"`
	AdminFirstName string `json:"admin_first_name" binding:"required"`
	AdminLastName  string `json:"admin_last_name" binding:"required"`
	AdminEmail     string `json:"admin_email" binding:"required,email"`
	AdminPassword  string `json:"admin_password" binding:"required"`
}
