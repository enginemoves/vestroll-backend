package password_reset_service

import (
	"context"
	"crypto/rand"
	"time"
	"math/big"

	"github.com/codeZe-us/vestroll-backend/internal/repository"
	"github.com/codeZe-us/vestroll-backend/internal/services/email_service"
	"github.com/codeZe-us/vestroll-backend/internal/services/sms_service"
)

type PasswordResetService struct {
	Repo         *repository.PasswordResetRepository
	EmailService *email_service.EmailService
	SMSService   *sms_service.SMSService
	TTL          time.Duration
}

func NewPasswordResetService(repo *repository.PasswordResetRepository, emailSvc *email_service.EmailService, smsSvc *sms_service.SMSService, ttl time.Duration) *PasswordResetService {
	return &PasswordResetService{
		Repo:         repo,
		EmailService: emailSvc,
		SMSService:   smsSvc,
		TTL:          ttl,
	}
}

func (s *PasswordResetService) GenerateAndSendResetCode(ctx context.Context, identifier, channel string) (string, error) {
	code, err := generateCode(6)
	if err != nil {
		return "", err
	}
	if err := s.Repo.StoreResetCode(ctx, identifier, code); err != nil {
		return "", err
	}
	if channel == "email" {
		if err := s.EmailService.SendOTP(ctx, identifier, code); err != nil {
			return "", err
		}
	} else if channel == "sms" {
		if err := s.SMSService.SendOTP(ctx, identifier, code); err != nil {
			return "", err
		}
	}
	return code, nil
}

func (s *PasswordResetService) VerifyResetCode(ctx context.Context, identifier, code string) (bool, error) {
	stored, err := s.Repo.GetResetCode(ctx, identifier)
	if err != nil {
		return false, err
	}
	return stored == code, nil
}

func (s *PasswordResetService) DeleteResetCode(ctx context.Context, identifier string) error {
	return s.Repo.DeleteResetCode(ctx, identifier)
}

func generateCode(length int) (string, error) {
	numbers := "0123456789"
	code := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbers))) )
		if err != nil {
			return "", err
		}
		code += string(numbers[num.Int64()])
	}
	return code, nil
}
