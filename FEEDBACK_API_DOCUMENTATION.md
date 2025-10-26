# Feedback API Documentation

## Overview
Feedback management APIs for post-match evaluations and player ratings.

## Endpoints

### 1. Create Feedback
**POST** `/api/v1/feedbacks`

Creates feedback for a user after completing a match.

**Headers:**
- `Authorization: Bearer {access_token}`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "bookingId": 1,
  "toUserId": 2,
  "rating": 5,
  "comment": "Great player! Very fair and skilled."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Feedback created successfully",
  "data": {
    "id": 1,
    "bookingId": 1,
    "fromUserId": 1,
    "toUserId": 2,
    "rating": 5,
    "comment": "Great player! Very fair and skilled.",
    "createdAt": "2024-01-25T10:00:00Z",
    "updatedAt": "2024-01-25T10:00:00Z",
    "booking": {
      "id": 1,
      "venueId": 1,
      "courtId": 1,
      "startTime": "2024-01-25T18:00:00Z",
      "endTime": "2024-01-25T20:00:00Z",
      "totalPrice": 400000,
      "status": "completed"
    },
    "fromUser": {
      "id": 1,
      "email": "user1@example.com",
      "fullName": "John Doe"
    },
    "toUser": {
      "id": 2,
      "email": "user2@example.com",
      "fullName": "Jane Smith"
    }
  }
}
```

### 2. Get User Feedbacks
**GET** `/api/v1/feedbacks/user/{id}`

Retrieves aggregated feedback summary for a specific user.

**Headers:**
- `Authorization: Bearer {access_token}`

**Response:**
```json
{
  "success": true,
  "message": "User feedbacks retrieved successfully",
  "data": {
    "userId": 2,
    "totalRatings": 15,
    "averageRating": 4.2,
    "ratingCounts": {
      "1": 0,
      "2": 1,
      "3": 2,
      "4": 5,
      "5": 7
    },
    "recentFeedbacks": [
      {
        "id": 1,
        "bookingId": 1,
        "fromUserId": 1,
        "toUserId": 2,
        "rating": 5,
        "comment": "Great player! Very fair and skilled.",
        "createdAt": "2024-01-25T10:00:00Z",
        "updatedAt": "2024-01-25T10:00:00Z",
        "fromUser": {
          "id": 1,
          "email": "user1@example.com",
          "fullName": "John Doe"
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
    ]
  }
}
```

### 3. Get Booking Feedbacks
**GET** `/api/v1/feedbacks/booking/{id}`

Retrieves all feedbacks for a specific booking.

**Headers:**
- `Authorization: Bearer {access_token}`

**Response:**
```json
{
  "success": true,
  "message": "Booking feedbacks retrieved successfully",
  "data": {
    "bookingId": 1,
    "feedbacks": [
      {
        "id": 1,
        "bookingId": 1,
        "fromUserId": 1,
        "toUserId": 2,
        "rating": 5,
        "comment": "Great player! Very fair and skilled.",
        "createdAt": "2024-01-25T10:00:00Z",
        "updatedAt": "2024-01-25T10:00:00Z",
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
    "total": 1
  }
}
```

## Business Rules

### Feedback Creation Rules
1. **Authentication Required**: Only authenticated users can create feedback
2. **Booking Validation**: The booking must exist and be in "completed" status
3. **User Validation**: Both from and to users must exist
4. **Duplicate Prevention**: Only one feedback per booking-user pair is allowed
5. **Rating Range**: Rating must be between 1-5 stars
6. **Comment Length**: Comment is optional but limited to 500 characters

### Feedback Retrieval Rules
1. **Public Access**: User feedback summaries are publicly accessible
2. **Booking Access**: Booking feedbacks are publicly accessible
3. **Aggregation**: User summaries include total ratings, average rating, and rating distribution
4. **Recent Feedbacks**: Limited to last 10 feedbacks in user summary

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request data",
  "error": "validation error details"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "User not authenticated",
  "error": "Invalid or missing access token"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Booking not found",
  "error": "Booking with ID 999 does not exist"
}
```

## Data Models

### Feedback Entity
```go
type Feedback struct {
    ID         uint64     `json:"id"`
    BookingID  uint64     `json:"bookingId"`
    FromUserID uint64     `json:"fromUserId"`
    ToUserID   uint64     `json:"toUserId"`
    Rating     *int       `json:"rating,omitempty"`
    Comment    *string    `json:"comment,omitempty"`
    CreatedAt  string     `json:"createdAt"`
    UpdatedAt  string     `json:"updatedAt"`
    Booking    *Booking   `json:"booking,omitempty"`
    FromUser   *User      `json:"fromUser,omitempty"`
    ToUser     *User      `json:"toUser,omitempty"`
}
```

### User Feedback Summary
```go
type UserFeedbackSummary struct {
    UserID           uint64             `json:"userId"`
    TotalRatings     int64              `json:"totalRatings"`
    AverageRating    float64            `json:"averageRating"`
    RatingCounts     map[int]int64      `json:"ratingCounts"`
    RecentFeedbacks  []*FeedbackResponse `json:"recentFeedbacks"`
}
```

### Booking Feedback List
```go
type BookingFeedbackList struct {
    BookingID uint64             `json:"bookingId"`
    Feedbacks []*FeedbackResponse `json:"feedbacks"`
    Total     int64              `json:"total"`
}
```

## Rating System
- **1 Star**: Poor experience
- **2 Stars**: Below average
- **3 Stars**: Average experience
- **4 Stars**: Good experience
- **5 Stars**: Excellent experience

## Use Cases
1. **Post-Match Evaluation**: Players can rate their opponents after completing a match
2. **Player Reputation**: Build player reputation based on feedback from other players
3. **Match Quality**: Track match quality through aggregated feedback
4. **Community Building**: Encourage fair play and good sportsmanship through feedback system
