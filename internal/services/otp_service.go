package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
)

type OTPService struct {
	otpRepo      *repository.OTPRepository
	smsService   *SMSService
	emailService *EmailService
	config       config.OTPConfig
}

func NewOTPService(
	otpRepo *repository.OTPRepository,
	smsService *SMSService,
	emailService *EmailService,
	config config.OTPConfig,
) *OTPService {
	return &OTPService{
		otpRepo:      otpRepo,
		smsService:   smsService,
		emailService: emailService,
		config:       config,
	}
}

func (s *OTPService) SendOTP(ctx context.Context, req models.OTPRequest) error {
	// Validate identifier format
	if err := s.validateIdentifier(req.Identifier, req.Type); err != nil {
		return err
	}

	// Check rate limiting
	allowed, err := s.otpRepo.CheckRateLimit(
		ctx,
		req.Identifier,
		s.config.RateLimit.MaxRequests,
		s.config.RateLimit.WindowSize,
	)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	if !allowed {
		remaining, _ := s.otpRepo.GetRemainingAttempts(ctx, req.Identifier, s.config.RateLimit.MaxRequests)
		return fmt.Errorf("rate limit exceeded. Try again later. Remaining attempts: %d", remaining)
	}

	// Generate OTP code
	code, err := s.generateOTPCode(s.config.Length)
	if err != nil {
		return fmt.Errorf("failed to generate OTP code: %w", err)
	}

	// Create OTP data
	otpData := models.OTPData{
		Code:      code,
		Type:      req.Type,
		ExpiresAt: time.Now().Add(s.config.TTL),
		Attempts:  0,
	}

	// Store OTP in Redis
	if err := s.otpRepo.StoreOTP(ctx, req.Identifier, otpData); err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	// Send OTP via appropriate channel
	switch req.Type {
	case models.OTPTypeSMS:
		if !s.smsService.IsConfigured() {
			return fmt.Errorf("SMS service is not configured")
		}
		if err := s.smsService.SendOTP(ctx, req.Identifier, code); err != nil {
			// Clean up stored OTP on send failure
			s.otpRepo.DeleteOTP(ctx, req.Identifier, req.Type)
			return fmt.Errorf("failed to send SMS OTP: %w", err)
		}
	case models.OTPTypeEmail:
		if !s.emailService.IsConfigured() {
			return fmt.Errorf("email service is not configured")
		}
		if err := s.emailService.SendOTP(ctx, req.Identifier, code); err != nil {
			// Clean up stored OTP on send failure
			s.otpRepo.DeleteOTP(ctx, req.Identifier, req.Type)
			return fmt.Errorf("failed to send email OTP: %w", err)
		}
	default:
		return fmt.Errorf("unsupported OTP type: %s", req.Type)
	}

	return nil
}

func (s *OTPService) VerifyOTP(ctx context.Context, req models.OTPVerificationRequest) error {
	// Validate identifier format
	if err := s.validateIdentifier(req.Identifier, req.Type); err != nil {
		return err
	}

	// Get stored OTP
	otpData, err := s.otpRepo.GetOTP(ctx, req.Identifier, req.Type)
	if err != nil {
		return fmt.Errorf("failed to retrieve OTP: %w", err)
	}
	if otpData == nil {
		return fmt.Errorf("OTP not found or expired")
	}

	// Check if OTP has expired
	if time.Now().After(otpData.ExpiresAt) {
		s.otpRepo.DeleteOTP(ctx, req.Identifier, req.Type)
		return fmt.Errorf("OTP has expired")
	}

	// Check attempt limits (max 3 attempts)
	if otpData.Attempts >= 3 {
		s.otpRepo.DeleteOTP(ctx, req.Identifier, req.Type)
		return fmt.Errorf("maximum verification attempts exceeded")
	}

	// Verify the code
	if otpData.Code != req.Code {
		// Increment attempt counter
		if err := s.otpRepo.IncrementAttempts(ctx, req.Identifier, req.Type); err != nil {
			return fmt.Errorf("failed to update attempt counter: %w", err)
		}
		return fmt.Errorf("invalid OTP code")
	}

	// OTP is valid, delete it to prevent reuse
	if err := s.otpRepo.DeleteOTP(ctx, req.Identifier, req.Type); err != nil {
		return fmt.Errorf("failed to clean up OTP: %w", err)
	}

	return nil
}

func (s *OTPService) generateOTPCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid OTP length")
	}

	// Generate random digits
	code := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}

	return code, nil
}

func (s *OTPService) validateIdentifier(identifier string, otpType models.OTPType) error {
	switch otpType {
	case models.OTPTypeSMS:
		// Basic phone number validation (should start with +)
		phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
		if !phoneRegex.MatchString(identifier) {
			return fmt.Errorf("invalid phone number format. Must be in international format (e.g., +1234567890)")
		}
	case models.OTPTypeEmail:
		// Basic email validation
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(identifier) {
			return fmt.Errorf("invalid email address format")
		}
	default:
		return fmt.Errorf("unsupported OTP type: %s", otpType)
	}
	return nil
}