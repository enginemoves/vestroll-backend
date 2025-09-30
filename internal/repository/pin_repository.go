package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/go-redis/redis/v8"
)

// PinRepository handles storage of user PINs in Redis
type PinRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPinRepository(client *redis.Client, ttl time.Duration) *PinRepository {
	return &PinRepository{client: client, ttl: ttl}
}

func (r *PinRepository) key(userID string) string {
	return fmt.Sprintf("user_pin:%s", userID)
}

// Save stores the hashed PIN and salt for a user
func (r *PinRepository) Save(ctx context.Context, userID string, data models.PinData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(userID), string(b), r.ttl).Err()
}

// Get retrieves the stored PIN data for a user
func (r *PinRepository) Get(ctx context.Context, userID string) (models.PinData, error) {
	val, err := r.client.Get(ctx, r.key(userID)).Result()
	if err != nil {
		return models.PinData{}, err
	}
	var data models.PinData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return models.PinData{}, err
	}
	return data, nil
}