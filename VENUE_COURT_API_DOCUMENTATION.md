# Venue & Court API Documentation

## 📍 Venue Management APIs

### 1. Create Venue
**POST** `/api/v1/venues`

Tạo sân mới (dành cho chủ sân)

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Tennis Center Ho Chi Minh",
  "description": "Premium tennis center with multiple courts",
  "address": "123 Nguyen Hue, District 1, Ho Chi Minh City",
  "latitude": 10.7769,
  "longitude": 106.7009,
  "phone": "+84901234567",
  "email": "info@tenniscenter.com",
  "website": "https://tenniscenter.com",
  "pricePerHour": 200000,
  "isActive": true
}
```

**Response:**
```json
{
  "status": "success",
  "code": 201,
  "message": "Venue created successfully",
  "data": {
    "id": 1,
    "ownerId": 1,
    "name": "Tennis Center Ho Chi Minh",
    "description": "Premium tennis center with multiple courts",
    "address": "123 Nguyen Hue, District 1, Ho Chi Minh City",
    "latitude": 10.7769,
    "longitude": 106.7009,
    "phone": "+84901234567",
    "email": "info@tenniscenter.com",
    "website": "https://tenniscenter.com",
    "pricePerHour": 200000,
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T09:00:00Z",
    "owner": {
      "id": 1,
      "fullName": "John Doe",
      "phone": "+84901234567",
      "email": "john@example.com"
    }
  }
}
```

### 2. Get Venues
**GET** `/api/v1/venues`

Lấy danh sách sân với các bộ lọc tùy chọn

**Query Parameters:**
- `latitude` (number, optional): Vĩ độ để lọc theo vị trí
- `longitude` (number, optional): Kinh độ để lọc theo vị trí  
- `radius` (number, optional): Bán kính tính bằng km
- `minPrice` (number, optional): Giá tối thiểu mỗi giờ
- `maxPrice` (number, optional): Giá tối đa mỗi giờ
- `isActive` (boolean, optional): Lọc theo trạng thái hoạt động
- `page` (int, optional): Số trang (mặc định: 1)
- `limit` (int, optional): Số item mỗi trang (mặc định: 20)

**Example Request:**
```
GET /api/v1/venues?latitude=10.7769&longitude=106.7009&radius=10&minPrice=100000&maxPrice=500000&isActive=true&page=1&limit=20
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Venues retrieved successfully",
  "data": {
    "venues": [
      {
        "id": 1,
        "ownerId": 1,
        "name": "Tennis Center Ho Chi Minh",
        "description": "Premium tennis center with multiple courts",
        "address": "123 Nguyen Hue, District 1, Ho Chi Minh City",
        "latitude": 10.7769,
        "longitude": 106.7009,
        "phone": "+84901234567",
        "email": "info@tenniscenter.com",
        "website": "https://tenniscenter.com",
        "pricePerHour": 200000,
        "isActive": true,
        "createdAt": "2024-01-15T09:00:00Z",
        "updatedAt": "2024-01-15T09:00:00Z",
        "owner": {
          "id": 1,
          "fullName": "John Doe",
          "phone": "+84901234567",
          "email": "john@example.com"
        },
        "courts": [
          {
            "id": 1,
            "venueId": 1,
            "name": "Court 1",
            "description": "Professional hard court",
            "courtType": "outdoor",
            "surface": "hard",
            "isActive": true,
            "createdAt": "2024-01-15T09:00:00Z",
            "updatedAt": "2024-01-15T09:00:00Z"
          }
        ]
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20,
    "totalPages": 1
  }
}
```

### 3. Get Venue by ID
**GET** `/api/v1/venues/{id}`

Lấy chi tiết sân theo ID

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Venue retrieved successfully",
  "data": {
    "id": 1,
    "ownerId": 1,
    "name": "Tennis Center Ho Chi Minh",
    "description": "Premium tennis center with multiple courts",
    "address": "123 Nguyen Hue, District 1, Ho Chi Minh City",
    "latitude": 10.7769,
    "longitude": 106.7009,
    "phone": "+84901234567",
    "email": "info@tenniscenter.com",
    "website": "https://tenniscenter.com",
    "pricePerHour": 200000,
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T09:00:00Z",
    "owner": {
      "id": 1,
      "fullName": "John Doe",
      "phone": "+84901234567",
      "email": "john@example.com"
    },
    "courts": [
      {
        "id": 1,
        "venueId": 1,
        "name": "Court 1",
        "description": "Professional hard court",
        "courtType": "outdoor",
        "surface": "hard",
        "isActive": true,
        "createdAt": "2024-01-15T09:00:00Z",
        "updatedAt": "2024-01-15T09:00:00Z"
      }
    ]
  }
}
```

### 4. Get My Venues
**GET** `/api/v1/venues/my`

Lấy danh sách sân của user hiện tại

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `page` (int, optional): Số trang (mặc định: 1)
- `limit` (int, optional): Số item mỗi trang (mặc định: 20)

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "My venues retrieved successfully",
  "data": {
    "venues": [
      {
        "id": 1,
        "ownerId": 1,
        "name": "Tennis Center Ho Chi Minh",
        "description": "Premium tennis center with multiple courts",
        "address": "123 Nguyen Hue, District 1, Ho Chi Minh City",
        "latitude": 10.7769,
        "longitude": 106.7009,
        "phone": "+84901234567",
        "email": "info@tenniscenter.com",
        "website": "https://tenniscenter.com",
        "pricePerHour": 200000,
        "isActive": true,
        "createdAt": "2024-01-15T09:00:00Z",
        "updatedAt": "2024-01-15T09:00:00Z",
        "owner": {
          "id": 1,
          "fullName": "John Doe",
          "phone": "+84901234567",
          "email": "john@example.com"
        },
        "courts": []
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20,
    "totalPages": 1
  }
}
```

### 5. Update Venue
**PUT** `/api/v1/venues/{id}`

Cập nhật thông tin sân

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Updated Tennis Center",
  "description": "Updated premium tennis center",
  "pricePerHour": 250000
}
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Venue updated successfully",
  "data": {
    "id": 1,
    "ownerId": 1,
    "name": "Updated Tennis Center",
    "description": "Updated premium tennis center",
    "address": "123 Nguyen Hue, District 1, Ho Chi Minh City",
    "latitude": 10.7769,
    "longitude": 106.7009,
    "phone": "+84901234567",
    "email": "info@tenniscenter.com",
    "website": "https://tenniscenter.com",
    "pricePerHour": 250000,
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

### 6. Delete Venue
**DELETE** `/api/v1/venues/{id}`

Xóa sân

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Venue deleted successfully",
  "data": null
}
```

## 🏟️ Court Management APIs

### 1. Create Court
**POST** `/api/v1/venues/{venueId}/courts`

Thêm sân con cho venue

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Court 1",
  "description": "Professional hard court",
  "courtType": "outdoor",
  "surface": "hard",
  "isActive": true
}
```

**Response:**
```json
{
  "status": "success",
  "code": 201,
  "message": "Court created successfully",
  "data": {
    "id": 1,
    "venueId": 1,
    "name": "Court 1",
    "description": "Professional hard court",
    "courtType": "outdoor",
    "surface": "hard",
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T09:00:00Z",
    "venue": {
      "id": 1,
      "name": "Tennis Center Ho Chi Minh",
      "address": "123 Nguyen Hue, District 1, Ho Chi Minh City"
    }
  }
}
```

### 2. Get Courts by Venue
**GET** `/api/v1/venues/{venueId}/courts`

Lấy danh sách sân con của venue

**Query Parameters:**
- `page` (int, optional): Số trang (mặc định: 1)
- `limit` (int, optional): Số item mỗi trang (mặc định: 20)

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Courts retrieved successfully",
  "data": {
    "courts": [
      {
        "id": 1,
        "venueId": 1,
        "name": "Court 1",
        "description": "Professional hard court",
        "courtType": "outdoor",
        "surface": "hard",
        "isActive": true,
        "createdAt": "2024-01-15T09:00:00Z",
        "updatedAt": "2024-01-15T09:00:00Z",
        "venue": {
          "id": 1,
          "name": "Tennis Center Ho Chi Minh",
          "address": "123 Nguyen Hue, District 1, Ho Chi Minh City"
        },
        "slots": []
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20
  }
}
```

### 3. Get Court by ID
**GET** `/api/v1/courts/{id}`

Lấy chi tiết sân con theo ID

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Court retrieved successfully",
  "data": {
    "id": 1,
    "venueId": 1,
    "name": "Court 1",
    "description": "Professional hard court",
    "courtType": "outdoor",
    "surface": "hard",
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T09:00:00Z",
    "venue": {
      "id": 1,
      "name": "Tennis Center Ho Chi Minh",
      "address": "123 Nguyen Hue, District 1, Ho Chi Minh City"
    },
    "slots": [
      {
        "id": 1,
        "courtId": 1,
        "startTime": "2024-01-15T09:00:00Z",
        "endTime": "2024-01-15T10:00:00Z",
        "price": 200000,
        "isActive": true,
        "createdAt": "2024-01-15T09:00:00Z",
        "updatedAt": "2024-01-15T09:00:00Z"
      }
    ]
  }
}
```

### 4. Update Court
**PUT** `/api/v1/courts/{id}`

Cập nhật thông tin sân con

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Updated Court 1",
  "description": "Updated professional hard court",
  "isActive": true
}
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Court updated successfully",
  "data": {
    "id": 1,
    "venueId": 1,
    "name": "Updated Court 1",
    "description": "Updated professional hard court",
    "courtType": "outdoor",
    "surface": "hard",
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

### 5. Delete Court
**DELETE** `/api/v1/courts/{id}`

Xóa sân con

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Court deleted successfully",
  "data": null
}
```

## ⏰ Availability Slots APIs

### 1. Create Slot
**POST** `/api/v1/courts/{courtId}/slots`

Tạo khung giờ rảnh cho sân con

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "startTime": "2024-01-15T09:00:00Z",
  "endTime": "2024-01-15T10:00:00Z",
  "price": 200000,
  "isActive": true
}
```

**Response:**
```json
{
  "status": "success",
  "code": 201,
  "message": "Availability slot created successfully",
  "data": {
    "id": 1,
    "courtId": 1,
    "startTime": "2024-01-15T09:00:00Z",
    "endTime": "2024-01-15T10:00:00Z",
    "price": 200000,
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T09:00:00Z",
    "court": {
      "id": 1,
      "name": "Court 1",
      "venueId": 1
    }
  }
}
```

### 2. Get Slots by Court
**GET** `/api/v1/courts/{courtId}/slots`

Lấy danh sách khung giờ rảnh của sân con

**Query Parameters:**
- `startDate` (string, optional): Ngày bắt đầu (RFC3339 format)
- `endDate` (string, optional): Ngày kết thúc (RFC3339 format)
- `isActive` (boolean, optional): Lọc theo trạng thái hoạt động
- `page` (int, optional): Số trang (mặc định: 1)
- `limit` (int, optional): Số item mỗi trang (mặc định: 20)

**Example Request:**
```
GET /api/v1/courts/1/slots?startDate=2024-01-15T00:00:00Z&endDate=2024-01-15T23:59:59Z&isActive=true&page=1&limit=20
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Availability slots retrieved successfully",
  "data": {
    "slots": [
      {
        "id": 1,
        "courtId": 1,
        "startTime": "2024-01-15T09:00:00Z",
        "endTime": "2024-01-15T10:00:00Z",
        "price": 200000,
        "isActive": true,
        "createdAt": "2024-01-15T09:00:00Z",
        "updatedAt": "2024-01-15T09:00:00Z",
        "court": {
          "id": 1,
          "name": "Court 1",
          "venueId": 1
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

### 3. Get Slot by ID
**GET** `/api/v1/slots/{id}`

Lấy chi tiết khung giờ rảnh theo ID

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Availability slot retrieved successfully",
  "data": {
    "id": 1,
    "courtId": 1,
    "startTime": "2024-01-15T09:00:00Z",
    "endTime": "2024-01-15T10:00:00Z",
    "price": 200000,
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T09:00:00Z",
    "court": {
      "id": 1,
      "name": "Court 1",
      "venueId": 1
    }
  }
}
```

### 4. Update Slot
**PUT** `/api/v1/slots/{id}`

Cập nhật thông tin khung giờ rảnh

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "startTime": "2024-01-15T10:00:00Z",
  "endTime": "2024-01-15T11:00:00Z",
  "price": 250000,
  "isActive": true
}
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Availability slot updated successfully",
  "data": {
    "id": 1,
    "courtId": 1,
    "startTime": "2024-01-15T10:00:00Z",
    "endTime": "2024-01-15T11:00:00Z",
    "price": 250000,
    "isActive": true,
    "createdAt": "2024-01-15T09:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  }
}
```

### 5. Delete Slot
**DELETE** `/api/v1/slots/{id}`

Xóa khung giờ rảnh

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": "success",
  "code": 200,
  "message": "Availability slot deleted successfully",
  "data": null
}
```

## 🔧 Data Models

### Venue Model
```json
{
  "id": "number",
  "ownerId": "number",
  "name": "string",
  "description": "string (optional)",
  "address": "string",
  "latitude": "number",
  "longitude": "number",
  "phone": "string",
  "email": "string (optional)",
  "website": "string (optional)",
  "pricePerHour": "number",
  "isActive": "boolean",
  "createdAt": "string (RFC3339)",
  "updatedAt": "string (RFC3339)",
  "owner": "User object (optional)",
  "courts": "Court array (optional)"
}
```

### Court Model
```json
{
  "id": "number",
  "venueId": "number",
  "name": "string",
  "description": "string (optional)",
  "courtType": "string (enum: indoor, outdoor)",
  "surface": "string (enum: hard, clay, grass, synthetic)",
  "isActive": "boolean",
  "createdAt": "string (RFC3339)",
  "updatedAt": "string (RFC3339)",
  "venue": "Venue object (optional)",
  "slots": "Slot array (optional)"
}
```

### Slot Model
```json
{
  "id": "number",
  "courtId": "number",
  "startTime": "string (RFC3339)",
  "endTime": "string (RFC3339)",
  "price": "number (optional)",
  "isActive": "boolean",
  "createdAt": "string (RFC3339)",
  "updatedAt": "string (RFC3339)",
  "court": "Court object (optional)"
}
```

## 🚨 Error Responses

### 400 Bad Request
```json
{
  "status": "fail",
  "code": 400,
  "message": "Invalid request data"
}
```

### 401 Unauthorized
```json
{
  "status": "fail",
  "code": 401,
  "message": "User not authenticated"
}
```

### 404 Not Found
```json
{
  "status": "fail",
  "code": 404,
  "message": "Venue not found"
}
```

### 409 Conflict
```json
{
  "status": "fail",
  "code": 409,
  "message": "Slot overlaps with existing slot"
}
```

## 📝 Notes

1. **Authentication**: Hầu hết các API yêu cầu authentication token
2. **Pagination**: Tất cả list APIs đều hỗ trợ pagination
3. **Location Filtering**: Venue APIs hỗ trợ lọc theo vị trí địa lý
4. **Time Validation**: Slot APIs kiểm tra thời gian hợp lệ và không trùng lặp
5. **Response Format**: Tất cả APIs đều sử dụng format response chuẩn với `status`, `code`, `message`, và `data`
6. **Auto-save**: Postman collection tự động lưu `venue_id`, `court_id`, và `slot_id` từ responses
