package repository

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

type testRedis2 struct {
	srv    *miniredis.Miniredis
	client *redis.Client
}

func newTestRedis2(t *testing.T) *testRedis2 {
	t.Helper()
	s, err := miniredis.Run()
	if err != nil { t.Fatalf("miniredis: %v", err) }
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	return &testRedis2{srv: s, client: client}
}

func (tr *testRedis2) Close() {
	_ = tr.client.Close()
	tr.srv.Close()
}

func TestUserRepository_UpsertGetAndVerify(t *testing.T) {
	tr := newTestRedis2(t)
	defer tr.Close()

	repo := NewUserRepository(tr.client)
	ctx := context.Background()

	u := User{UserID: "u123", Email: "test@example.com", EmailVerified: false}
	if err := repo.Upsert(ctx, u); err != nil { t.Fatalf("Upsert: %v", err) }

	got, err := repo.Get(ctx, "u123")
	if err != nil { t.Fatalf("Get: %v", err) }
	if got == nil || got.Email != u.Email || got.EmailVerified {
		t.Fatalf("Get mismatch: got=%+v", got)
	}

	verAt := time.Now()
	if err := repo.SetEmailVerified(ctx, "u123", "test@example.com", verAt); err != nil {
		t.Fatalf("SetEmailVerified: %v", err)
	}
	got, err = repo.Get(ctx, "u123")
	if err != nil { t.Fatalf("Get2: %v", err) }
	if got == nil || !got.EmailVerified {
		t.Fatalf("expected verified user, got=%+v", got)
	}
}