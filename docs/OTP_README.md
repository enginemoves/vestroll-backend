# OTP Verification System

## Overview

This implementation provides a complete OTP (One-Time Password) verification system for multi-factor authentication with the following features:

- **Dual delivery methods**: SMS and Email
- **6-digit numeric codes** with 5-minute expiration
- **Redis storage** with automatic TTL cleanup
- **Rate limiting** to prevent abuse
- **Comprehensive error handling** with appropriate HTTP status codes
- **Input validation** for phone numbers and email addresses

## API Endpoints

### POST `/api/v1/auth/send-otp`

Generates and sends an OTP code to the specified identifier.

**Request Body:**
```json
{
  "identifier": "+1234567890",  // Phone (international format) or email
  "type": "sms"                 // "sms" or "email"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "OTP sent successfully"
}
```

**Error Responses:**
- `400` - Validation error (invalid phone/email format)
- `429` - Rate limit exceeded
- `503` - Service unavailable (SMS/Email service not configured)

### POST `/api/v1/auth/verify-otp`

Verifies an OTP code against the stored value.

**Request Body:**
```json
{
  "identifier": "+1234567890",
  "code": "123456",
  "type": "sms"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "OTP verified successfully"
}
```

**Error Responses:**
- `400` - Invalid OTP code, expired OTP, or validation error
- `429` - Maximum verification attempts exceeded

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure the following:

#### Required for SMS (Twilio):
```bash
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_FROM_PHONE=+1234567890
```

#### Required for Email (SMTP):
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM_EMAIL=noreply@vestroll.com
SMTP_FROM_NAME=VestRoll
```

#### Required for Redis:
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

#### Optional Configuration:
```bash
OTP_RATE_LIMIT_MAX=5
OTP_RATE_LIMIT_WINDOW_MINUTES=15
```

## Security Features

### Rate Limiting
- **Global rate limiting**: 10 requests per minute per IP for all OTP endpoints
- **Per-identifier rate limiting**: 5 OTP requests per 15 minutes per phone/email
- **Verification attempts**: Maximum 3 attempts per OTP before deletion

### Code Security
- **Cryptographically secure** random number generation
- **One-time use**: OTP codes are deleted after successful verification
- **Time-based expiration**: 5-minute TTL with automatic cleanup
- **Attempt tracking**: Failed attempts are counted and limited

### Input Validation
- **Phone numbers**: Must be in international format (`+1234567890`)
- **Email addresses**: Standard email format validation
- **OTP codes**: Must be exactly 6 digits

## Architecture

### Components

1. **Models** (`internal/models/otp.go`)
   - Request/response structures
   - OTP data models

2. **Repository** (`internal/repository/otp_repository.go`)
   - Redis storage operations
   - Rate limiting logic

3. **Services**
   - `internal/services/otp_service.go` - Core OTP logic
   - `internal/services/sms_service.go` - Twilio SMS integration
   - `internal/services/email_service.go` - SMTP email integration

4. **Handlers** (`internal/handlers/otp_handler.go`)
   - HTTP request handling
   - Error response mapping

5. **Middleware** (`internal/middleware/rate_limit.go`)
   - Token bucket rate limiting
   - IP-based request throttling

### Data Flow

1. **Send OTP**:
   ```
   Request → Validation → Rate Check → Generate Code → Store Redis → Send SMS/Email
   ```

2. **Verify OTP**:
   ```
   Request → Validation → Retrieve Redis → Check Expiry → Verify Code → Delete OTP
   ```

## Development Setup

### Prerequisites
- Go 1.25+
- Redis server
- Twilio account (for SMS)
- SMTP server access (for email)

### Quick Start

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Start Redis** (using Docker):
   ```bash
   docker run -d -p 6379:6379 redis:alpine
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

### Testing the API

#### Send SMS OTP:
```bash
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"identifier": "+1234567890", "type": "sms"}'
```

#### Send Email OTP:
```bash
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"identifier": "user@example.com", "type": "email"}'
```

#### Verify OTP:
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"identifier": "+1234567890", "code": "123456", "type": "sms"}'
```

## Production Considerations

### Security
- Use environment variables for all secrets
- Enable HTTPS in production
- Consider implementing additional fraud detection
- Monitor for unusual patterns

### Scalability
- Redis clustering for high availability
- Consider using Redis Streams for audit logging
- Implement proper logging and monitoring

### Compliance
- Ensure compliance with SMS regulations (opt-out mechanisms)
- Implement proper data retention policies
- Consider GDPR/privacy requirements

## Error Handling

The system provides detailed error responses with appropriate HTTP status codes:

| Error | Status Code | Description |
|-------|-------------|-------------|
| `validation_error` | 400 | Invalid request format or data |
| `invalid_otp` | 400 | Wrong OTP code |
| `otp_expired` | 400 | OTP has expired |
| `rate_limit_exceeded` | 429 | Too many requests |
| `max_attempts_exceeded` | 429 | Too many verification attempts |
| `service_unavailable` | 503 | SMS/Email service not configured |

## Monitoring

Key metrics to monitor:
- OTP generation/verification rates
- Success/failure ratios
- Rate limiting triggers
- Service availability (Twilio, SMTP, Redis)
- Response times