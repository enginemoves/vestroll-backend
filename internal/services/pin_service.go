package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"regexp"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
)

// PINService encapsulates PIN setup and authentication logic
type PINService struct {
	repo *repository.PinRepository
}

func NewPINService(repo *repository.PinRepository) *PINService {
	return &PINService{repo: repo}
}

// ValidatePINFormat ensures PIN is 4-6 digits numeric
func (s *PINService) ValidatePINFormat(pin string) error {
	if len(pin) < 4 || len(pin) > 6 {
		return errors.New("pin must be 4 to 6 digits")
	}
	matched, _ := regexp.MatchString("^[0-9]+$", pin)
	if !matched {
		return errors.New("pin must contain only digits")
	}
	return nil
}

func generateSalt(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashPIN(pin, salt string) string {
	sum := sha256.Sum256([]byte(pin + ":" + salt))
	return hex.EncodeToString(sum[:])
}

// SetupPIN validates and stores the user's PIN (salted+hashed)
func (s *PINService) SetupPIN(ctx context.Context, req models.SetupPINRequest) error {
	if req.UserID == "" {
		return errors.New("user_id is required")
	}
	if err := s.ValidatePINFormat(req.PIN); err != nil {
		return err
	}
	salt, err := generateSalt(16)
	if err != nil {
		return err
	}
	hash := hashPIN(req.PIN, salt)
	data := models.PinData{Salt: salt, Hash: hash, CreatedAt: time.Now()}
	return s.repo.Save(ctx, req.UserID, data)
}

// LoginPIN verifies the provided PIN against stored hash
func (s *PINService) LoginPIN(ctx context.Context, req models.LoginPINRequest) error {
	if req.UserID == "" {
		return errors.New("user_id is required")
	}
	if err := s.ValidatePINFormat(req.PIN); err != nil {
		return err
	}
	stored, err := s.repo.Get(ctx, req.UserID)
	if err != nil {
		return errors.New("pin not set")
	}
	candidate := hashPIN(req.PIN, stored.Salt)
	if subtle.ConstantTimeCompare([]byte(candidate), []byte(stored.Hash)) != 1 {
		return errors.New("invalid pin")
	}
	return nil
}