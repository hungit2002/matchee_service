package entity

import "time"

// CreateMatchPostRequest represents the request payload for creating match posts
type CreateMatchPostRequest struct {
	DesiredLevel   *string    `json:"desiredLevel" binding:"omitempty,oneof=beginner average good pro"`
	Location       *string    `json:"location" binding:"omitempty,max=255"`
	Latitude       *float64   `json:"latitude" binding:"omitempty"`
	Longitude      *float64   `json:"longitude" binding:"omitempty"`
	VenueID        *uint64    `json:"venueId" binding:"omitempty"`
	DesiredTime    *time.Time `json:"desiredTime" binding:"omitempty"`
	PricePerPerson *float64   `json:"pricePerPerson" binding:"omitempty,min=0"`
	MaxPlayers     int        `json:"maxPlayers" binding:"min=2,max=20"`
	Note           *string    `json:"note" binding:"omitempty,max=1000"`
}

// UpdateMatchPostRequest represents the request payload for updating match posts
type UpdateMatchPostRequest struct {
	DesiredLevel   *string    `json:"desiredLevel" binding:"omitempty,oneof=beginner average good pro"`
	Location       *string    `json:"location" binding:"omitempty,max=255"`
	Latitude       *float64   `json:"latitude" binding:"omitempty"`
	Longitude      *float64   `json:"longitude" binding:"omitempty"`
	VenueID        *uint64    `json:"venueId" binding:"omitempty"`
	DesiredTime    *time.Time `json:"desiredTime" binding:"omitempty"`
	PricePerPerson *float64   `json:"pricePerPerson" binding:"omitempty,min=0"`
	MaxPlayers     *int       `json:"maxPlayers" binding:"omitempty,min=2,max=20"`
	Note           *string    `json:"note" binding:"omitempty,max=1000"`
	Status         *string    `json:"status" binding:"omitempty,oneof=open matched cancelled done"`
}

// MatchPostResponse represents the response payload for match post
type MatchPostResponse struct {
	ID             uint64        `json:"id"`
	UserID         uint64        `json:"userId"`
	DesiredLevel   *string       `json:"desiredLevel,omitempty"`
	Location       *string       `json:"location,omitempty"`
	Latitude       *float64      `json:"latitude,omitempty"`
	Longitude      *float64      `json:"longitude,omitempty"`
	VenueID        *uint64       `json:"venueId,omitempty"`
	DesiredTime    *string       `json:"desiredTime,omitempty"`
	PricePerPerson *float64      `json:"pricePerPerson,omitempty"`
	MaxPlayers     int           `json:"maxPlayers"`
	Note           *string       `json:"note,omitempty"`
	Status         string        `json:"status"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
	User           *User         `json:"user,omitempty"`
	Venue          *Venue        `json:"venue,omitempty"`
	MatchGroups    []*MatchGroup `json:"matchGroups,omitempty"`
}

// MatchPostListRequest represents the request payload for listing match posts
type MatchPostListRequest struct {
	DesiredLevel *string  `json:"desiredLevel" form:"desiredLevel" binding:"omitempty,oneof=beginner average good pro"`
	Location     *string  `json:"location" form:"location" binding:"omitempty"`
	Latitude     *float64 `json:"latitude" form:"latitude" binding:"omitempty"`
	Longitude    *float64 `json:"longitude" form:"longitude" binding:"omitempty"`
	Radius       *float64 `json:"radius" form:"radius" binding:"omitempty,min=0"`
	VenueID      *uint64  `json:"venueId" form:"venueId" binding:"omitempty"`
	Status       *string  `json:"status" form:"status" binding:"omitempty,oneof=open matched cancelled done"`
	Page         int      `json:"page" form:"page"`
	Limit        int      `json:"limit" form:"limit"`
}

// MatchPostListResponse represents the response payload for match post list
type MatchPostListResponse struct {
	MatchPosts []*MatchPostResponse `json:"matchPosts"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"totalPages"`
}

// MatchPostSuggestRequest represents the request payload for match post suggestions
type MatchPostSuggestRequest struct {
	UserID       uint64   `json:"userId" binding:"required"`
	DesiredLevel *string  `json:"desiredLevel" binding:"omitempty,oneof=beginner average good pro"`
	Latitude     *float64 `json:"latitude" binding:"omitempty"`
	Longitude    *float64 `json:"longitude" binding:"omitempty"`
	Radius       *float64 `json:"radius" binding:"omitempty,min=0"`
	Limit        int      `json:"limit" binding:"omitempty,min=1,max=50"`
}

// MatchPostSuggestResponse represents the response payload for match post suggestions
type MatchPostSuggestResponse struct {
	Suggestions []*MatchPostResponse `json:"suggestions"`
	Total       int64                `json:"total"`
}
