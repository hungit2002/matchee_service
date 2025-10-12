# Matchee Player Profile API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

### Player Profile Endpoints

#### 1. Create/Update Player Profile
**POST** `/player/profile`

Tạo hoặc cập nhật hồ sơ người chơi với trình độ, vị trí, giới tính.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "level": "good",
  "gender": "male",
  "preferredLocation": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "bio": "Passionate tennis player looking for matches"
}
```

**Response (200):**
```json
{
  "id": 1,
  "userId": 1,
  "level": "good",
  "gender": "male",
  "preferredLocation": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "bio": "Passionate tennis player looking for matches",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
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
```

#### 2. Get Current User's Player Profile
**GET** `/player/profile`

Lấy hồ sơ người chơi của user hiện tại.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "id": 1,
  "userId": 1,
  "level": "good",
  "gender": "male",
  "preferredLocation": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "bio": "Passionate tennis player looking for matches",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "user": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com"
  }
}
```

#### 3. Get Player Profile by ID
**GET** `/player/profile/{id}`

Lấy chi tiết hồ sơ người chơi theo ID.

**Parameters:**
- `id` (path): Player Profile ID

**Response (200):**
```json
{
  "id": 1,
  "userId": 1,
  "level": "good",
  "gender": "male",
  "preferredLocation": "Ho Chi Minh City",
  "latitude": 10.8231,
  "longitude": 106.6297,
  "bio": "Passionate tennis player looking for matches",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "user": {
    "id": 1,
    "fullName": "Nguyễn Văn A",
    "phone": "0123456789",
    "email": "user@example.com"
  }
}
```

#### 4. Get Player Suggestions
**GET** `/player/suggestions`

Gợi ý đối thủ phù hợp dựa theo level + location.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `level` (optional): Player level (`beginner`, `average`, `good`, `pro`)
- `latitude` (optional): Latitude for location-based search
- `longitude` (optional): Longitude for location-based search
- `radius` (optional): Search radius in kilometers (default: 10)
- `limit` (optional): Maximum number of suggestions (default: 20)

**Example Requests:**
```bash
# Get suggestions by level
GET /api/v1/player/suggestions?level=good&limit=5

# Get suggestions by location
GET /api/v1/player/suggestions?latitude=10.8231&longitude=106.6297&radius=5&limit=5

# Get suggestions with combined criteria
GET /api/v1/player/suggestions?level=average&latitude=10.8231&longitude=106.6297&radius=10&limit=10
```

**Response (200):**
```json
{
  "players": [
    {
      "id": 2,
      "userId": 2,
      "level": "good",
      "gender": "female",
      "preferredLocation": "District 1, Ho Chi Minh City",
      "latitude": 10.7769,
      "longitude": 106.7009,
      "bio": "Tennis enthusiast looking for practice partners",
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
```

## Field Descriptions

### Player Level Options
- `beginner`: Người mới bắt đầu
- `average`: Trình độ trung bình
- `good`: Trình độ khá
- `pro`: Trình độ chuyên nghiệp

### Gender Options
- `male`: Nam
- `female`: Nữ
- `other`: Khác

### Location Fields
- `latitude`: Vĩ độ (-90 đến 90)
- `longitude`: Kinh độ (-180 đến 180)
- `preferredLocation`: Tên địa điểm ưa thích
- `radius`: Bán kính tìm kiếm (km)

## Error Responses

### Common Error Format
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes
- `200` - Success
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid or missing token)
- `404` - Not Found (profile not found)
- `500` - Internal Server Error

## Test Examples

### 1. Create Player Profile
```bash
curl -X POST "http://localhost:8080/api/v1/player/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "level": "good",
    "gender": "male",
    "preferredLocation": "Ho Chi Minh City",
    "latitude": 10.8231,
    "longitude": 106.6297,
    "bio": "Passionate tennis player"
  }'
```

### 2. Get Player Suggestions
```bash
curl -X GET "http://localhost:8080/api/v1/player/suggestions?level=good&limit=5" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 3. Get Nearby Players
```bash
curl -X GET "http://localhost:8080/api/v1/player/suggestions?latitude=10.8231&longitude=106.6297&radius=5" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Database Schema

### player_profiles Table
```sql
CREATE TABLE player_profiles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    level ENUM('beginner','average','good','pro') DEFAULT 'average',
    gender ENUM('male','female','other') NULL,
    preferred_location VARCHAR(255) NULL,
    latitude DECIMAL(10,7) NULL,
    longitude DECIMAL(10,7) NULL,
    bio TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_level (level),
    INDEX idx_location (latitude, longitude)
);
```

## Features

1. **Level-based Matching**: Tìm kiếm người chơi cùng trình độ
2. **Location-based Matching**: Tìm kiếm người chơi gần vị trí
3. **Combined Criteria**: Kết hợp level và location
4. **Distance Calculation**: Tính toán khoảng cách chính xác
5. **Flexible Search**: Hỗ trợ nhiều tiêu chí tìm kiếm
6. **User Privacy**: Chỉ hiển thị thông tin cần thiết
