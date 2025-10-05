package models

import "time"

// EmailVerificationTokenPayload is the data stored with a verification token
// It lets us verify and update the correct user after a token is submitted
// without requiring the client to send user_id explicitly.
//
type EmailVerificationTokenPayload struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
}

// VerifyEmailRequest is the payload for POST /api/v1/auth/verify-email
// Client submits the token received via email
//
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmailResponse is the success response
//
type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
