# API Response Examples

## Health Check

### Request
```bash
GET /health
```

### Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Service is healthy",
  "data": {
    "app": "matchee-services",
    "version": "1.0.0",
    "status": "running"
  }
}
```

## User Registration

### Request
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "fullName": "Nguyễn Văn A",
  "phone": "0123456789",
  "email": "user@example.com",
  "password": "password123"
}
```

### Success Response (201)
```json
{
  "status": "success",
  "code": 201,
  "message": "User registered successfully",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
    "user": {
      "id": 1,
      "fullName": "Nguyễn Văn A",
      "phone": "0123456789",
      "email": "user@example.com",
      "userRoles": [
        {
          "userId": 1,
          "roleId": 1,
          "role": {
            "id": 1,
            "name": "player"
          }
        }
      ]
    }
  }
}
```

### Error Response (409)
```json
{
  "status": "fail",
  "code": 409,
  "message": "phone number already exists",
  "data": null
}
```

## User Login

### Request
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "phone": "0123456789",
  "password": "password123"
}
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "Login successful",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
    "user": {
      "id": 1,
      "fullName": "Nguyễn Văn A",
      "phone": "0123456789",
      "email": "user@example.com"
    }
  }
}
```

### Error Response (401)
```json
{
  "status": "fail",
  "code": 401,
  "message": "invalid credentials",
  "data": null
}
```

## Get Current User

### Request
```bash
GET /api/v1/users/me
Authorization: Bearer <access_token>
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "User information retrieved successfully",
  "data": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com",
    "userRoles": [
      {
        "userId": 1,
        "roleId": 1,
        "role": {
          "id": 1,
          "name": "player"
        }
      }
    ]
  }
}
```

### Error Response (401)
```json
{
  "status": "fail",
  "code": 401,
  "message": "User not authenticated",
  "data": null
}
```

## Create Player Profile

### Request
```bash
POST /api/v1/player/profile
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "level": "good",
  "gender": "male",
  "preferredLocation": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "bio": "Passionate tennis player"
}
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "Player profile created/updated successfully",
  "data": {
    "id": 1,
    "userId": 1,
    "level": "good",
    "gender": "male",
    "preferredLocation": "Ho Chi Minh City",
    "latitude": 10.8231,
    "longitude": 106.6297,
    "bio": "Passionate tennis player",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "user": {
      "id": 1,
      "fullName": "Nguyễn Văn A",
      "phone": "0123456789",
      "email": "user@example.com"
    }
  }
}
```

### Error Response (400)
```json
{
  "status": "fail",
  "code": 400,
  "message": "Key: 'CreateUpdatePlayerProfileRequest.Level' Error:Field validation for 'Level' failed on the 'oneof' tag",
  "data": null
}
```

## Get Player Suggestions

### Request
```bash
GET /api/v1/player/suggestions?level=good&latitude=10.8231&longitude=106.6297&radius=10&limit=5
Authorization: Bearer <access_token>
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "Player suggestions retrieved successfully",
  "data": {
    "players": [
      {
        "id": 2,
        "userId": 2,
        "level": "good",
        "gender": "female",
        "preferredLocation": "District 1, Ho Chi Minh City",
        "latitude": 10.7769,
        "longitude": 106.7009,
        "distance": 2.5,
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-01T00:00:00Z",
        "user": {
          "id": 2,
          "fullName": "Trần Thị B",
          "phone": "0987654321",
          "email": "user2@example.com"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 5
  }
}
```

## Update Profile

### Request
```bash
PUT /api/v1/users/me
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "fullName": "Nguyễn Văn B",
  "email": "newemail@example.com"
}
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "fullName": "Nguyễn Văn B",
    "phone": "0123456789",
    "email": "newemail@example.com",
    "userRoles": [
      {
        "userId": 1,
        "roleId": 1,
        "role": {
          "id": 1,
          "name": "player"
        }
      }
    ]
  }
}
```

### Error Response (409)
```json
{
  "status": "fail",
  "code": 409,
  "message": "email already exists",
  "data": null
}
```

## Change Password

### Request
```bash
POST /api/v1/users/change-password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "currentPassword": "password123",
  "newPassword": "newpassword456"
}
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "Password changed successfully",
  "data": null
}
```

### Error Response (401)
```json
{
  "status": "fail",
  "code": 401,
  "message": "current password is incorrect",
  "data": null
}
```

## Logout

### Request
```bash
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refreshToken": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
}
```

### Success Response (200)
```json
{
  "status": "success",
  "code": 200,
  "message": "Logged out successfully",
  "data": null
}
```

## Common Error Responses

### Validation Error (400)
```json
{
  "status": "fail",
  "code": 400,
  "message": "Key: 'RegisterRequest.Phone' Error:Field validation for 'Phone' failed on the 'len' tag",
  "data": null
}
```

### Not Found (404)
```json
{
  "status": "fail",
  "code": 404,
  "message": "Player profile not found",
  "data": null
}
```

### Internal Server Error (500)
```json
{
  "status": "fail",
  "code": 500,
  "message": "Internal server error",
  "data": null
}
```
