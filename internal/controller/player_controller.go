package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/middleware"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PlayerController struct {
	playerUC usecase.PlayerUsecase
}

func NewPlayerController(playerUC usecase.PlayerUsecase) *PlayerController {
	return &PlayerController{
		playerUC: playerUC,
	}
}

// CreateOrUpdatePlayerProfile handles creating or updating player profile
// @Summary Create or update player profile
// @Description Create or update the current user's player profile with level, position, gender, etc.
// @Tags player
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.CreateUpdatePlayerProfileRequest true "Player profile data"
// @Success 200 {object} entity.PlayerProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/player/profile [post]
func (pc *PlayerController) CreateOrUpdatePlayerProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req entity.CreateUpdatePlayerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := pc.playerUC.CreateOrUpdatePlayerProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPlayerProfileByID handles getting player profile by ID
// @Summary Get player profile by ID
// @Description Get detailed player profile information by profile ID
// @Tags player
// @Accept json
// @Produce json
// @Param id path int true "Player Profile ID"
// @Success 200 {object} entity.PlayerProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/player/profile/{id} [get]
func (pc *PlayerController) GetPlayerProfileByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	response, err := pc.playerUC.GetPlayerProfileByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPlayerProfileByUserID handles getting current user's player profile
// @Summary Get current user's player profile
// @Description Get the current user's player profile information
// @Tags player
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.PlayerProfileResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/player/profile [get]
func (pc *PlayerController) GetPlayerProfileByUserID(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	response, err := pc.playerUC.GetPlayerProfileByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPlayerSuggestions handles getting player suggestions
// @Summary Get player suggestions
// @Description Get suggested players based on level, location, and other criteria
// @Tags player
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.PlayerSuggestionRequest true "Suggestion criteria"
// @Success 200 {object} entity.PlayerSuggestionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/player/suggestions [get]
func (pc *PlayerController) GetPlayerSuggestions(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req entity.PlayerSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := pc.playerUC.GetPlayerSuggestions(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPlayerSuggestionsQuery handles getting player suggestions via query parameters
// @Summary Get player suggestions (query params)
// @Description Get suggested players based on query parameters
// @Tags player
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param level query string false "Player level (beginner, average, good, pro)"
// @Param latitude query number false "Latitude for location-based search"
// @Param longitude query number false "Longitude for location-based search"
// @Param radius query number false "Search radius in kilometers (default: 10)"
// @Param limit query int false "Maximum number of suggestions (default: 20)"
// @Success 200 {object} entity.PlayerSuggestionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/player/suggestions [get]
func (pc *PlayerController) GetPlayerSuggestionsQuery(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	req := &entity.PlayerSuggestionRequest{}

	if level := c.Query("level"); level != "" {
		req.Level = &level
	}

	if latStr := c.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			req.Latitude = &lat
		}
	}

	if lonStr := c.Query("longitude"); lonStr != "" {
		if lon, err := strconv.ParseFloat(lonStr, 64); err == nil {
			req.Longitude = &lon
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			req.Radius = &radius
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = &limit
		}
	}

	response, err := pc.playerUC.GetPlayerSuggestions(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
