# Admin API Documentation

## Overview
Admin management APIs for system administration, user management, venue approval, and analytics.

## Authentication
All admin endpoints require admin role authentication. Use the `admin_access_token` in the Authorization header.

## Endpoints

### 1. Get Revenue Statistics
**GET** `/api/v1/admin/statistics/revenue`

Retrieves revenue statistics for admin dashboard.

**Headers:**
- `Authorization: Bearer {admin_access_token}`

**Query Parameters:**
- `fromDate` (optional): Start date in YYYY-MM-DD format
- `toDate` (optional): End date in YYYY-MM-DD format

**Response:**
```json
{
  "success": true,
  "message": "Revenue statistics retrieved successfully",
  "data": {
    "period": "2024-01-01 to 2024-01-31",
    "data": [
      {
        "date": "2024-01-25",
        "totalRevenue": 1200000,
        "bookingCount": 15,
        "averageOrderValue": 80000
      }
    ],
    "total": 1200000
  }
}
```

### 2. Get Users List
**GET** `/api/v1/admin/users`

Retrieves paginated list of users for admin management.

**Headers:**
- `Authorization: Bearer {admin_access_token}`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `status` (optional): Filter by status (`active`, `locked`)
- `role` (optional): Filter by role (`user`, `venue_owner`, `admin`)
- `search` (optional): Search by email or name

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": 1,
        "email": "user@example.com",
        "fullName": "John Doe",
        "phone": "+84901234567",
        "role": "user",
        "status": "active",
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-25T10:00:00Z",
        "lastLoginAt": "2024-01-25T09:00:00Z",
        "bookingCount": 5,
        "totalSpent": 400000
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 20,
    "totalPages": 5
  }
}
```

### 3. Lock User
**PUT** `/api/v1/admin/users/{id}/lock`

Locks or unlocks a user account.

**Headers:**
- `Authorization: Bearer {admin_access_token}`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "locked": true,
  "reason": "Violation of terms of service"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User locked successfully",
  "data": null
}
```

### 4. Approve Venue
**PUT** `/api/v1/admin/venues/{id}/approve`

Approves or rejects a venue for operation.

**Headers:**
- `Authorization: Bearer {admin_access_token}`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "approved": true,
  "reason": "Venue meets all requirements"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Venue approved successfully",
  "data": null
}
```

### 5. Get Feedbacks
**GET** `/api/v1/admin/feedbacks`

Retrieves paginated list of feedbacks for admin review.

**Headers:**
- `Authorization: Bearer {admin_access_token}`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `rating` (optional): Filter by rating (1-5)
- `fromDate` (optional): Start date in YYYY-MM-DD format
- `toDate` (optional): End date in YYYY-MM-DD format

**Response:**
```json
{
  "success": true,
  "message": "Feedbacks retrieved successfully",
  "data": {
    "feedbacks": [
      {
        "id": 1,
        "bookingId": 1,
        "fromUserId": 1,
        "toUserId": 2,
        "rating": 5,
        "comment": "Great player! Very fair and skilled.",
        "createdAt": "2024-01-25T10:00:00Z",
        "fromUser": {
          "id": 1,
          "email": "user1@example.com",
          "fullName": "John Doe"
        },
        "toUser": {
          "id": 2,
          "email": "user2@example.com",
          "fullName": "Jane Smith"
        },
        "booking": {
          "id": 1,
          "venueId": 1,
          "courtId": 1,
          "startTime": "2024-01-25T18:00:00Z",
          "endTime": "2024-01-25T20:00:00Z",
          "totalPrice": 400000,
          "status": "completed"
        }
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 20,
    "totalPages": 3
  }
}
```

## Admin Permissions

### Required Role
- All admin endpoints require `admin` role in JWT token
- Non-admin users will receive 403 Forbidden response

### Admin Capabilities
1. **Revenue Analytics**: View system revenue statistics
2. **User Management**: View, search, and lock/unlock users
3. **Venue Approval**: Approve or reject venue registrations
4. **Feedback Monitoring**: Review user feedback and ratings

## Business Rules

### User Management
- **Locking**: Soft delete user account (sets `deleted_at`)
- **Unlocking**: Restore user account (clears `deleted_at`)
- **Search**: Search by email or full name
- **Filtering**: Filter by status (active/locked) and role

### Venue Management
- **Approval**: Set venue `is_active` status
- **Rejection**: Disable venue operation
- **Reason Tracking**: Optional reason for approval/rejection

### Revenue Statistics
- **Data Source**: Based on successful payments
- **Date Range**: Flexible date filtering
- **Metrics**: Daily revenue, booking count, average order value

### Feedback Monitoring
- **Rating Filter**: Filter by specific star ratings
- **Date Range**: Filter by creation date
- **Full Details**: Includes user and booking information

## Error Responses

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Invalid token",
  "error": "JWT validation failed"
}
```

### 403 Forbidden
```json
{
  "success": false,
  "message": "Admin access required",
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "User not found",
  "error": "User with ID 999 does not exist"
}
```

## Security Considerations

1. **Admin Token**: Use dedicated admin JWT tokens with admin role
2. **Audit Logging**: All admin actions should be logged
3. **Rate Limiting**: Implement rate limiting for admin endpoints
4. **IP Whitelisting**: Consider IP whitelisting for admin access
5. **Session Management**: Implement proper session timeout for admin users

## Use Cases

1. **System Monitoring**: Track revenue and user activity
2. **User Support**: Lock problematic users, review accounts
3. **Venue Management**: Approve new venues, manage venue status
4. **Quality Control**: Monitor user feedback, identify issues
5. **Business Intelligence**: Analyze revenue trends and user behavior
