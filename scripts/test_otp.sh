#!/bin/bash

echo "🧪 Testing VestRoll OTP Endpoints"
echo "================================="

BASE_URL="http://localhost:8080/api/v1/auth"

echo ""
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/../../../health" | jq '.' || echo "Response: $(curl -s http://localhost:8080/health)"

echo ""
echo "2. Testing send SMS OTP (will fail without proper Twilio config)..."
curl -X POST "$BASE_URL/send-otp" \
  -H "Content-Type: application/json" \
  -d '{"identifier": "+1234567890", "type": "sms"}' \
  -w "\nStatus: %{http_code}\n" || echo "Failed to send request"

echo ""
echo "3. Testing send Email OTP (will fail without proper SMTP config)..."
curl -X POST "$BASE_URL/send-otp" \
  -H "Content-Type: application/json" \
  -d '{"identifier": "test@example.com", "type": "email"}' \
  -w "\nStatus: %{http_code}\n" || echo "Failed to send request"

echo ""
echo "4. Testing invalid phone format..."
curl -X POST "$BASE_URL/send-otp" \
  -H "Content-Type: application/json" \
  -d '{"identifier": "1234567890", "type": "sms"}' \
  -w "\nStatus: %{http_code}\n" || echo "Failed to send request"

echo ""
echo "5. Testing invalid email format..."
curl -X POST "$BASE_URL/send-otp" \
  -H "Content-Type: application/json" \
  -d '{"identifier": "invalid-email", "type": "email"}' \
  -w "\nStatus: %{http_code}\n" || echo "Failed to send request"

echo ""
echo "6. Testing verify OTP with non-existent code..."
curl -X POST "$BASE_URL/verify-otp" \
  -H "Content-Type: application/json" \
  -d '{"identifier": "+1234567890", "code": "123456", "type": "sms"}' \
  -w "\nStatus: %{http_code}\n" || echo "Failed to send request"

echo ""
echo "7. Testing invalid OTP format..."
curl -X POST "$BASE_URL/verify-otp" \
  -H "Content-Type: application/json" \
  -d '{"identifier": "+1234567890", "code": "12345", "type": "sms"}' \
  -w "\nStatus: %{http_code}\n" || echo "Failed to send request"

echo ""
echo "✅ Test completed!"
echo ""
echo "📝 Notes:"
echo "- SMS/Email OTP sending will fail without proper service configuration"
echo "- Validation errors should return status 400"
echo "- Non-existent OTP verification should return status 400"
echo "- The server gracefully handles missing Redis/service configurations"