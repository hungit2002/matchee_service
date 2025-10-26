package entity

// CreateNotificationRequest is used to create a new notification
type CreateNotificationRequest struct {
	UserID  uint64  `json:"userId" binding:"required"`
	Title   *string `json:"title" binding:"omitempty,max=255"`
	Message *string `json:"message" binding:"omitempty"`
}

// NotificationResponse is the DTO returned to clients
type NotificationResponse struct {
	ID        uint64  `json:"id"`
	UserID    uint64  `json:"userId"`
	Title     *string `json:"title,omitempty"`
	Message   *string `json:"message,omitempty"`
	IsRead    bool    `json:"isRead"`
	CreatedAt string  `json:"createdAt"`
}

// NotificationListResponse contains a list of notifications
type NotificationListResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Total         int64                   `json:"total"`
	Page          int                     `json:"page"`
	Limit         int                     `json:"limit"`
	TotalPages    int                     `json:"totalPages"`
}
