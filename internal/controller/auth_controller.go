package controller

import (
	"net/http"

	"matchee/services/internal/entity"
	"matchee/services/internal/middleware"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUC usecase.AuthUsecase
}

func NewAuthController(authUC usecase.AuthUsecase) *AuthController {
	return &AuthController{
		authUC: authUC,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with default player role
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Registration data"
// @Success 201 {object} entity.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := ac.authUC.Register(&req)
	if err != nil {
		if err.Error() == "phone number already exists" || err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, entity.ConflictResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("User registered successfully", response))
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Login credentials"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := ac.authUC.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Login successful", response))
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req entity.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := ac.authUC.RefreshToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Token refreshed successfully", response))
}

// Logout handles user logout
// @Summary Logout user
// @Description Revoke refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RefreshTokenRequest true "Refresh token to revoke"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	var req entity.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	if err := ac.authUC.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Logged out successfully", nil))
}

// GetCurrentUser handles getting current user info
// @Summary Get current user
// @Description Get information of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/me [get]
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	user, err := ac.authUC.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	c.JSON(http.StatusOK, entity.OKResponse("User information retrieved successfully", user))
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} entity.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/users/me [put]
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	var req entity.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	user, err := ac.authUC.UpdateProfile(userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else if err.Error() == "phone number already exists" || err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, entity.ConflictResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	c.JSON(http.StatusOK, entity.OKResponse("Profile updated successfully", user))
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change current user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/change-password [post]
func (ac *AuthController) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	var req entity.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	if err := ac.authUC.ChangePassword(userID, &req); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else if err.Error() == "current password is incorrect" {
			c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Password changed successfully", nil))
}
