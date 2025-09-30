package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/go-redis/redis/v8"
)

// BusinessProfileRepository stores business profiles in Redis
// Key pattern: business_profile:{user_id}
type BusinessProfileRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewBusinessProfileRepository(client *redis.Client, ttl time.Duration) *BusinessProfileRepository {
	return &BusinessProfileRepository{client: client, ttl: ttl}
}

func (r *BusinessProfileRepository) Save(ctx context.Context, profile models.BusinessProfile) error {
	key := fmt.Sprintf("business_profile:%s", profile.UserID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *BusinessProfileRepository) Get(ctx context.Context, userID string) (*models.BusinessProfile, error) {
	key := fmt.Sprintf("business_profile:%s", userID)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var profile models.BusinessProfile
	if err := json.Unmarshal(val, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}