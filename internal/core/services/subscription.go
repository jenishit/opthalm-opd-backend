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
	return "subscription:" + clinicID.String()
}

func (s *SubscriptionService) GetByClinicID(
	ctx context.Context,
	clinicID uuid.UUID,
) (*domain.Subscription, error) {
	key := s.cacheKey(clinicID)

	// Cache lookup
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var sub domain.Subscription
			if json.Unmarshal([]byte(cached), &sub) == nil {
				return &sub, nil
			}
		} else if err != redis.Nil {
			// Ignore Redis errors and fall back to DB
		}
	}

	// DB fallback
	sub, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	// Populate cache
	if s.redis != nil && sub != nil {
		if b, err := json.Marshal(sub); err == nil {
			_ = s.redis.Set(ctx, key, b, subscriptionCacheTTL).Err()
		}
	}

	return sub, nil
}

func (s *SubscriptionService) Upsert(
	ctx context.Context,
	sub *domain.Subscription,
) (*domain.Subscription, error) {
	if sub.PlanName == "" {
		sub.PlanName = "trial"
	}
	if sub.Status == "" {
		sub.Status = domain.SubscriptionTrialing
	}

	updated, err := s.repo.Upsert(ctx, sub)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.redis != nil {
		_ = s.redis.Del(ctx, s.cacheKey(updated.ClinicID)).Err()
	}

	return updated, nil
}
