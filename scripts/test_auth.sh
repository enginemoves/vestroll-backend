#!/bin/bash

# Test script for VestRoll Auth API
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1/auth"

echo "=== VestRoll Auth API Test ==="
echo ""

# Test 1: Register a new user
echo "1. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@vestroll.com",
    "password": "TestPassword123!",
    "full_name": "Test User"
  }')

echo "Register Response: $REGISTER_RESPONSE"
echo ""

# Extract token from response (basic extraction)
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Extracted Token: ${TOKEN:0:50}..."
echo ""

# Test 2: Try to register the same user again (should fail)
echo "2. Testing duplicate email registration..."
DUPLICATE_RESPONSE=$(curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@vestroll.com",
    "password": "AnotherPassword123!",
    "full_name": "Another User"
  }')

echo "Duplicate Registration Response: $DUPLICATE_RESPONSE"
echo ""

# Test 3: Login with correct credentials
echo "3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@vestroll.com",
    "password": "TestPassword123!"
  }')

echo "Login Response: $LOGIN_RESPONSE"
echo ""

# Test 4: Login with wrong password
echo "4. Testing login with wrong password..."
WRONG_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@vestroll.com",
    "password": "WrongPassword123!"
  }')

echo "Wrong Login Response: $WRONG_LOGIN_RESPONSE"
echo ""

# Test 5: Register with invalid email
echo "5. Testing registration with invalid email..."
INVALID_EMAIL_RESPONSE=$(curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "TestPassword123!",
    "full_name": "Test User"
  }')

echo "Invalid Email Response: $INVALID_EMAIL_RESPONSE"
echo ""

# Test 6: Register with weak password
echo "6. Testing registration with weak password..."
WEAK_PASSWORD_RESPONSE=$(curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test2@vestroll.com",
    "password": "weak",
    "full_name": "Test User"
  }')

echo "Weak Password Response: $WEAK_PASSWORD_RESPONSE"
echo ""

echo "=== Test Complete ==="
echo ""
echo "Expected Results:"
echo "1. Registration should succeed with 201 status"
echo "2. Duplicate registration should fail with 409 status"
echo "3. Login should succeed with 200 status"
echo "4. Wrong password login should fail with 401 status"
echo "5. Invalid email should fail with 400 status"
echo "6. Weak password should fail with 400 status"
