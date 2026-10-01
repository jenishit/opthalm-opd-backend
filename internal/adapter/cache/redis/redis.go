package redis

import (
	"context"
	"fmt"

	"github.com/jenishit/opthalm-opd-backend/internal/adapter/config"
	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	return client, nil
}
