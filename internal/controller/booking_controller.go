package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingUC usecase.BookingUsecase
}

func NewBookingController(bookingUC usecase.BookingUsecase) *BookingController {
	return &BookingController{bookingUC: bookingUC}
}

// CreateBooking
// @Summary Create booking
// @Description Create a booking for a match group or venue/court
// @Tags bookings
// @Accept json
// @Produce json
// @Param request body entity.CreateBookingRequest true "Booking data"
// @Success 201 {object} entity.BookingResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/bookings [post]
func (bc *BookingController) CreateBooking(c *gin.Context) {
	var req entity.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	resp, err := bc.bookingUC.CreateBooking(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, entity.CreatedResponse("Booking created successfully", resp))
}

// GetBookings
// @Summary Get bookings
// @Description List bookings filtered by user or venue
// @Tags bookings
// @Produce json
// @Param userId query int false "User ID"
// @Param venueId query int false "Venue ID"
// @Param status query string false "Status filter (reserved, paid, completed, cancelled)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} entity.BookingListResponse
// @Router /api/v1/bookings [get]
func (bc *BookingController) GetBookings(c *gin.Context) {
	req := &entity.BookingListRequest{}
	if userIDStr := c.Query("userId"); userIDStr != "" {
		if v, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			req.UserID = &v
		}
	}
	if venueIDStr := c.Query("venueId"); venueIDStr != "" {
		if v, err := strconv.ParseUint(venueIDStr, 10, 64); err == nil {
			req.VenueID = &v
		}
	}
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if pageStr := c.Query("page"); pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil {
			req.Page = v
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = v
		}
	}
	resp, err := bc.bookingUC.GetBookings(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Bookings retrieved successfully", resp))
}

// GetBookingByID
// @Summary Get booking details
// @Description Get detailed booking information by ID
// @Tags bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} entity.BookingResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/bookings/{id} [get]
func (bc *BookingController) GetBookingByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid booking ID"))
		return
	}
	resp, err := bc.bookingUC.GetBookingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Booking retrieved successfully", resp))
}

// UpdateBookingStatus
// @Summary Update booking status
// @Description Update status: reserved -> paid -> completed, or cancelled
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param request body entity.UpdateBookingStatusRequest true "Status data"
// @Success 200 {object} entity.BookingResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/bookings/{id}/status [put]
func (bc *BookingController) UpdateBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid booking ID"))
		return
	}
	var req entity.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	resp, err := bc.bookingUC.UpdateBookingStatus(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Booking status updated successfully", resp))
}

// DeleteBooking
// @Summary Delete booking
// @Description Cancel a booking (set status to cancelled or delete)
// @Tags bookings
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/bookings/{id} [delete]
func (bc *BookingController) DeleteBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid booking ID"))
		return
	}
	if err := bc.bookingUC.DeleteBooking(id); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Booking deleted successfully", nil))
}
