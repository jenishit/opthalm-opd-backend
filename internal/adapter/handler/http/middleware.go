package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationType       = "Bearer"
	authorizationPayloadKey = "authorization_payload"
	sessionContextKey       = "session_state"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-agent-code")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
func authMiddleware(token port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		isEmpty := len(authorizationHeader) == 0

		if isEmpty {
			handleAbort(ctx, domain.ErrEmptyAuthorizationHeader)
			return

		}

		fields := strings.Fields(authorizationHeader)

		isValid := len(fields) == 2

		if !isValid {
			handleAbort(ctx, domain.ErrInvalidAuthorizationHeader)
			return
		}
		currentAuthorizationtype := fields[0]

		if currentAuthorizationtype != authorizationType {
			handleAbort(ctx, domain.ErrInvalidAuthorizationType)
			return

		}
		accessToken := fields[1]

		payload, err := token.VerifyAccessToken(accessToken)

		if err != nil {
			handleAbort(ctx, err)
			return

		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			validationError(ctx, domain.ErrEmptyAuthorizationHeader)
			return
		}

		userPayload, ok := payload.(*domain.TokenPayload)
		if !ok {
			validationError(ctx, domain.ErrInvalidAuthorizationHeader)
			return
		}

		if userPayload.RoleName != "ROLE_ADMIN" {
			handleAbort(ctx, domain.ErrUnauthorized)
			return

		}
		ctx.Next()

	}
}
