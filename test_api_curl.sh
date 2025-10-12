#!/bin/bash

# Matchee Auth API Test Script
# Base URL
BASE_URL="http://localhost:8080"

echo "🚀 Testing Matchee Auth API..."
echo "================================"

# Variables to store tokens
ACCESS_TOKEN=""
REFRESH_TOKEN=""

# Function to make API calls
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    
    echo "📡 $method $endpoint"
    if [ -n "$data" ]; then
        echo "📤 Data: $data"
    fi
    if [ -n "$headers" ]; then
        echo "📋 Headers: $headers"
    fi
    
    if [ -n "$data" ] && [ -n "$headers" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "$headers" \
            -d "$data")
    elif [ -n "$data" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    elif [ -n "$headers" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "$headers")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint")
    fi
    
    echo "📥 Response: $response"
    echo "---"
    
    # Extract tokens if present
    if echo "$response" | grep -q "accessToken"; then
        ACCESS_TOKEN=$(echo "$response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
        REFRESH_TOKEN=$(echo "$response" | grep -o '"refreshToken":"[^"]*"' | cut -d'"' -f4)
        echo "🔑 Tokens extracted and saved"
    fi
    
    # Check for role information
    if echo "$response" | grep -q "userRoles"; then
        echo "👤 User roles information included in response"
    fi
}

echo ""
echo "1️⃣ Health Check"
make_request "GET" "/health"

echo ""
echo "2️⃣ Register User (with default player role)"
make_request "POST" "/api/v1/auth/register" '{
  "fullName": "Nguyễn Văn A",
  "phone": "0123456789",
  "email": "user@example.com",
  "password": "password123"
}'

echo ""
echo "3️⃣ Login User"
make_request "POST" "/api/v1/auth/login" '{
  "phone": "0123456789",
  "password": "password123"
}'

if [ -n "$ACCESS_TOKEN" ]; then
    echo ""
    echo "4️⃣ Get Current User (with token)"
    make_request "GET" "/api/v1/users/me" "" "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "5️⃣ Update Profile (with token)"
    make_request "PUT" "/api/v1/users/me" '{
      "fullName": "Nguyễn Văn B",
      "email": "newemail@example.com",
      "phone": "0987654321"
    }' "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "6️⃣ Change Password (with token)"
    make_request "POST" "/api/v1/users/change-password" '{
      "currentPassword": "password123",
      "newPassword": "newpassword456"
    }' "Authorization: Bearer $ACCESS_TOKEN"
    
    if [ -n "$REFRESH_TOKEN" ]; then
        echo ""
        echo "7️⃣ Refresh Token"
        make_request "POST" "/api/v1/auth/refresh" "{
          \"refreshToken\": \"$REFRESH_TOKEN\"
        }"
        
        echo ""
        echo "8️⃣ Logout"
        make_request "POST" "/api/v1/auth/logout" "{
          \"refreshToken\": \"$REFRESH_TOKEN\"
        }"
    fi
else
    echo "❌ No access token available, skipping authenticated requests"
fi

echo ""
echo "✅ API testing completed!"
echo ""
echo "💡 Tips:"
echo "- Make sure the server is running on $BASE_URL"
echo "- Check the database connection"
echo "- Verify JWT_SECRET is set in environment"
echo "- Run migrations: make migrate-up"
