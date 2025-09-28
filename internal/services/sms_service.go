package services

import (
	"context"
	"fmt"

	"github.com/codeZe-us/vestroll-backend/internal/config"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService struct {
	client    *twilio.RestClient
	fromPhone string
}

func NewSMSService(cfg config.TwilioConfig) *SMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})

	return &SMSService{
		client:    client,
		fromPhone: cfg.FromPhone,
	}
}

func (s *SMSService) SendOTP(ctx context.Context, phoneNumber, code string) error {
	if s.client == nil || s.fromPhone == "" {
		return fmt.Errorf("SMS service not properly configured")
	}

	message := fmt.Sprintf("Your VestRoll verification code is: %s. This code expires in 5 minutes.", code)

	params := &api.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(s.fromPhone)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	return err
}

func (s *SMSService) IsConfigured() bool {
	return s.client != nil && s.fromPhone != ""
}