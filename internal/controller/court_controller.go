package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CourtController struct {
	courtUC usecase.CourtUsecase
}

func NewCourtController(courtUC usecase.CourtUsecase) *CourtController {
	return &CourtController{courtUC: courtUC}
}

// CreateCourt handles creating a new court
// @Summary Create court
// @Description Create a new court for a venue
// @Tags courts
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param request body entity.CreateCourtRequest true "Court data"
// @Success 201 {object} entity.CourtResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/venues/{id}/courts [post]
func (cc *CourtController) CreateCourt(c *gin.Context) {
	venueIDStr := c.Param("id")
	venueID, err := strconv.ParseUint(venueIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid venue ID"))
		return
	}

	var req entity.CreateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := cc.courtUC.CreateCourt(venueID, &req)
	if err != nil {
		if err.Error() == "venue not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("Court created successfully", response))
}

// GetCourtsByVenue handles getting courts for a venue
// @Summary Get courts by venue
// @Description Get list of courts for a specific venue
// @Tags courts
// @Produce json
// @Param id path int true "Venue ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {array} entity.CourtResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/venues/{id}/courts [get]
func (cc *CourtController) GetCourtsByVenue(c *gin.Context) {
	venueIDStr := c.Param("id")
	venueID, err := strconv.ParseUint(venueIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid venue ID"))
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

	response, total, err := cc.courtUC.GetCourtsByVenue(venueID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Courts retrieved successfully", gin.H{
		"courts": response,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}))
}

// GetCourtByID handles getting court details
// @Summary Get court details
// @Description Get detailed court information by ID
// @Tags courts
// @Produce json
// @Param id path int true "Court ID"
// @Success 200 {object} entity.CourtResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/courts/{id} [get]
func (cc *CourtController) GetCourtByID(c *gin.Context) {
	courtIDStr := c.Param("id")
	courtID, err := strconv.ParseUint(courtIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid court ID"))
		return
	}

	response, err := cc.courtUC.GetCourtByID(courtID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Court retrieved successfully", response))
}

// UpdateCourt handles updating court information
// @Summary Update court
// @Description Update court information
// @Tags courts
// @Accept json
// @Produce json
// @Param id path int true "Court ID"
// @Param request body entity.UpdateCourtRequest true "Court update data"
// @Success 200 {object} entity.CourtResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/courts/{id} [put]
func (cc *CourtController) UpdateCourt(c *gin.Context) {
	courtIDStr := c.Param("id")
	courtID, err := strconv.ParseUint(courtIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid court ID"))
		return
	}

	var req entity.UpdateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := cc.courtUC.UpdateCourt(courtID, &req)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Court updated successfully", response))
}

// DeleteCourt handles deleting a court
// @Summary Delete court
// @Description Delete a court
// @Tags courts
// @Param id path int true "Court ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/courts/{id} [delete]
func (cc *CourtController) DeleteCourt(c *gin.Context) {
	courtIDStr := c.Param("id")
	courtID, err := strconv.ParseUint(courtIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid court ID"))
		return
	}

	err = cc.courtUC.DeleteCourt(courtID)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Court deleted successfully", nil))
}
