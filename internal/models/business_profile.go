package models

import "time"

// AccountType represents the type of account
const (
	AccountTypeContractor string = "contractor"
)

// BusinessAddress holds business location details
type BusinessAddress struct {
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country" binding:"required"`
}

// BusinessContact holds contact details for the business
type BusinessContact struct {
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required,min=7,max=20"`
}

// BusinessDetailsRequest is the payload for creating/updating business details
// For now, only contractor account type is supported
type BusinessDetailsRequest struct {
	UserID             string           `json:"user_id" binding:"required"`
	AccountType        string           `json:"account_type" binding:"required,oneof=contractor"`
	BusinessName       string           `json:"business_name" binding:"required"`
	RegistrationNumber string           `json:"registration_number" binding:"required"`
	TaxID              string           `json:"tax_id" binding:"required"`
	Address            BusinessAddress  `json:"address" binding:"required"`
	Contact            BusinessContact  `json:"contact" binding:"required"`
}

// BusinessProfile persists the normalized business details
type BusinessProfile struct {
	UserID             string          `json:"user_id"`
	AccountType        string          `json:"account_type"`
	BusinessName       string          `json:"business_name"`
	RegistrationNumber string          `json:"registration_number"`
	TaxID              string          `json:"tax_id"`
	Address            BusinessAddress `json:"address"`
	Contact            BusinessContact `json:"contact"`
	Completed          bool            `json:"completed"`
	CompletionPercent  int             `json:"completion_percent"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// BusinessProfileResponse is returned by the endpoint
type BusinessProfileResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Profile BusinessProfile `json:"profile"`
}