package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type PasswordResetRepository struct {
	client *redis.Client
	resetTTL time.Duration
}

func NewPasswordResetRepository(client *redis.Client, resetTTL time.Duration) *PasswordResetRepository {
	return &PasswordResetRepository{
		client: client,
		resetTTL: resetTTL,
	}
}

func (r *PasswordResetRepository) StoreResetCode(ctx context.Context, identifier, code string) error {
	key := fmt.Sprintf("password_reset:%s", identifier)
	data, _ := json.Marshal(map[string]interface{}{
		"code": code,
		"created_at": time.Now().Unix(),
	})
	return r.client.Set(ctx, key, data, r.resetTTL).Err()
}

func (r *PasswordResetRepository) GetResetCode(ctx context.Context, identifier string) (string, error) {
	key := fmt.Sprintf("password_reset:%s", identifier)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return "", err
	}
	code, ok := result["code"].(string)
	if !ok {
		return "", fmt.Errorf("reset code not found")
	}
	return code, nil
}

func (r *PasswordResetRepository) DeleteResetCode(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("password_reset:%s", identifier)
	return r.client.Del(ctx, key).Err()
}
