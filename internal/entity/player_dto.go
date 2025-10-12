package entity

// CreateUpdatePlayerProfileRequest represents the request payload for creating/updating player profile
type CreateUpdatePlayerProfileRequest struct {
	Level             string   `json:"level" binding:"required,oneof=beginner average good pro"`
	Gender            *string  `json:"gender" binding:"omitempty,oneof=male female other"`
	PreferredLocation *string  `json:"preferredLocation" binding:"omitempty,max=255"`
	Latitude          *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude         *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Bio               *string  `json:"bio" binding:"omitempty,max=1000"`
}

// PlayerProfileResponse represents the response payload for player profile
type PlayerProfileResponse struct {
	ID                uint64   `json:"id"`
	UserID            uint64   `json:"userId"`
	Level             string   `json:"level"`
	Gender            *string  `json:"gender,omitempty"`
	PreferredLocation *string  `json:"preferredLocation,omitempty"`
	Latitude          *float64 `json:"latitude,omitempty"`
	Longitude         *float64 `json:"longitude,omitempty"`
	Bio               *string  `json:"bio,omitempty"`
	Distance          *float64 `json:"distance,omitempty"` // Distance in kilometers
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
	User              *User    `json:"user,omitempty"`
}

// PlayerSuggestionRequest represents the request payload for player suggestions
type PlayerSuggestionRequest struct {
	Level     *string  `json:"level" binding:"omitempty,oneof=beginner average good pro"`
	Latitude  *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Radius    *float64 `json:"radius" binding:"omitempty,min=0.1,max=50"` // in kilometers
	Limit     *int     `json:"limit" binding:"omitempty,min=1,max=50"`
}

// PlayerSuggestionResponse represents the response payload for player suggestions
type PlayerSuggestionResponse struct {
	Players []PlayerProfileResponse `json:"players"`
	Total   int                     `json:"total"`
	Page    int                     `json:"page"`
	Limit   int                     `json:"limit"`
}
