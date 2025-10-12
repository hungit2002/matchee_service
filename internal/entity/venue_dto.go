package entity

import "time"

// CreateVenueRequest represents the request payload for creating a venue
type CreateVenueRequest struct {
	Name         string  `json:"name" binding:"required,min=3,max=100"`
	Description  *string `json:"description" binding:"omitempty,max=500"`
	Address      string  `json:"address" binding:"required,min=10,max=255"`
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Phone        string  `json:"phone" binding:"required,e164"`
	Email        *string `json:"email" binding:"omitempty,email"`
	Website      *string `json:"website" binding:"omitempty,url"`
	PricePerHour float64 `json:"pricePerHour" binding:"required,min=0"`
	IsActive     bool    `json:"isActive"`
}

// UpdateVenueRequest represents the request payload for updating a venue
type UpdateVenueRequest struct {
	Name         *string  `json:"name" binding:"omitempty,min=3,max=100"`
	Description  *string  `json:"description" binding:"omitempty,max=500"`
	Address      *string  `json:"address" binding:"omitempty,min=10,max=255"`
	Latitude     *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Phone        *string  `json:"phone" binding:"omitempty,e164"`
	Email        *string  `json:"email" binding:"omitempty,email"`
	Website      *string  `json:"website" binding:"omitempty,url"`
	PricePerHour *float64 `json:"pricePerHour" binding:"omitempty,min=0"`
	IsActive     *bool    `json:"isActive"`
}

// VenueResponse represents the response payload for venue
type VenueResponse struct {
	ID           uint64   `json:"id"`
	OwnerID      uint64   `json:"ownerId"`
	Name         string   `json:"name"`
	Description  *string  `json:"description,omitempty"`
	Address      string   `json:"address"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Phone        string   `json:"phone"`
	Email        *string  `json:"email,omitempty"`
	Website      *string  `json:"website,omitempty"`
	PricePerHour float64  `json:"pricePerHour"`
	IsActive     bool     `json:"isActive"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
	Owner        *User    `json:"owner,omitempty"`
	Courts       []*Court `json:"courts,omitempty"`
}

// VenueListRequest represents the request payload for listing venues
type VenueListRequest struct {
	Latitude  *float64 `json:"latitude" form:"latitude"`
	Longitude *float64 `json:"longitude" form:"longitude"`
	Radius    *float64 `json:"radius" form:"radius"` // in kilometers
	MinPrice  *float64 `json:"minPrice" form:"minPrice"`
	MaxPrice  *float64 `json:"maxPrice" form:"maxPrice"`
	IsActive  *bool    `json:"isActive" form:"isActive"`
	Page      int      `json:"page" form:"page"`
	Limit     int      `json:"limit" form:"limit"`
}

// VenueListResponse represents the response payload for venue list
type VenueListResponse struct {
	Venues     []*VenueResponse `json:"venues"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

// CreateCourtRequest represents the request payload for creating a court
type CreateCourtRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	CourtType   string  `json:"courtType" binding:"required,oneof=indoor outdoor"`
	Surface     string  `json:"surface" binding:"required,oneof=hard clay grass synthetic"`
	IsActive    bool    `json:"isActive"`
}

// UpdateCourtRequest represents the request payload for updating a court
type UpdateCourtRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	CourtType   *string `json:"courtType" binding:"omitempty,oneof=indoor outdoor"`
	Surface     *string `json:"surface" binding:"omitempty,oneof=hard clay grass synthetic"`
	IsActive    *bool   `json:"isActive"`
}

// CourtResponse represents the response payload for court
type CourtResponse struct {
	ID          uint64              `json:"id"`
	VenueID     uint64              `json:"venueId"`
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	CourtType   string              `json:"courtType"`
	Surface     string              `json:"surface"`
	IsActive    bool                `json:"isActive"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
	Venue       *Venue              `json:"venue,omitempty"`
	Slots       []*AvailabilitySlot `json:"slots,omitempty"`
}

// CreateSlotRequest represents the request payload for creating availability slots
type CreateSlotRequest struct {
	StartTime   time.Time `json:"startTime" binding:"required"`
	EndTime     time.Time `json:"endTime" binding:"required"`
	IsAvailable bool      `json:"isAvailable"`
}

// UpdateSlotRequest represents the request payload for updating availability slots
type UpdateSlotRequest struct {
	StartTime   *time.Time `json:"startTime" binding:"omitempty"`
	EndTime     *time.Time `json:"endTime" binding:"omitempty"`
	IsAvailable *bool      `json:"isAvailable"`
}

// SlotResponse represents the response payload for availability slot
type SlotResponse struct {
	ID          uint64 `json:"id"`
	CourtID     uint64 `json:"courtId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Court       *Court `json:"court,omitempty"`
}

// SlotListRequest represents the request payload for listing slots
type SlotListRequest struct {
	StartDate   *time.Time `json:"startDate" form:"startDate"`
	EndDate     *time.Time `json:"endDate" form:"endDate"`
	IsAvailable *bool      `json:"isAvailable" form:"isAvailable"`
	Page        int        `json:"page" form:"page"`
	Limit       int        `json:"limit" form:"limit"`
}

// SlotListResponse represents the response payload for slot list
type SlotListResponse struct {
	Slots      []*SlotResponse `json:"slots"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"totalPages"`
}
