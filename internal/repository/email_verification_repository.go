package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/go-redis/redis/v8"
)

// EmailVerificationRepository stores verification tokens in Redis
// Key: email_verify_token:{token} => JSON payload (user_id, email, issued_at)
// TTL: configured
//
type EmailVerificationRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewEmailVerificationRepository(client *redis.Client, ttl time.Duration) *EmailVerificationRepository {
	return &EmailVerificationRepository{client: client, ttl: ttl}
}

func (r *EmailVerificationRepository) tokenKey(token string) string {
	return fmt.Sprintf("email_verify_token:%s", token)
}

func (r *EmailVerificationRepository) StoreToken(ctx context.Context, token string, payload models.EmailVerificationTokenPayload) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.tokenKey(token), b, r.ttl).Err()
}

func (r *EmailVerificationRepository) GetToken(ctx context.Context, token string) (*models.EmailVerificationTokenPayload, error) {
	val, err := r.client.Get(ctx, r.tokenKey(token)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var payload models.EmailVerificationTokenPayload
	if err := json.Unmarshal(val, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (r *EmailVerificationRepository) DeleteToken(ctx context.Context, token string) error {
	return r.client.Del(ctx, r.tokenKey(token)).Err()
}
