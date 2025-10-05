package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// User represents a minimal user record for verification state.
// Stored in Redis as JSON at key user:{userID}
//
type User struct {
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	EmailVerified  bool      `json:"email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at,omitempty"`
}

type UserRepository struct {
	client *redis.Client
}

func NewUserRepository(client *redis.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) key(userID string) string { return fmt.Sprintf("user:%s", userID) }

// Upsert sets/updates user fields
func (r *UserRepository) Upsert(ctx context.Context, user User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(user.UserID), b, 0).Err()
}

// Get retrieves a user; returns (nil, nil) if not found
func (r *UserRepository) Get(ctx context.Context, userID string) (*User, error) {
	val, err := r.client.Get(ctx, r.key(userID)).Bytes()
	if err != nil {
		if err == redis.Nil { return nil, nil }
		return nil, err
	}
	var u User
	if err := json.Unmarshal(val, &u); err != nil { return nil, err }
	return &u, nil
}

// SetEmailVerified marks the user as verified and persists
func (r *UserRepository) SetEmailVerified(ctx context.Context, userID, email string, verifiedAt time.Time) error {
	u, err := r.Get(ctx, userID)
	if err != nil { return err }
	if u == nil {
		u = &User{UserID: userID, Email: email}
	}
	u.Email = email
	u.EmailVerified = true
	u.EmailVerifiedAt = verifiedAt
	return r.Upsert(ctx, *u)
}
