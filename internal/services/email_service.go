package services

import (
	"context"
	"fmt"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer    *gomail.Dialer
	fromEmail string
	fromName  string
}

func NewEmailService(cfg config.SMTPConfig) *EmailService {
	if cfg.Username == "" || cfg.Password == "" {
		return &EmailService{} // Return unconfigured service
	}

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	
	return &EmailService{
		dialer:    dialer,
		fromEmail: cfg.FromEmail,
		fromName:  cfg.FromName,
	}
}

func (e *EmailService) SendOTP(ctx context.Context, email, code string) error {
	if e.dialer == nil {
		return fmt.Errorf("email service not properly configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", e.fromName, e.fromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "VestRoll - Verification Code")

	body := fmt.Sprintf(`
	<html>
	<body>
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<div style="background-color: #f8f9fa; padding: 20px; text-align: center;">
				<h1 style="color: #333; margin: 0;">VestRoll</h1>
			</div>
			<div style="padding: 30px 20px;">
				<h2 style="color: #333; text-align: center;">Verification Code</h2>
				<p style="color: #666; font-size: 16px;">Your verification code is:</p>
				<div style="background-color: #f8f9fa; border: 2px dashed #dee2e6; padding: 20px; text-align: center; margin: 20px 0;">
					<span style="font-size: 32px; font-weight: bold; color: #007bff; letter-spacing: 5px;">%s</span>
				</div>
				<p style="color: #666; font-size: 14px; text-align: center;">
					This code will expire in 5 minutes. If you didn't request this code, please ignore this email.
				</p>
			</div>
			<div style="background-color: #f8f9fa; padding: 15px; text-align: center; font-size: 12px; color: #666;">
				© 2025 VestRoll. All rights reserved.
			</div>
		</div>
	</body>
	</html>
	`, code)

	m.SetBody("text/html", body)

	// Also set plain text version
	plainText := fmt.Sprintf("Your VestRoll verification code is: %s\n\nThis code expires in 5 minutes.\n\nIf you didn't request this code, please ignore this email.", code)
	m.AddAlternative("text/plain", plainText)

	return e.dialer.DialAndSend(m)
}

func (e *EmailService) IsConfigured() bool {
	return e.dialer != nil
}