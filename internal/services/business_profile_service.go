package services

import (
	"context"
	"strings"
	"time"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/repository"
)

// BusinessProfileService handles validation and persistence of business profiles
// Only contractor account type is currently supported
// Completion tracking ensures required fields are provided

type BusinessProfileService struct {
	repo *repository.BusinessProfileRepository
}

func NewBusinessProfileService(repo *repository.BusinessProfileRepository) *BusinessProfileService {
	return &BusinessProfileService{repo: repo}
}

// ValidateContractor checks required fields for contractor business profiles
func (s *BusinessProfileService) ValidateContractor(req models.BusinessDetailsRequest) error {
	// Additional business rules can be applied here
	// e.g., registration number pattern, tax ID format per country
	if strings.TrimSpace(req.BusinessName) == "" {
		return ErrValidation("business_name is required")
	}
	if strings.TrimSpace(req.RegistrationNumber) == "" {
		return ErrValidation("registration_number is required")
	}
	if strings.TrimSpace(req.TaxID) == "" {
		return ErrValidation("tax_id is required")
	}
	return nil
}

// BuildProfile builds the persistence model and calculates completion
func (s *BusinessProfileService) BuildProfile(req models.BusinessDetailsRequest) models.BusinessProfile {
	completedFields := 0
	requiredFields := 8 // business_name, registration_number, tax_id, address.street, address.city, address.country, contact.email, contact.phone

	if req.BusinessName != "" { completedFields++ }
	if req.RegistrationNumber != "" { completedFields++ }
	if req.TaxID != "" { completedFields++ }
	if req.Address.Street != "" { completedFields++ }
	if req.Address.City != "" { completedFields++ }
	if req.Address.Country != "" { completedFields++ }
	if req.Contact.Email != "" { completedFields++ }
	if req.Contact.Phone != "" { completedFields++ }

	percent := int(float64(completedFields) / float64(requiredFields) * 100)

	return models.BusinessProfile{
		UserID:             req.UserID,
		AccountType:        req.AccountType,
		BusinessName:       req.BusinessName,
		RegistrationNumber: req.RegistrationNumber,
		TaxID:              req.TaxID,
		Address:            req.Address,
		Contact:            req.Contact,
		Completed:          percent == 100,
		CompletionPercent:  percent,
		UpdatedAt:          time.Now(),
	}
}

// Save persists the business profile
func (s *BusinessProfileService) Save(ctx context.Context, profile models.BusinessProfile) error {
	return s.repo.Save(ctx, profile)
}

// ErrValidation creates a simple validation error
func ErrValidation(msg string) error { return &validationError{msg: msg} }

type validationError struct{ msg string }

func (e *validationError) Error() string { return "validation_error: " + e.msg }