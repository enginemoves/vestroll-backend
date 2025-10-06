package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OTP      OTPConfig
	Twilio   TwilioConfig
	SMTP     SMTPConfig
	Google   GoogleConfig
	Apple    AppleConfig
  EmailVerification  EmailVerificationConfig
}
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type AppleConfig struct {
	ClientID     string
	ClientSecret string
	TeamID       string
	KeyID        string
	PrivateKey   string
	RedirectURL  string
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type OTPConfig struct {
	Length    int
	TTL       time.Duration
	RateLimit RateLimitConfig
}

type RateLimitConfig struct {
	MaxRequests int
	WindowSize  time.Duration
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromPhone  string
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

type EmailVerificationConfig struct {
	TTL         time.Duration
	LinkBaseURL string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "vestroll"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
			TTL:    time.Duration(getEnvAsInt("JWT_TTL_HOURS", 24)) * time.Hour,
		},
		OTP: OTPConfig{
			Length: 6,
			TTL:    5 * time.Minute,
			RateLimit: RateLimitConfig{
				MaxRequests: getEnvAsInt("OTP_RATE_LIMIT_MAX", 5),
				WindowSize:  time.Duration(getEnvAsInt("OTP_RATE_LIMIT_WINDOW_MINUTES", 15)) * time.Minute,
			},
		},
		Twilio: TwilioConfig{
			AccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
			AuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
			FromPhone:  getEnv("TWILIO_FROM_PHONE", ""),
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:      getEnvAsInt("SMTP_PORT", 587),
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			FromEmail: getEnv("SMTP_FROM_EMAIL", "noreply@vestroll.com"),
			FromName:  getEnv("SMTP_FROM_NAME", "VestRoll"),
		},
		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		Apple: AppleConfig{
			ClientID:     getEnv("APPLE_CLIENT_ID", ""),
			ClientSecret: getEnv("APPLE_CLIENT_SECRET", ""),
			TeamID:       getEnv("APPLE_TEAM_ID", ""),
			KeyID:        getEnv("APPLE_KEY_ID", ""),
			PrivateKey:   getEnv("APPLE_PRIVATE_KEY", ""),
			RedirectURL:  getEnv("APPLE_REDIRECT_URL", ""),
		EmailVerification: EmailVerificationConfig{
			TTL:         time.Duration(getEnvAsInt("EMAIL_VERIFICATION_TTL_MINUTES", 60*24)) * time.Minute,
			LinkBaseURL: getEnv("EMAIL_VERIFICATION_LINK_BASE_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
