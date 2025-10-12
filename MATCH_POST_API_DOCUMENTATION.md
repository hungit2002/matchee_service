# Match Post API Documentation

## Tổng quan
API Match Post cho phép người chơi tạo và quản lý các bài đăng tìm đối thủ, tìm nhóm ghép để chơi tennis.

## Base URL
```
{{base_url}}/api/v1/match-posts
```

## Authentication
Một số endpoint yêu cầu authentication thông qua Bearer token:
```
Authorization: Bearer {{access_token}}
```

## Endpoints

### 1. Tạo bài đăng tìm đối
**POST** `/api/v1/match-posts`

Tạo bài đăng mới để tìm đối thủ chơi tennis.

#### Headers
```
Content-Type: application/json
Authorization: Bearer {{access_token}}
```

#### Request Body
```json
{
  "desiredLevel": "intermediate",
  "location": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "venueId": 1,
  "desiredTime": "2024-01-20T18:00:00Z",
  "pricePerPerson": 100000,
  "maxPlayers": 4,
  "note": "Looking for intermediate level players for evening match"
}
```

#### Response
```json
{
  "status": "success",
  "code": 201,
  "message": "Match post created successfully",
  "data": {
    "id": 1,
    "userId": 1,
    "desiredLevel": "intermediate",
    "location": "Ho Chi Minh City",
    "latitude": 10.8231,
    "longitude": 106.6297,
    "venueId": 1,
    "desiredTime": "2024-01-20T18:00:00Z",
    "pricePerPerson": 100000,
    "maxPlayers": 4,
    "note": "Looking for intermediate level players for evening match",
    "status": "open",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z",
    "user": {
      "id": 1,
      "fullName": "John Doe",
      "email": "john@example.com",
      "phone": "0123456789"
    },
    "venue": {
      "id": 1,
      "name": "Tennis Club HCMC",
      "address": "123 Tennis Street, HCMC"
    }
  }
}
```

### 2. Lấy danh sách bài đăng
**GET** `/api/v1/match-posts`

Lấy danh sách bài đăng với các bộ lọc tùy chọn.

#### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| desiredLevel | string | No | Trình độ mong muốn (beginner, average, good, pro) |
| location | string | No | Vị trí tìm kiếm |
| latitude | number | No | Vĩ độ cho tìm kiếm theo vị trí |
| longitude | number | No | Kinh độ cho tìm kiếm theo vị trí |
| radius | number | No | Bán kính tìm kiếm (km) |
| venueId | integer | No | ID của venue |
| status | string | No | Trạng thái (open, matched, cancelled, done) |
| page | integer | No | Số trang (mặc định: 1) |
| limit | integer | No | Số item mỗi trang (mặc định: 20) |

#### Example Request
```
GET /api/v1/match-posts?desiredLevel=intermediate&location=Ho Chi Minh City&latitude=10.8231&longitude=106.6297&radius=10&status=open&page=1&limit=20
```

#### Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Match posts retrieved successfully",
  "data": {
    "matchPosts": [
      {
        "id": 1,
        "userId": 1,
        "desiredLevel": "intermediate",
        "location": "Ho Chi Minh City",
        "latitude": 10.8231,
        "longitude": 106.6297,
        "venueId": 1,
        "desiredTime": "2024-01-20T18:00:00Z",
        "pricePerPerson": 100000,
        "maxPlayers": 4,
        "note": "Looking for intermediate level players",
        "status": "open",
        "createdAt": "2024-01-15T10:00:00Z",
        "updatedAt": "2024-01-15T10:00:00Z",
        "user": {
          "id": 1,
          "fullName": "John Doe",
          "email": "john@example.com"
        },
        "venue": {
          "id": 1,
          "name": "Tennis Club HCMC"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20,
    "totalPages": 1
  }
}
```

### 3. Lấy chi tiết bài đăng
**GET** `/api/v1/match-posts/{id}`

Lấy thông tin chi tiết của một bài đăng.

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | ID của bài đăng |

#### Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Match post retrieved successfully",
  "data": {
    "id": 1,
    "userId": 1,
    "desiredLevel": "intermediate",
    "location": "Ho Chi Minh City",
    "latitude": 10.8231,
    "longitude": 106.6297,
    "venueId": 1,
    "desiredTime": "2024-01-20T18:00:00Z",
    "pricePerPerson": 100000,
    "maxPlayers": 4,
    "note": "Looking for intermediate level players",
    "status": "open",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z",
    "user": {
      "id": 1,
      "fullName": "John Doe",
      "email": "john@example.com",
      "phone": "0123456789"
    },
    "venue": {
      "id": 1,
      "name": "Tennis Club HCMC",
      "address": "123 Tennis Street, HCMC"
    },
    "matchGroups": []
  }
}
```

### 4. Cập nhật bài đăng
**PUT** `/api/v1/match-posts/{id}`

Cập nhật thông tin bài đăng (chỉ chủ bài đăng mới có thể cập nhật).

#### Headers
```
Content-Type: application/json
Authorization: Bearer {{access_token}}
```

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | ID của bài đăng |

#### Request Body
```json
{
  "desiredLevel": "advanced",
  "maxPlayers": 6,
  "note": "Updated: Looking for advanced level players",
  "status": "open"
}
```

#### Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Match post updated successfully",
  "data": {
    "id": 1,
    "userId": 1,
    "desiredLevel": "advanced",
    "location": "Ho Chi Minh City",
    "latitude": 10.8231,
    "longitude": 106.6297,
    "venueId": 1,
    "desiredTime": "2024-01-20T18:00:00Z",
    "pricePerPerson": 100000,
    "maxPlayers": 6,
    "note": "Updated: Looking for advanced level players",
    "status": "open",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:30:00Z",
    "user": {
      "id": 1,
      "fullName": "John Doe",
      "email": "john@example.com"
    },
    "venue": {
      "id": 1,
      "name": "Tennis Club HCMC"
    }
  }
}
```

### 5. Xóa bài đăng
**DELETE** `/api/v1/match-posts/{id}`

Xóa bài đăng (chỉ chủ bài đăng mới có thể xóa).

#### Headers
```
Authorization: Bearer {{access_token}}
```

#### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | ID của bài đăng |

#### Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Match post deleted successfully",
  "data": null
}
```

### 6. Gợi ý bài đăng phù hợp
**GET** `/api/v1/match-posts/suggest`

Lấy danh sách gợi ý bài đăng phù hợp với người dùng.

#### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| userId | integer | Yes | ID của người dùng |
| desiredLevel | string | No | Trình độ mong muốn |
| latitude | number | No | Vĩ độ của người dùng |
| longitude | number | No | Kinh độ của người dùng |
| radius | number | No | Bán kính tìm kiếm (km) |
| limit | integer | No | Số lượng gợi ý (mặc định: 10, tối đa: 50) |

#### Example Request
```
GET /api/v1/match-posts/suggest?userId=1&desiredLevel=intermediate&latitude=10.8231&longitude=106.6297&radius=5&limit=10
```

#### Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Match post suggestions retrieved successfully",
  "data": {
    "suggestions": [
      {
        "id": 2,
        "userId": 3,
        "desiredLevel": "intermediate",
        "location": "District 1, HCMC",
        "latitude": 10.8231,
        "longitude": 106.6297,
        "venueId": 2,
        "desiredTime": "2024-01-21T19:00:00Z",
        "pricePerPerson": 120000,
        "maxPlayers": 4,
        "note": "Evening tennis session",
        "status": "open",
        "createdAt": "2024-01-15T11:00:00Z",
        "updatedAt": "2024-01-15T11:00:00Z",
        "user": {
          "id": 3,
          "fullName": "Jane Smith",
          "email": "jane@example.com"
        },
        "venue": {
          "id": 2,
          "name": "Premium Tennis Center"
        }
      }
    ],
    "total": 1
  }
}
```

## Error Responses

### 400 Bad Request
```json
{
  "status": "error",
  "code": 400,
  "message": "Invalid request data",
  "data": null
}
```

### 401 Unauthorized
```json
{
  "status": "error",
  "code": 401,
  "message": "User not authenticated",
  "data": null
}
```

### 404 Not Found
```json
{
  "status": "error",
  "code": 404,
  "message": "Match post not found",
  "data": null
}
```

## Data Models

### MatchPost
| Field | Type | Description |
|-------|------|-------------|
| id | integer | ID duy nhất của bài đăng |
| userId | integer | ID của người tạo bài đăng |
| desiredLevel | string | Trình độ mong muốn (beginner, average, good, pro) |
| location | string | Vị trí tìm kiếm |
| latitude | number | Vĩ độ |
| longitude | number | Kinh độ |
| venueId | integer | ID của venue (tùy chọn) |
| desiredTime | string | Thời gian mong muốn (ISO 8601) |
| pricePerPerson | number | Giá mỗi người |
| maxPlayers | integer | Số người chơi tối đa |
| note | string | Ghi chú thêm |
| status | string | Trạng thái (open, matched, cancelled, done) |
| createdAt | string | Thời gian tạo (ISO 8601) |
| updatedAt | string | Thời gian cập nhật (ISO 8601) |
| user | object | Thông tin người tạo |
| venue | object | Thông tin venue |
| matchGroups | array | Danh sách nhóm match |

## Status Codes
- `open`: Bài đăng đang mở, chưa có đủ người
- `matched`: Đã tìm đủ người chơi
- `cancelled`: Bài đăng đã bị hủy
- `done`: Trận đấu đã hoàn thành

## Notes
- Tất cả thời gian đều sử dụng định dạng ISO 8601 (UTC)
- Tọa độ sử dụng hệ thống WGS84
- Bán kính tìm kiếm tính bằng kilomet
- Chỉ chủ bài đăng mới có thể cập nhật hoặc xóa bài đăng của mình
