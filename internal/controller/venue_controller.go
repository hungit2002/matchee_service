package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/middleware"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type VenueController struct {
	venueUC usecase.VenueUsecase
}

func NewVenueController(venueUC usecase.VenueUsecase) *VenueController {
	return &VenueController{venueUC: venueUC}
}

// CreateVenue handles creating a new venue
// @Summary Create venue
// @Description Create a new venue (for venue owners)
// @Tags venues
// @Accept json
// @Produce json
// @Param request body entity.CreateVenueRequest true "Venue data"
// @Success 201 {object} entity.VenueResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/venues [post]
func (vc *VenueController) CreateVenue(c *gin.Context) {
	ownerID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	var req entity.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := vc.venueUC.CreateVenue(ownerID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("Venue created successfully", response))
}

// UpdateVenue handles updating venue information
// @Summary Update venue
// @Description Update venue information
// @Tags venues
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param request body entity.UpdateVenueRequest true "Venue update data"
// @Success 200 {object} entity.VenueResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/venues/{id} [put]
func (vc *VenueController) UpdateVenue(c *gin.Context) {
	venueIDStr := c.Param("id")
	venueID, err := strconv.ParseUint(venueIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid venue ID"))
		return
	}

	var req entity.UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := vc.venueUC.UpdateVenue(venueID, &req)
	if err != nil {
		if err.Error() == "venue not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Venue updated successfully", response))
}

// GetVenueByID handles getting venue details
// @Summary Get venue details
// @Description Get detailed venue information by ID
// @Tags venues
// @Produce json
// @Param id path int true "Venue ID"
// @Success 200 {object} entity.VenueResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/venues/{id} [get]
func (vc *VenueController) GetVenueByID(c *gin.Context) {
	venueIDStr := c.Param("id")
	venueID, err := strconv.ParseUint(venueIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid venue ID"))
		return
	}

	response, err := vc.venueUC.GetVenueByID(venueID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Venue retrieved successfully", response))
}

// GetVenues handles getting venue list with filters
// @Summary Get venues
// @Description Get list of venues with optional filters (location, price, etc.)
// @Tags venues
// @Produce json
// @Param latitude query number false "Latitude for location filter"
// @Param longitude query number false "Longitude for location filter"
// @Param radius query number false "Radius in kilometers for location filter"
// @Param minPrice query number false "Minimum price per hour"
// @Param maxPrice query number false "Maximum price per hour"
// @Param isActive query boolean false "Filter by active status"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} entity.VenueListResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/venues [get]
func (vc *VenueController) GetVenues(c *gin.Context) {
	// Parse query parameters
	req := &entity.VenueListRequest{}

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

	if minPriceStr := c.Query("minPrice"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			req.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("maxPrice"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			req.MaxPrice = &maxPrice
		}
	}

	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
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

	response, err := vc.venueUC.GetVenuesByLocation(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Venues retrieved successfully", response))
}

// GetMyVenues handles getting venues owned by current user
// @Summary Get my venues
// @Description Get venues owned by current user
// @Tags venues
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} entity.VenueListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/venues/my [get]
func (vc *VenueController) GetMyVenues(c *gin.Context) {
	ownerID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}

	// Parse query parameters
	page := 1
	limit := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	response, err := vc.venueUC.GetVenuesByOwner(ownerID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("My venues retrieved successfully", response))
}

// DeleteVenue handles deleting a venue
// @Summary Delete venue
// @Description Delete a venue
// @Tags venues
// @Param id path int true "Venue ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/venues/{id} [delete]
func (vc *VenueController) DeleteVenue(c *gin.Context) {
	venueIDStr := c.Param("id")
	venueID, err := strconv.ParseUint(venueIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid venue ID"))
		return
	}

	err = vc.venueUC.DeleteVenue(venueID)
	if err != nil {
		if err.Error() == "venue not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Venue deleted successfully", nil))
}
