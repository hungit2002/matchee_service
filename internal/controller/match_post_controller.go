package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MatchPostController struct {
	matchPostUC usecase.MatchPostUsecase
}

func NewMatchPostController(matchPostUC usecase.MatchPostUsecase) *MatchPostController {
	return &MatchPostController{matchPostUC: matchPostUC}
}

// CreateMatchPost handles creating a new match post
// @Summary Create match post
// @Description Create a new match post for finding opponents
// @Tags match-posts
// @Accept json
// @Produce json
// @Param request body entity.CreateMatchPostRequest true "Match post data"
// @Success 201 {object} entity.MatchPostResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/match-posts [post]
func (mpc *MatchPostController) CreateMatchPost(c *gin.Context) {
	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}

	var req entity.CreateMatchPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := mpc.matchPostUC.CreateMatchPost(userIDUint, &req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "venue not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("Match post created successfully", response))
}

// GetMatchPosts handles getting match posts with filters
// @Summary Get match posts
// @Description Get list of match posts with optional filters
// @Tags match-posts
// @Produce json
// @Param desiredLevel query string false "Filter by desired level (beginner, average, good, pro)"
// @Param location query string false "Filter by location"
// @Param latitude query number false "Latitude for location-based search"
// @Param longitude query number false "Longitude for location-based search"
// @Param radius query number false "Search radius in kilometers"
// @Param venueId query int false "Filter by venue ID"
// @Param status query string false "Filter by status (open, matched, cancelled, done)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} entity.MatchPostListResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/match-posts [get]
func (mpc *MatchPostController) GetMatchPosts(c *gin.Context) {
	// Parse query parameters
	req := &entity.MatchPostListRequest{}

	if desiredLevel := c.Query("desiredLevel"); desiredLevel != "" {
		req.DesiredLevel = &desiredLevel
	}

	if location := c.Query("location"); location != "" {
		req.Location = &location
	}

	if latitudeStr := c.Query("latitude"); latitudeStr != "" {
		if latitude, err := strconv.ParseFloat(latitudeStr, 64); err == nil {
			req.Latitude = &latitude
		}
	}

	if longitudeStr := c.Query("longitude"); longitudeStr != "" {
		if longitude, err := strconv.ParseFloat(longitudeStr, 64); err == nil {
			req.Longitude = &longitude
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			req.Radius = &radius
		}
	}

	if venueIDStr := c.Query("venueId"); venueIDStr != "" {
		if venueID, err := strconv.ParseUint(venueIDStr, 10, 64); err == nil {
			req.VenueID = &venueID
		}
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	response, err := mpc.matchPostUC.GetMatchPostsByLocation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Match posts retrieved successfully", response))
}

// GetMatchPostByID handles getting match post details
// @Summary Get match post details
// @Description Get detailed match post information by ID
// @Tags match-posts
// @Produce json
// @Param id path int true "Match post ID"
// @Success 200 {object} entity.MatchPostResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/match-posts/{id} [get]
func (mpc *MatchPostController) GetMatchPostByID(c *gin.Context) {
	matchPostIDStr := c.Param("id")
	matchPostID, err := strconv.ParseUint(matchPostIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid match post ID"))
		return
	}

	response, err := mpc.matchPostUC.GetMatchPostByID(matchPostID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Match post retrieved successfully", response))
}

// UpdateMatchPost handles updating match post information
// @Summary Update match post
// @Description Update match post information
// @Tags match-posts
// @Accept json
// @Produce json
// @Param id path int true "Match post ID"
// @Param request body entity.UpdateMatchPostRequest true "Match post update data"
// @Success 200 {object} entity.MatchPostResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/match-posts/{id} [put]
func (mpc *MatchPostController) UpdateMatchPost(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}

	matchPostIDStr := c.Param("id")
	matchPostID, err := strconv.ParseUint(matchPostIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid match post ID"))
		return
	}

	var req entity.UpdateMatchPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := mpc.matchPostUC.UpdateMatchPost(matchPostID, userIDUint, &req)
	if err != nil {
		if err.Error() == "match post not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else if err.Error() == "unauthorized to update this match post" {
			c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse(err.Error()))
		} else if err.Error() == "venue not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Match post updated successfully", response))
}

// DeleteMatchPost handles deleting a match post
// @Summary Delete match post
// @Description Delete a match post (set status to cancelled)
// @Tags match-posts
// @Param id path int true "Match post ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/match-posts/{id} [delete]
func (mpc *MatchPostController) DeleteMatchPost(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	userIDUint, ok := userID.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}

	matchPostIDStr := c.Param("id")
	matchPostID, err := strconv.ParseUint(matchPostIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid match post ID"))
		return
	}

	err = mpc.matchPostUC.DeleteMatchPost(matchPostID, userIDUint)
	if err != nil {
		if err.Error() == "match post not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else if err.Error() == "unauthorized to delete this match post" {
			c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Match post deleted successfully", nil))
}

// GetSuggestedMatchPosts handles getting suggested match posts
// @Summary Get suggested match posts
// @Description Get match post suggestions for a user based on their preferences
// @Tags match-posts
// @Produce json
// @Param userId query int true "User ID"
// @Param desiredLevel query string false "Desired level filter (beginner, average, good, pro)"
// @Param latitude query number false "User latitude"
// @Param longitude query number false "User longitude"
// @Param radius query number false "Search radius in kilometers"
// @Param limit query int false "Number of suggestions (default: 10, max: 50)"
// @Success 200 {object} entity.MatchPostSuggestResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/match-posts/suggest [get]
func (mpc *MatchPostController) GetSuggestedMatchPosts(c *gin.Context) {
	// Parse query parameters
	req := &entity.MatchPostSuggestRequest{}

	if userIDStr := c.Query("userId"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			req.UserID = userID
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("User ID is required"))
		return
	}

	if desiredLevel := c.Query("desiredLevel"); desiredLevel != "" {
		req.DesiredLevel = &desiredLevel
	}

	if latitudeStr := c.Query("latitude"); latitudeStr != "" {
		if latitude, err := strconv.ParseFloat(latitudeStr, 64); err == nil {
			req.Latitude = &latitude
		}
	}

	if longitudeStr := c.Query("longitude"); longitudeStr != "" {
		if longitude, err := strconv.ParseFloat(longitudeStr, 64); err == nil {
			req.Longitude = &longitude
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			req.Radius = &radius
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	response, err := mpc.matchPostUC.GetSuggestedMatchPosts(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Match post suggestions retrieved successfully", response))
}
