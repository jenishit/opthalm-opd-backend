package http

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
	"github.com/redis/go-redis/v9"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationType       = "Bearer"
	authorizationPayloadKey = "authorization_payload"
	requestIDKey            = "request_id"
	requestIDHeader         = "X-Request-ID"
)

// requestIDMiddleware assigns a request ID (reusing one supplied by the
// caller/proxy in X-Request-ID, if any) so every log line for a request can
// be correlated, and so a client can hand support a single ID to search for.
func requestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetHeader(requestIDHeader)
