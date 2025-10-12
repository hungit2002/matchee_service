# Matchee Services API Documentation

## Authentication & User Management APIs

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### 1. Register User
**POST** `/auth/register`

Tạo tài khoản mới với role mặc định là "player". Role sẽ được tự động gán cho user mới.

**Request Body:**
```json
{
  "fullName": "Nguyễn Văn A",
  "phone": "0123456789",
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
  "user": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 2. Login User
**POST** `/auth/login`

Đăng nhập và nhận JWT token.

**Request Body:**
```json
{
  "phone": "0123456789",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
  "user": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 3. Refresh Token
**POST** `/auth/refresh`

Làm mới access token bằng refresh token.

**Request Body:**
```json
{
  "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
}
```

**Response (200):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1",
  "user": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 4. Logout
**POST** `/auth/logout`

Đăng xuất và xóa refresh token.

**Request Body:**
```json
{
  "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
}
```

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

### User Management Endpoints

#### 5. Get Current User
**GET** `/users/me`

Lấy thông tin user hiện tại (cần authentication).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "id": 1,
  "fullName": "Nguyễn Văn A",
  "phone": "0123456789",
  "email": "user@example.com",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "userRoles": [
    {
      "userId": 1,
      "roleId": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z",
      "role": {
        "id": 1,
        "name": "player",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    }
  ]
}
```

#### 6. Update Profile
**PUT** `/users/me`

Cập nhật thông tin cá nhân (cần authentication).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "fullName": "Nguyễn Văn B",
  "email": "newemail@example.com",
  "phone": "0987654321"
}
```

**Response (200):**
```json
{
  "id": 1,
  "fullName": "Nguyễn Văn B",
  "phone": "0987654321",
  "email": "newemail@example.com",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z",
  "userRoles": [
    {
      "userId": 1,
      "roleId": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z",
      "role": {
        "id": 1,
        "name": "player",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z"
      }
    }
  ]
}
```

#### 7. Change Password
**POST** `/users/change-password`

Đổi mật khẩu (cần authentication).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "currentPassword": "oldpassword123",
  "newPassword": "newpassword456"
}
```

**Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

## Error Responses

### Common Error Format
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid credentials or token)
- `404` - Not Found
- `409` - Conflict (duplicate phone/email)
- `500` - Internal Server Error

## Environment Variables

Thêm các biến môi trường sau vào file `.env`:

```env
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=15m
REFRESH_EXPIRY=168h
```

## Security Features

1. **Password Hashing**: Sử dụng bcrypt để hash password
2. **JWT Tokens**: Access token có thời hạn ngắn (15 phút)
3. **Refresh Tokens**: Có thời hạn dài hơn (7 ngày) và được lưu trong database
4. **Token Revocation**: Refresh token bị xóa khi logout hoặc đổi password
5. **Input Validation**: Validation cho tất cả input fields
6. **SQL Injection Protection**: Sử dụng GORM ORM

## Database Schema

### auth_tokens Table
```sql
CREATE TABLE auth_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY unique_token (token(255)),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```
