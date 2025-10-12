# Matchee Auth API - CURL Commands

## Base URL
```bash
BASE_URL="http://localhost:8080"
```

## 1. Health Check
```bash
curl -X GET "$BASE_URL/health"
```

## 2. Register User
```bash
curl -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com",
    "password": "password123"
  }'
```

## 3. Login User
```bash
curl -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "0123456789",
    "password": "password123"
  }'
```

**Response sẽ chứa accessToken và refreshToken - lưu lại để sử dụng cho các request tiếp theo.**

## 4. Get Current User (Cần Authentication)
```bash
# Thay YOUR_ACCESS_TOKEN bằng token thực tế từ login response
curl -X GET "$BASE_URL/api/v1/users/me" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 5. Update Profile (Cần Authentication)
```bash
curl -X PUT "$BASE_URL/api/v1/users/me" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "fullName": "Nguyễn Văn B",
    "email": "newemail@example.com",
    "phone": "0987654321"
  }'
```

## 6. Change Password (Cần Authentication)
```bash
curl -X POST "$BASE_URL/api/v1/users/change-password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "currentPassword": "password123",
    "newPassword": "newpassword456"
  }'
```

## 7. Refresh Token
```bash
# Thay YOUR_REFRESH_TOKEN bằng refresh token thực tế
curl -X POST "$BASE_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "YOUR_REFRESH_TOKEN"
  }'
```

## 8. Logout
```bash
curl -X POST "$BASE_URL/api/v1/auth/logout" \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "YOUR_REFRESH_TOKEN"
  }'
```

## Test Sequence (Thứ tự test)

### Bước 1: Kiểm tra server
```bash
curl -X GET "http://localhost:8080/health"
```

### Bước 2: Đăng ký user mới
```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Test User",
    "phone": "0123456789",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Bước 3: Đăng nhập
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "0123456789",
    "password": "password123"
  }'
```

**Lưu lại accessToken và refreshToken từ response**

### Bước 4: Test các API cần authentication
```bash
# Lấy thông tin user
curl -X GET "http://localhost:8080/api/v1/users/me" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Cập nhật profile
curl -X PUT "http://localhost:8080/api/v1/users/me" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "fullName": "Updated Name",
    "email": "updated@example.com"
  }'

# Đổi mật khẩu
curl -X POST "http://localhost:8080/api/v1/users/change-password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "currentPassword": "password123",
    "newPassword": "newpassword456"
  }'
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "invalid credentials"
}
```

### 409 Conflict
```json
{
  "error": "phone number already exists"
}
```

## Environment Setup

1. **Start Database**: Đảm bảo MySQL đang chạy
2. **Run Migrations**: `make migrate-up`
3. **Set Environment Variables**:
   ```bash
   export JWT_SECRET="your-secret-key-change-in-production"
   export JWT_EXPIRY="15m"
   export REFRESH_EXPIRY="168h"
   ```
4. **Start Server**: `make run`

## Tips

- **Token Expiry**: Access token có thời hạn 15 phút, refresh token có thời hạn 7 ngày
- **Password Requirements**: Tối thiểu 6 ký tự
- **Phone Format**: Phải có đúng 10 số
- **Email**: Optional, nhưng nếu có thì phải đúng format email
- **Authentication**: Tất cả API `/users/*` đều cần Bearer token
