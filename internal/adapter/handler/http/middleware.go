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
		if id == "" {
			id = uuid.NewString()
		}
		ctx.Set(requestIDKey, id)
		ctx.Header(requestIDHeader, id)
		ctx.Next()
	}
}

// requestLoggerMiddleware logs one structured line per request (method,
// path, status, latency, and the authenticated user/clinic when present)
// via log/slog, replacing gin's plain-text default logger.
func requestLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		if raw := ctx.Request.URL.RawQuery; raw != "" {
			path = path + "?" + raw
		}

		ctx.Next()

		status := ctx.Writer.Status()
		attrs := []any{
			"request_id", ctx.GetString(requestIDKey),
			"method", ctx.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", ctx.ClientIP(),
		}

		if payload, exists := ctx.Get(authorizationPayloadKey); exists {
			if p, ok := payload.(*domain.TokenPayload); ok {
				attrs = append(attrs, "user_id", p.UserId, "clinic_id", p.ClinicID)
			}
		}
		if len(ctx.Errors) > 0 {
			attrs = append(attrs, "gin_errors", ctx.Errors.String())
		}

		msg := "request handled"
		switch {
		case status >= 500:
			slog.Error(msg, attrs...)
		case status >= 400:
			slog.Warn(msg, attrs...)
		default:
			slog.Info(msg, attrs...)
		}
	}
}

// rateLimitMiddleware enforces a fixed-window limit of at most limit
// requests per window, keyed by client IP, backed by Redis (INCR + EXPIRE)
// so the limit holds across multiple app instances, not just in one
// process's memory. It fails open (lets the request through) if Redis is
// unreachable — a rate limiter should never be the reason the whole app goes
// down when the cache does.
func rateLimitMiddleware(rdb *redis.Client, bucket string, limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if rdb == nil {
			ctx.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:%s:%s", bucket, ctx.ClientIP())

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			slog.Warn("rate limiter: redis unavailable, failing open", "bucket", bucket, "error", err)
			ctx.Next()
			return
		}
		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			ttl, err := rdb.TTL(ctx, key).Result()
			if err == nil && ttl > 0 {
				ctx.Header("Retry-After", fmt.Sprintf("%.0f", ttl.Seconds()))
			}
			handleAbort(ctx, domain.ErrTooManyRequests)
			return
		}

		ctx.Next()
	}
}

