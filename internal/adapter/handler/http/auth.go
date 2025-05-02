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

	res, err := h.authService.Signup(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, res)
}

// Refresh godoc
//
//	@Summary		Refresh an access token
//	@Description	Exchanges a valid, not-yet-used refresh token for a new access+refresh pair, rotating the old one. Presenting an already-rotated-out refresh token is treated as theft and revokes every session for that user.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	response{data=domain.LoginResponse}
//	@Failure		401		{object}	errorResponse	"invalid, expired, or reused refresh token"
//	@Router			/auth/refresh [post]
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var req domain.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	res, err := h.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, res)
}

// Logout godoc
//
//	@Summary		Log out
//	@Description	Revokes the session backing the given refresh token. The already-issued access token keeps working until it expires (at most TOKEN_DURATION) — access tokens are stateless and not checked against sessions on every request.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	response
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	var req domain.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if err := h.authService.Logout(ctx, req.RefreshToken); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Logged out successfully"})
}

// RequestPasswordReset godoc
//
//	@Summary		Request a password reset code
//	@Description	Always returns the same success response regardless of whether the email belongs to an account, so this endpoint can't be used to enumerate registered emails. If it does, a single-use code valid for 1 hour is emailed.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.RequestPasswordResetRequest	true	"Account email"
//	@Success		200		{object}	response
//	@Failure		429		{object}	errorResponse	"rate limited"
//	@Router			/auth/password-reset/request [post]
func (h *AuthHandler) RequestPasswordReset(ctx *gin.Context) {
	var req domain.RequestPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if err := h.authService.RequestPasswordReset(ctx, req.Email); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "If an account exists for that email, a reset code has been sent."})
}

// ConfirmPasswordReset godoc
//
//	@Summary		Confirm a password reset
//	@Description	Applies a new password using the code from RequestPasswordReset, then revokes every existing session for that user (forcing re-login everywhere).
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.ConfirmPasswordResetRequest	true	"Reset code + new password"
//	@Success		200		{object}	response
//	@Failure		401		{object}	errorResponse	"invalid, expired, or already-used code"
//	@Router			/auth/password-reset/confirm [post]
func (h *AuthHandler) ConfirmPasswordReset(ctx *gin.Context) {
	var req domain.ConfirmPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if err := h.authService.ConfirmPasswordReset(ctx, req.Token, req.NewPassword); err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Password has been reset. Please log in again."})
}

// RequestEmailVerification godoc
//
//	@Summary		Resend the email verification code
//	@Description	Sends a fresh verification code (valid 24h) to the authenticated caller's own email address, invalidating any previously issued one.
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response
//	@Failure		401	{object}	errorResponse
//	@Router			/auth/email/verify/resend [post]
func (h *AuthHandler) RequestEmailVerification(ctx *gin.Context) {
	payload, ok := currentUser(ctx)
	if !ok {
		return
	}

