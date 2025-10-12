#!/bin/bash

# Matchee Player Profile API Test Script
# Base URL
BASE_URL="http://localhost:8080"

echo "🏓 Testing Matchee Player Profile API..."
echo "========================================"

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
}

echo ""
echo "1️⃣ Health Check"
make_request "GET" "/health"

echo ""
echo "2️⃣ Login User (to get token)"
make_request "POST" "/api/v1/auth/login" '{
  "phone": "0123456789",
  "password": "password123"
}'

if [ -n "$ACCESS_TOKEN" ]; then
    echo ""
    echo "3️⃣ Create/Update Player Profile"
    make_request "POST" "/api/v1/player/profile" '{
      "level": "good",
      "gender": "male",
      "preferredLocation": "Ho Chi Minh City",
      "latitude": 10.8231,
      "longitude": 106.6297,
      "bio": "Passionate tennis player looking for matches"
    }' "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "4️⃣ Get Current User Player Profile"
    make_request "GET" "/api/v1/player/profile" "" "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "5️⃣ Get Player Suggestions (by level)"
    make_request "GET" "/api/v1/player/suggestions?level=good&limit=5" "" "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "6️⃣ Get Player Suggestions (by location)"
    make_request "GET" "/api/v1/player/suggestions?latitude=10.8231&longitude=106.6297&radius=5&limit=5" "" "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "7️⃣ Get Player Suggestions (combined criteria)"
    make_request "GET" "/api/v1/player/suggestions?level=average&latitude=10.8231&longitude=106.6297&radius=10&limit=10" "" "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "8️⃣ Update Player Profile"
    make_request "POST" "/api/v1/player/profile" '{
      "level": "pro",
      "gender": "male",
      "preferredLocation": "District 1, Ho Chi Minh City",
      "latitude": 10.7769,
      "longitude": 106.7009,
      "bio": "Professional tennis player with 10+ years experience"
    }' "Authorization: Bearer $ACCESS_TOKEN"
    
    echo ""
    echo "9️⃣ Get Updated Profile"
    make_request "GET" "/api/v1/player/profile" "" "Authorization: Bearer $ACCESS_TOKEN"
    
else
    echo "❌ No access token available, skipping authenticated requests"
    echo "💡 Please run the auth test first to get tokens"
fi

echo ""
echo "✅ Player Profile API testing completed!"
echo ""
echo "💡 Tips:"
echo "- Make sure the server is running on $BASE_URL"
echo "- Ensure you have a valid user account"
echo "- Check the database connection"
echo "- Verify JWT_SECRET is set in environment"
