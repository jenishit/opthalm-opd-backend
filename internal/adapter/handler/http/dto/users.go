package dto

type CreateUser struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"` //the acceptable string type is of email type
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone"`
	RoleName  string `json:"role_name" binding:"required"`
}
