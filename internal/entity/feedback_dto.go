package entity

// CreateFeedbackRequest represents the request to create a feedback
type CreateFeedbackRequest struct {
	BookingID uint64  `json:"bookingId" binding:"required"`
	ToUserID  uint64  `json:"toUserId" binding:"required"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"`
	Comment   *string `json:"comment" binding:"omitempty,max=500"`
}

// FeedbackResponse represents feedback response
type FeedbackResponse struct {
	ID         uint64   `json:"id"`
	BookingID  uint64   `json:"bookingId"`
	FromUserID uint64   `json:"fromUserId"`
	ToUserID   uint64   `json:"toUserId"`
	Rating     *int     `json:"rating,omitempty"`
	Comment    *string  `json:"comment,omitempty"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
	Booking    *Booking `json:"booking,omitempty"`
	FromUser   *User    `json:"fromUser,omitempty"`
	ToUser     *User    `json:"toUser,omitempty"`
}

// UserFeedbackSummary represents aggregated feedback for a user
type UserFeedbackSummary struct {
	UserID          uint64              `json:"userId"`
	TotalRatings    int64               `json:"totalRatings"`
	AverageRating   float64             `json:"averageRating"`
	RatingCounts    map[int]int64       `json:"ratingCounts"` // 1-5 star counts
	RecentFeedbacks []*FeedbackResponse `json:"recentFeedbacks"`
}

// BookingFeedbackList represents feedbacks for a specific booking
type BookingFeedbackList struct {
	BookingID uint64              `json:"bookingId"`
	Feedbacks []*FeedbackResponse `json:"feedbacks"`
	Total     int64               `json:"total"`
}
