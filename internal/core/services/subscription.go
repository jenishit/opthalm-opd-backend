package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
	"github.com/redis/go-redis/v9"
)

const subscriptionCacheTTL = 45 * time.Second

type SubscriptionService struct {
	repo  port.SubscriptionRepository
	redis *redis.Client
}

func NewSubscriptionService(
	r port.SubscriptionRepository,
	rc *redis.Client,
) *SubscriptionService {
	return &SubscriptionService{
		repo:  r,
		redis: rc,
	}
}

func (s *SubscriptionService) cacheKey(clinicID uuid.UUID) string {
