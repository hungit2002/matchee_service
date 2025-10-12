package controller

import (
	"net/http"
	"strconv"
	"time"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SlotController struct {
	slotUC usecase.SlotUsecase
}

func NewSlotController(slotUC usecase.SlotUsecase) *SlotController {
	return &SlotController{slotUC: slotUC}
}

// CreateSlot handles creating a new availability slot
// @Summary Create availability slot
// @Description Create a new availability slot for a court
// @Tags slots
// @Accept json
// @Produce json
// @Param id path int true "Court ID"
// @Param request body entity.CreateSlotRequest true "Slot data"
// @Success 201 {object} entity.SlotResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/courts/{id}/slots [post]
func (sc *SlotController) CreateSlot(c *gin.Context) {
	courtIDStr := c.Param("id")
	courtID, err := strconv.ParseUint(courtIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid court ID"))
		return
	}

	var req entity.CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := sc.slotUC.CreateSlot(courtID, &req)
	if err != nil {
		if err.Error() == "court not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("Availability slot created successfully", response))
}

// GetSlotsByCourt handles getting availability slots for a court
// @Summary Get slots by court
// @Description Get list of availability slots for a specific court
// @Tags slots
// @Produce json
// @Param id path int true "Court ID"
// @Param startDate query string false "Start date filter (RFC3339 format)"
// @Param endDate query string false "End date filter (RFC3339 format)"
// @Param isAvailable query boolean false "Filter by availability status"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {object} entity.SlotListResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/courts/{id}/slots [get]
func (sc *SlotController) GetSlotsByCourt(c *gin.Context) {
	courtIDStr := c.Param("id")
	courtID, err := strconv.ParseUint(courtIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid court ID"))
		return
	}

	// Parse query parameters
	req := &entity.SlotListRequest{}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			req.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			req.EndDate = &endDate
		}
	}

	if isAvailableStr := c.Query("isAvailable"); isAvailableStr != "" {
		if isAvailable, err := strconv.ParseBool(isAvailableStr); err == nil {
			req.IsAvailable = &isAvailable
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

	response, err := sc.slotUC.GetSlotsByCourt(courtID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Availability slots retrieved successfully", response))
}

// GetSlotByID handles getting slot details
// @Summary Get slot details
// @Description Get detailed slot information by ID
// @Tags slots
// @Produce json
// @Param id path int true "Slot ID"
// @Success 200 {object} entity.SlotResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/slots/{id} [get]
func (sc *SlotController) GetSlotByID(c *gin.Context) {
	slotIDStr := c.Param("id")
	slotID, err := strconv.ParseUint(slotIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid slot ID"))
		return
	}

	response, err := sc.slotUC.GetSlotByID(slotID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Availability slot retrieved successfully", response))
}

// UpdateSlot handles updating slot information
// @Summary Update slot
// @Description Update availability slot information
// @Tags slots
// @Accept json
// @Produce json
// @Param id path int true "Slot ID"
// @Param request body entity.UpdateSlotRequest true "Slot update data"
// @Success 200 {object} entity.SlotResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/slots/{id [put]
func (sc *SlotController) UpdateSlot(c *gin.Context) {
	slotIDStr := c.Param("id")
	slotID, err := strconv.ParseUint(slotIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid slot ID"))
		return
	}

	var req entity.UpdateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	response, err := sc.slotUC.UpdateSlot(slotID, &req)
	if err != nil {
		if err.Error() == "slot not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Availability slot updated successfully", response))
}

// DeleteSlot handles deleting a slot
// @Summary Delete slot
// @Description Delete an availability slot
// @Tags slots
// @Param id path int true "Slot ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/slots/{id} [delete]
func (sc *SlotController) DeleteSlot(c *gin.Context) {
	slotIDStr := c.Param("id")
	slotID, err := strconv.ParseUint(slotIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid slot ID"))
		return
	}

	err = sc.slotUC.DeleteSlot(slotID)
	if err != nil {
		if err.Error() == "slot not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Availability slot deleted successfully", nil))
}
