# Postman Testing Guide - Matchee APIs

## 📁 Files để Import

1. **Matchee_Auth_API.postman_collection.json** - Collection cho Authentication APIs
2. **Matchee_Environment.postman_environment.json** - Environment cho Auth APIs
3. **Matchee_Player_Profile_API.postman_collection.json** - Collection cho Player Profile APIs
4. **Matchee_Player_Environment.postman_environment.json** - Environment cho Player Profile APIs

## 🚀 Cách Import và Sử dụng

### Bước 1: Import Collections
1. Mở Postman
2. Click **Import** button
3. Chọn các file `.json` collection files
4. Click **Import**

### Bước 2: Import Environments
1. Click **Environments** tab
2. Click **Import**
3. Chọn các file environment `.json` files
4. Click **Import**

### Bước 3: Chọn Environment
1. Click dropdown environment ở góc trên bên phải
2. Chọn **Matchee Environment** hoặc **Matchee Player Environment**

## 🧪 Test Flow

### Authentication Flow
1. **Login User** - Lấy access token
2. **Register User** - Tạo user mới (optional)
3. **Get Current User** - Kiểm tra user info
4. **Update Profile** - Cập nhật profile
5. **Change Password** - Đổi mật khẩu

### Player Profile Flow
1. **Login User** - Lấy access token
2. **Create/Update Player Profile** - Tạo/cập nhật hồ sơ người chơi
3. **Get Current User Player Profile** - Lấy hồ sơ hiện tại
4. **Get Player Suggestions** - Lấy gợi ý đối thủ

## 📋 Test Scenarios

### Scenario 1: Complete Auth Flow
```
1. Health Check
2. Register User
3. Login User
4. Get Current User
5. Update Profile
6. Change Password
7. Logout
```

### Scenario 2: Complete Player Profile Flow
```
1. Health Check
2. Login User
3. Create Player Profile
4. Get Player Profile
5. Get Suggestions by Level
6. Get Suggestions by Location
7. Get Combined Suggestions
8. Update Player Profile
```

## 🔧 Environment Variables

### Auth Environment
- `base_url`: http://localhost:8080
- `access_token`: (auto-filled from login)
- `refresh_token`: (auto-filled from login)
- `user_phone`: 0123456789
- `user_password`: password123

### Player Environment
- `base_url`: http://localhost:8080
- `access_token`: (auto-filled from login)
- `user_id`: (auto-filled from profile)
- `player_profile_id`: (auto-filled from profile)
- `test_latitude`: 10.8231
- `test_longitude`: 106.6297
- `test_radius`: 5

## 🎯 Test Cases

### 1. Authentication Tests
- ✅ Login với valid credentials
- ✅ Login với invalid credentials
- ✅ Register user mới
- ✅ Get current user info
- ✅ Update profile
- ✅ Change password
- ✅ Refresh token
- ✅ Logout

### 2. Player Profile Tests
- ✅ Create player profile
- ✅ Update player profile
- ✅ Get current user profile
- ✅ Get profile by ID
- ✅ Get suggestions by level
- ✅ Get suggestions by location
- ✅ Get combined suggestions
- ✅ Validation errors

## 📊 Expected Responses

### Successful Login Response
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
  "user": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com"
  }
}
```

### Successful Player Profile Response
```json
{
  "id": 1,
  "userId": 1,
  "level": "good",
  "gender": "male",
  "preferredLocation": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "bio": "Passionate tennis player",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Successful Suggestions Response
```json
{
  "players": [
    {
      "id": 2,
      "userId": 2,
      "level": "good",
      "gender": "female",
      "preferredLocation": "District 1",
      "latitude": 10.7769,
      "longitude": 106.7009,
      "distance": 2.5,
      "user": {
        "id": 2,
        "fullName": "Trần Thị B",
        "phone": "0987654321"
      }
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 5
}
```

## 🚨 Error Testing

### Test Error Cases
1. **401 Unauthorized** - Không có token hoặc token invalid
2. **400 Bad Request** - Validation errors
3. **404 Not Found** - Profile không tồn tại
4. **409 Conflict** - Duplicate phone/email

### Common Error Responses
```json
{
  "error": "User not authenticated"
}
```

```json
{
  "error": "validation error message"
}
```

## 💡 Tips

1. **Auto-save Tokens**: Collection tự động lưu tokens từ login response
2. **Environment Variables**: Sử dụng variables để dễ dàng thay đổi values
3. **Test Scripts**: Mỗi request có test scripts để validate responses
4. **Complete Flow**: Sử dụng "Complete Player Profile Flow" để test toàn bộ flow
5. **Debug**: Check console logs để debug issues

## 🔍 Debugging

### Check Environment Variables
1. Click **Environments** tab
2. Select environment
3. Verify all variables are set correctly

### Check Response Logs
1. Open **Console** (View → Show Postman Console)
2. Check logs for debugging info
3. Verify token extraction

### Common Issues
- **Token not saved**: Check if login response contains `accessToken`
- **401 errors**: Verify token is valid and not expired
- **Empty suggestions**: Check if there are other players in database
- **Location errors**: Verify latitude/longitude values are valid

## 📝 Notes

- Tất cả Player Profile APIs đều cần authentication
- Suggestions API hỗ trợ nhiều query parameters
- Distance calculation sử dụng Haversine formula
- Level options: beginner, average, good, pro
- Gender options: male, female, other
