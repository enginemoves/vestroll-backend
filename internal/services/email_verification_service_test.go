package services

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/go-redis/redis/v8"
)

type fakeEmailSvc struct{
	configured bool
	sent []struct{ email, token, link string }
}

func (f *fakeEmailSvc) IsConfigured() bool { return f.configured }
func (f *fakeEmailSvc) SendVerificationEmail(ctx context.Context, email, token, link string) error {
	f.sent = append(f.sent, struct{ email, token, link string }{email, token, link})
	return nil
}

func newRedisClient(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	s, err := miniredis.Run()
	if err != nil { t.Fatalf("miniredis: %v", err) }
	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	return s, client
}

func TestEmailVerificationService_GenerateAndSend_And_Verify(t *testing.T) {
	srv, client := newRedisClient(t)
	defer func(){ client.Close(); srv.Close() }()

	repo := repository.NewEmailVerificationRepository(client, time.Minute)
	userRepo := repository.NewUserRepository(client)
	fake := &fakeEmailSvc{configured: true}
	cfg := config.EmailVerificationConfig{TTL: time.Minute, LinkBaseURL: "https://app.example.com/verify"}

	svc := NewEmailVerificationService(repo, userRepo, fake, cfg)
	ctx := context.Background()

	// Generate
	token, err := svc.GenerateAndSend(ctx, "user-1", "user1@example.com")
	if err != nil { t.Fatalf("GenerateAndSend: %v", err) }
	if token == "" { t.Fatalf("expected non-empty token") }
	if len(fake.sent) != 1 { t.Fatalf("expected 1 email sent, got %d", len(fake.sent)) }
	if fake.sent[0].email != "user1@example.com" || fake.sent[0].token != token { t.Fatalf("email content mismatch: %+v", fake.sent[0]) }

	// Verify
	if err := svc.Verify(ctx, token); err != nil { t.Fatalf("Verify: %v", err) }
	// User should be marked verified
	u, err := userRepo.Get(ctx, "user-1")
	if err != nil { t.Fatalf("user Get: %v", err) }
	if u == nil || !u.EmailVerified { t.Fatalf("expected verified user, got=%+v", u) }
	// Token should be deleted
	p, err := repo.GetToken(ctx, token)
	if err != nil { t.Fatalf("GetToken: %v", err) }
	if p != nil { t.Fatalf("expected token deleted, got=%+v", p) }
}

func TestEmailVerificationService_InvalidToken(t *testing.T) {
	srv, client := newRedisClient(t)
	defer func(){ client.Close(); srv.Close() }()

	repo := repository.NewEmailVerificationRepository(client, time.Minute)
	userRepo := repository.NewUserRepository(client)
	fake := &fakeEmailSvc{configured: true}
	cfg := config.EmailVerificationConfig{TTL: time.Minute}
	svc := NewEmailVerificationService(repo, userRepo, fake, cfg)

	if err := svc.Verify(context.Background(), "does-not-exist"); err == nil {
		t.Fatalf("expected error for invalid token")
	}
}