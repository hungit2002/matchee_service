package entity

import "time"

// CreateBookingRequest represents the request to create a booking
type CreateBookingRequest struct {
	MatchGroupID *uint64    `json:"matchGroupId" binding:"omitempty"`
	VenueID      *uint64    `json:"venueId" binding:"omitempty"`
	CourtID      *uint64    `json:"courtId" binding:"omitempty"`
	StartTime    *time.Time `json:"startTime" binding:"required"`
	EndTime      *time.Time `json:"endTime" binding:"required"`
	TotalPrice   *float64   `json:"totalPrice" binding:"omitempty,min=0"`
}

// UpdateBookingStatusRequest represents the request to update booking status
type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=reserved paid completed cancelled"`
}

// BookingListRequest represents filters for listing bookings
type BookingListRequest struct {
	UserID  *uint64 `json:"userId" form:"userId"`
	VenueID *uint64 `json:"venueId" form:"venueId"`
	Status  *string `json:"status" form:"status" binding:"omitempty,oneof=reserved paid completed cancelled"`
	Page    int     `json:"page" form:"page"`
	Limit   int     `json:"limit" form:"limit"`
}

// BookingResponse represents booking response
type BookingResponse struct {
	ID           uint64      `json:"id"`
	MatchGroupID *uint64     `json:"matchGroupId,omitempty"`
	VenueID      *uint64     `json:"venueId,omitempty"`
	CourtID      *uint64     `json:"courtId,omitempty"`
	StartTime    *string     `json:"startTime,omitempty"`
	EndTime      *string     `json:"endTime,omitempty"`
	TotalPrice   *float64    `json:"totalPrice,omitempty"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"createdAt"`
	UpdatedAt    string      `json:"updatedAt"`
	MatchGroup   *MatchGroup `json:"matchGroup,omitempty"`
	Venue        *Venue      `json:"venue,omitempty"`
	Court        *Court      `json:"court,omitempty"`
}

// BookingListResponse holds paginated list of bookings
type BookingListResponse struct {
	Bookings   []*BookingResponse `json:"bookings"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}
