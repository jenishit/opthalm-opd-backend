package dto

type CreateRole struct { //binding required is the compulsory values
	RoleName string `json:"role_name" binding:"required"`
}
