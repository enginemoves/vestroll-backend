package models

type PasswordResetRequest struct {
	Identifier string `json:"identifier" binding:"required"` // email or phone
}

type VerifyResetCodeRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Code       string `json:"code" binding:"required,len=6"`
}

type ResetPasswordRequest struct {
	Identifier      string `json:"identifier" binding:"required"`
	Code            string `json:"code" binding:"required,len=6"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// Response types

type PasswordResetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
