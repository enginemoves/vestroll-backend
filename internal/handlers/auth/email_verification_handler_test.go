package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/handlers/auth"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// fakeEmailSvc implements the minimal interface used by EmailVerificationService
// for sending emails.
type fakeEmailSvc struct{
	configured bool
}

func (f *fakeEmailSvc) IsConfigured() bool { return f.configured }
func (f *fakeEmailSvc) SendVerificationEmail(ctx context.Context, email, token, link string) error { return nil }

func setupHandlerTest(t *testing.T) (*gin.Engine, *services.EmailVerificationService, *repository.EmailVerificationRepository, *repository.UserRepository, func()) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	srv, err := miniredis.Run()
	if err != nil { t.Fatalf("miniredis: %v", err) }
	client := redis.NewClient(&redis.Options{Addr: srv.Addr()})

	emailRepo := repository.NewEmailVerificationRepository(client, time.Minute)
	userRepo := repository.NewUserRepository(client)
	fake := &fakeEmailSvc{configured: true}
	cfg := config.EmailVerificationConfig{TTL: time.Minute}
	svc := services.NewEmailVerificationService(emailRepo, userRepo, fake, cfg)

	h := auth.NewEmailVerificationHandler(svc)
	r := gin.New()
	r.POST("/api/v1/auth/verify-email", h.VerifyEmail)

	cleanup := func(){ client.Close(); srv.Close() }
	return r, svc, emailRepo, userRepo, cleanup
}

type verifyReq struct {
	Token string `json:"token"`
}

func doRequest(r *gin.Engine, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-email", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestVerifyEmail_Success(t *testing.T) {
	r, svc, _, userRepo, cleanup := setupHandlerTest(t)
	defer cleanup()

	// Generate token first (simulates registration flow)
	token, err := svc.GenerateAndSend(context.Background(), "user-42", "user42@example.com")
	if err != nil { t.Fatalf("GenerateAndSend: %v", err) }

	rec := doRequest(r, verifyReq{Token: token})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	u, err := userRepo.Get(context.Background(), "user-42")
	if err != nil { t.Fatalf("userRepo.Get: %v", err) }
	if u == nil || !u.EmailVerified {
		t.Fatalf("expected verified user, got=%+v", u)
	}
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	r, _, _, _, cleanup := setupHandlerTest(t)
	defer cleanup()

	rec := doRequest(r, verifyReq{Token: "nope"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestVerifyEmail_MissingToken(t *testing.T) {
	r, _, _, _, cleanup := setupHandlerTest(t)
	defer cleanup()

	// Missing token field
	rec := doRequest(r, map[string]any{"not_token": "x"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rec.Code, rec.Body.String())
	}
}
