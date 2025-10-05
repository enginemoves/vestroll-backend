package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
)

// EmailVerificationService coordinates token generation, email sending and verification
//
type EmailVerificationService struct {
	repo         *repository.EmailVerificationRepository
	userRepo     *repository.UserRepository
	emailService interface{
		IsConfigured() bool
		SendVerificationEmail(ctx context.Context, email, token, linkBase string) error
	}
	cfg config.EmailVerificationConfig
}

func NewEmailVerificationService(repo *repository.EmailVerificationRepository, userRepo *repository.UserRepository, emailSvc interface{
	IsConfigured() bool
	SendVerificationEmail(ctx context.Context, email, token, linkBase string) error
}, cfg config.EmailVerificationConfig) *EmailVerificationService {
	return &EmailVerificationService{repo: repo, userRepo: userRepo, emailService: emailSvc, cfg: cfg}
}

// GenerateAndSend creates a token for the user and sends an email
func (s *EmailVerificationService) GenerateAndSend(ctx context.Context, userID, email string) (string, error) {
	token, err := generateSecureToken(32)
	if err != nil { return "", err }
	payload := models.EmailVerificationTokenPayload{UserID: userID, Email: email, IssuedAt: time.Now()}
	if err := s.repo.StoreToken(ctx, token, payload); err != nil { return "", err }
	if s.emailService != nil && s.emailService.IsConfigured() {
		if err := s.emailService.SendVerificationEmail(ctx, email, token, s.cfg.LinkBaseURL); err != nil { return "", err }
	}
	return token, nil
}

// Verify consumes the token and sets the user's verification state
func (s *EmailVerificationService) Verify(ctx context.Context, token string) error {
	payload, err := s.repo.GetToken(ctx, token)
	if err != nil { return err }
	if payload == nil { return ErrInvalidOrExpiredToken }
	// Mark verified
	if err := s.userRepo.SetEmailVerified(ctx, payload.UserID, payload.Email, time.Now()); err != nil { return err }
	// Delete token to prevent reuse
	_ = s.repo.DeleteToken(ctx, token)
	return nil
}

func generateSecureToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil { return "", err }
	return hex.EncodeToString(b), nil
}

// Errors
var ErrInvalidOrExpiredToken = &verificationError{msg: "invalid_or_expired_token"}

type verificationError struct{ msg string }
func (e *verificationError) Error() string { return e.msg }
