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
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	email_service "github.com/codeZe-us/vestroll-backend/internal/services/email_service"
	sms_service "github.com/codeZe-us/vestroll-backend/internal/services/sms_service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func setupHandlerTest(t *testing.T) (*gin.Engine, *repository.OTPRepository, *repository.UserRepository, *services.OTPService, func()) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	srv, err := miniredis.Run()
	if err != nil { t.Fatalf("miniredis: %v", err) }
	client := redis.NewClient(&redis.Options{Addr: srv.Addr()})

	otpRepo := repository.NewOTPRepository(client, time.Minute)
	userRepo := repository.NewUserRepository(client)
	otpCfg := config.OTPConfig{Length: 6, TTL: time.Minute, RateLimit: config.RateLimitConfig{MaxRequests: 5, WindowSize: time.Minute}}
	otpSvc := services.NewOTPService(otpRepo, sms_service.NewSMSService(config.TwilioConfig{}), email_service.NewEmailService(config.SMTPConfig{}), otpCfg)

	h := auth.NewEmailVerificationHandler(otpSvc, userRepo)
	r := gin.New()
	r.POST("/api/auth/verify-email", h.VerifyEmail)

	cleanup := func(){ client.Close(); srv.Close() }
	return r, otpRepo, userRepo, otpSvc, cleanup
}

type verifyReq struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Code   string `json:"code"`
}

func doRequest(r *gin.Engine, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/verify-email", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestVerifyEmail_Success(t *testing.T) {
	r, otpRepo, userRepo, _, cleanup := setupHandlerTest(t)
	defer cleanup()

	// Seed OTP
	code := "123456"
	email := "user42@example.com"
if err := otpRepo.StoreOTP(context.Background(), email, models.OTPData{Code: code, Type: models.OTPTypeEmail, ExpiresAt: time.Now().Add(time.Minute), Attempts: 0}); err != nil {
		t.Fatalf("StoreOTP: %v", err)
	}

	rec := doRequest(r, verifyReq{UserID: "user-42", Email: email, Code: code})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}
u, err := userRepo.Get(context.Background(), "user-42")
	if err != nil { t.Fatalf("userRepo.Get: %v", err) }
	if u == nil || !u.EmailVerified { t.Fatalf("expected verified user, got=%+v", u) }
}

func TestVerifyEmail_InvalidCode(t *testing.T) {
	r, otpRepo, _, _, cleanup := setupHandlerTest(t)
	defer cleanup()

	email := "user43@example.com"
if err := otpRepo.StoreOTP(context.Background(), email, models.OTPData{Code: "111111", Type: models.OTPTypeEmail, ExpiresAt: time.Now().Add(time.Minute), Attempts: 0}); err != nil {
		t.Fatalf("StoreOTP: %v", err)
	}
	rec := doRequest(r, verifyReq{UserID: "user-43", Email: email, Code: "222222"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestVerifyEmail_NonDigitCode(t *testing.T) {
	r, _, _, _, cleanup := setupHandlerTest(t)
	defer cleanup()

	rec := doRequest(r, verifyReq{UserID: "user-44", Email: "user44@example.com", Code: "12ab56"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

