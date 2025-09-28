package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/go-redis/redis/v8"
)

type OTPRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewOTPRepository(client *redis.Client, ttl time.Duration) *OTPRepository {
	return &OTPRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *OTPRepository) StoreOTP(ctx context.Context, identifier string, otpData models.OTPData) error {
	key := fmt.Sprintf("otp:%s:%s", string(otpData.Type), identifier)
	
	data, err := json.Marshal(otpData)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *OTPRepository) GetOTP(ctx context.Context, identifier string, otpType models.OTPType) (*models.OTPData, error) {
	key := fmt.Sprintf("otp:%s:%s", string(otpType), identifier)
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // OTP not found or expired
		}
		return nil, err
	}

	var otpData models.OTPData
	if err := json.Unmarshal([]byte(data), &otpData); err != nil {
		return nil, err
	}

	return &otpData, nil
}

func (r *OTPRepository) DeleteOTP(ctx context.Context, identifier string, otpType models.OTPType) error {
	key := fmt.Sprintf("otp:%s:%s", string(otpType), identifier)
	return r.client.Del(ctx, key).Err()
}

func (r *OTPRepository) IncrementAttempts(ctx context.Context, identifier string, otpType models.OTPType) error {
	otpData, err := r.GetOTP(ctx, identifier, otpType)
	if err != nil || otpData == nil {
		return err
	}

	otpData.Attempts++
	return r.StoreOTP(ctx, identifier, *otpData)
}

// Rate limiting methods
func (r *OTPRepository) CheckRateLimit(ctx context.Context, identifier string, maxRequests int, windowSize time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:otp:%s", identifier)
	
	current, err := r.client.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if current >= maxRequests {
		return false, nil // Rate limit exceeded
	}

	// Increment counter
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, windowSize)
	_, err = pipe.Exec(ctx)
	
	return err == nil, err
}

func (r *OTPRepository) GetRemainingAttempts(ctx context.Context, identifier string, maxRequests int) (int, error) {
	key := fmt.Sprintf("rate_limit:otp:%s", identifier)
	
	current, err := r.client.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return maxRequests, err
	}

	remaining := maxRequests - current
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}