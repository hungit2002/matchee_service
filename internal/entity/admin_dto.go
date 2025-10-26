package entity

// RevenueStatistics represents revenue statistics for admin
type RevenueStatistics struct {
	Date              string  `json:"date"`
	TotalRevenue      float64 `json:"totalRevenue"`
	BookingCount      int64   `json:"bookingCount"`
	AverageOrderValue float64 `json:"averageOrderValue"`
}

// RevenueStatisticsResponse represents revenue statistics response
type RevenueStatisticsResponse struct {
	Period string               `json:"period"`
	Data   []*RevenueStatistics `json:"data"`
	Total  float64              `json:"total"`
}

// AdminUserListRequest represents filters for listing users
type AdminUserListRequest struct {
	Page   int     `json:"page" form:"page"`
	Limit  int     `json:"limit" form:"limit"`
	Status *string `json:"status" form:"status" binding:"omitempty,oneof=active locked"`
	Role   *string `json:"role" form:"role" binding:"omitempty,oneof=user venue_owner admin"`
	Search *string `json:"search" form:"search"`
}

// AdminUserListResponse represents paginated user list for admin
type AdminUserListResponse struct {
	Users      []*AdminUserResponse `json:"users"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"totalPages"`
}

// AdminUserResponse represents user data for admin
type AdminUserResponse struct {
	ID           uint64  `json:"id"`
	Email        string  `json:"email"`
	FullName     *string `json:"fullName,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Role         string  `json:"role"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
	LastLoginAt  *string `json:"lastLoginAt,omitempty"`
	BookingCount int64   `json:"bookingCount"`
	TotalSpent   float64 `json:"totalSpent"`
}

// LockUserRequest represents request to lock/unlock user
type LockUserRequest struct {
	Locked bool    `json:"locked" binding:"required"`
	Reason *string `json:"reason" binding:"omitempty,max=500"`
}

// VenueApprovalRequest represents request to approve venue
type VenueApprovalRequest struct {
	Approved bool    `json:"approved" binding:"required"`
	Reason   *string `json:"reason" binding:"omitempty,max=500"`
}

// AdminFeedbackListRequest represents filters for listing feedbacks
type AdminFeedbackListRequest struct {
	Page     int     `json:"page" form:"page"`
	Limit    int     `json:"limit" form:"limit"`
	Rating   *int    `json:"rating" form:"rating" binding:"omitempty,min=1,max=5"`
	FromDate *string `json:"fromDate" form:"fromDate"`
	ToDate   *string `json:"toDate" form:"toDate"`
}

// AdminFeedbackListResponse represents paginated feedback list for admin
type AdminFeedbackListResponse struct {
	Feedbacks  []*AdminFeedbackResponse `json:"feedbacks"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"totalPages"`
}

// AdminFeedbackResponse represents feedback data for admin
type AdminFeedbackResponse struct {
	ID         uint64   `json:"id"`
	BookingID  uint64   `json:"bookingId"`
	FromUserID uint64   `json:"fromUserId"`
	ToUserID   uint64   `json:"toUserId"`
	Rating     *int     `json:"rating,omitempty"`
	Comment    *string  `json:"comment,omitempty"`
	CreatedAt  string   `json:"createdAt"`
	FromUser   *User    `json:"fromUser,omitempty"`
	ToUser     *User    `json:"toUser,omitempty"`
	Booking    *Booking `json:"booking,omitempty"`
}
