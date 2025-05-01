package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

type AuthHandler struct {
	authService port.AuthService
}

func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
//
//	@Summary		Log in
//	@Description	Exchanges email + password for an access token (short-lived JWT) and a refresh token (long-lived, rotating).
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.Login	true	"Credentials"
//	@Success		200		{object}	response{data=domain.LoginResponse}
//	@Failure		401		{object}	errorResponse	"invalid credentials"
//	@Failure		429		{object}	errorResponse	"rate limited"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req domain.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	res, err := h.authService.Login(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, res)
}

// Signup godoc
//
//	@Summary		Sign up a new clinic
//	@Description	Creates a new clinic (tenant), its first ROLE_ADMIN user, and a 14-day trial subscription, atomically. Returns a normal login response so the new admin is immediately authenticated.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.SignupRequest	true	"New clinic + admin details"
//	@Success		200		{object}	response{data=domain.LoginResponse}
//	@Failure		400		{object}	errorResponse
//	@Failure		409		{object}	errorResponse	"email already registered"
//	@Failure		429		{object}	errorResponse	"rate limited"
//	@Router			/signup [post]
func (h *AuthHandler) Signup(ctx *gin.Context) {
	var req domain.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
