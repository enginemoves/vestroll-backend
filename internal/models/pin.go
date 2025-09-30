package models

import "time"

// SetupPINRequest is the payload to set up a user's PIN
type SetupPINRequest struct {
	UserID string `json:"user_id" binding:"required"`
	PIN    string `json:"pin" binding:"required"`
}

// LoginPINRequest is the payload to authenticate a user using PIN
type LoginPINRequest struct {
	UserID string `json:"user_id" binding:"required"`
	PIN    string `json:"pin" binding:"required"`
}

// PinData represents the stored PIN data (salt + hashed PIN)
type PinData struct {
	Salt      string    `json:"salt"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

// SetupPINResponse is returned by the setup PIN endpoint
type SetupPINResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginPINResponse is returned by the login PIN endpoint
type LoginPINResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}