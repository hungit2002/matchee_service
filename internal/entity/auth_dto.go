package entity

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	FullName string `json:"fullName" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,len=10"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,len=10"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the request payload for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// AuthResponse represents the response payload for authentication
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

// ChangePasswordRequest represents the request payload for password change
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}

// UpdateProfileRequest represents the request payload for profile update
type UpdateProfileRequest struct {
	FullName string  `json:"fullName" binding:"omitempty,min=2,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    string  `json:"phone" binding:"omitempty,len=10"`
}
