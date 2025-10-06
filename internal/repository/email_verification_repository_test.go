package repository

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/go-redis/redis/v8"
)

type testRedis struct {
	srv    *miniredis.Miniredis
	client *redis.Client
}

func newTestRedis(t *testing.T) *testRedis {
	t.Helper()
	s, err := miniredis.Run()
	if err != nil { t.Fatalf("miniredis: %v", err) }
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	return &testRedis{srv: s, client: client}
}

func (tr *testRedis) Close() {
	_ = tr.client.Close()
	tr.srv.Close()
}

func TestEmailVerificationRepository_StoreGetDelete(t *testing.T) {
	tr := newTestRedis(t)
	defer tr.Close()

	repo := NewEmailVerificationRepository(tr.client, time.Minute)
	ctx := context.Background()

	token := "token123"
	payload := models.EmailVerificationTokenPayload{UserID: "u1", Email: "u1@example.com", IssuedAt: time.Now()}

	// Store
	if err := repo.StoreToken(ctx, token, payload); err != nil {
		t.Fatalf("StoreToken error: %v", err)
	}

	// Get
	got, err := repo.GetToken(ctx, token)
	if err != nil { t.Fatalf("GetToken error: %v", err) }
	if got == nil || got.UserID != payload.UserID || got.Email != payload.Email {
		t.Fatalf("GetToken mismatch: got=%+v want=%+v", got, payload)
	}

	// Delete
	if err := repo.DeleteToken(ctx, token); err != nil { t.Fatalf("DeleteToken error: %v", err) }
	got, err = repo.GetToken(ctx, token)
	if err != nil { t.Fatalf("GetToken after delete error: %v", err) }
	if got != nil { t.Fatalf("expected nil after delete, got=%+v", got) }
}

func TestEmailVerificationRepository_TTLExpiry(t *testing.T) {
	tr := newTestRedis(t)
	defer tr.Close()

	repo := NewEmailVerificationRepository(tr.client, time.Minute)
	ctx := context.Background()

	token := "expiring"
	payload := models.EmailVerificationTokenPayload{UserID: "u2", Email: "u2@example.com", IssuedAt: time.Now()}
	if err := repo.StoreToken(ctx, token, payload); err != nil {
		t.Fatalf("StoreToken: %v", err)
	}
	// Advance time beyond TTL
	tr.srv.FastForward(time.Minute + time.Second)
	got, err := repo.GetToken(ctx, token)
	if err == nil && got != nil {
		t.Fatalf("expected expired token to be nil, got=%+v", got)
	}
}