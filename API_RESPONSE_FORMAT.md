# Matchee API Response Format

## Standard Response Structure

Tất cả API responses đều tuân theo chuẩn format sau:

```json
{
  "status": "success" | "fail",
  "code": 200 | 400 | 401 | 404 | 409 | 500,
  "message": "string",
  "data": any | null
}
```

## Response Fields

### status
- **"success"**: Request thành công
- **"fail"**: Request thất bại

### code
- **200**: OK - Request thành công
- **201**: Created - Tạo mới thành công
- **400**: Bad Request - Lỗi validation hoặc request không hợp lệ
- **401**: Unauthorized - Không có quyền truy cập
- **404**: Not Found - Không tìm thấy resource
- **409**: Conflict - Xung đột dữ liệu (duplicate)
- **500**: Internal Server Error - Lỗi server

### message
- Mô tả chi tiết về kết quả của request
- Luôn có giá trị, không bao giờ null

### data
- Dữ liệu trả về (có thể là object, array, hoặc null)
- Chỉ có khi status = "success"

## Response Examples

### Success Responses

#### 1. Health Check
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

#### 2. User Registration
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
      "email": "user@example.com"
    }
  }
}
```

#### 3. Login Success
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

#### 4. Get Current User
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

#### 5. Player Profile Created
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
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 6. Player Suggestions
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
}
```

### Error Responses

#### 1. Validation Error (400)
```json
{
  "status": "fail",
  "code": 400,
  "message": "Key: 'RegisterRequest.FullName' Error:Field validation for 'FullName' failed on the 'required' tag",
  "data": null
}
```

#### 2. Unauthorized (401)
```json
{
  "status": "fail",
  "code": 401,
  "message": "User not authenticated",
  "data": null
}
```

#### 3. Not Found (404)
```json
{
  "status": "fail",
  "code": 404,
  "message": "Player profile not found",
  "data": null
}
```

#### 4. Conflict (409)
```json
{
  "status": "fail",
  "code": 409,
  "message": "phone number already exists",
  "data": null
}
```

#### 5. Internal Server Error (500)
```json
{
  "status": "fail",
  "code": 500,
  "message": "Internal server error",
  "data": null
}
```

## HTTP Status Codes Mapping

| HTTP Status | Response Code | Description |
|-------------|---------------|-------------|
| 200 | 200 | OK - Success |
| 201 | 201 | Created - Resource created |
| 400 | 400 | Bad Request - Validation error |
| 401 | 401 | Unauthorized - Authentication required |
| 404 | 404 | Not Found - Resource not found |
| 409 | 409 | Conflict - Duplicate resource |
| 500 | 500 | Internal Server Error - Server error |

## Implementation Details

### Response Helper Functions

```go
// Success responses
entity.OKResponse(message, data)
entity.CreatedResponse(message, data)

// Error responses
entity.BadRequestResponse(message)
entity.UnauthorizedResponse(message)
entity.NotFoundResponse(message)
entity.ConflictResponse(message)
entity.InternalServerErrorResponse(message)
```

### Usage in Controllers

```go
// Success response
c.JSON(http.StatusOK, entity.OKResponse("Operation successful", data))

// Error response
c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid input"))
```

## Benefits

1. **Consistency**: Tất cả APIs đều có format response giống nhau
2. **Clarity**: Dễ dàng hiểu được kết quả của request
3. **Error Handling**: Xử lý lỗi một cách rõ ràng và nhất quán
4. **Frontend Integration**: Dễ dàng tích hợp với frontend
5. **Debugging**: Dễ dàng debug và troubleshoot

## Migration Notes

- Tất cả existing APIs đã được cập nhật để sử dụng chuẩn response format
- Health check endpoint cũng đã được cập nhật
- Postman collections sẽ cần được cập nhật để phản ánh format mới
- Frontend code cần được cập nhật để handle response format mới
