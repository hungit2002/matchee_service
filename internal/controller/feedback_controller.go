package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FeedbackController struct {
	feedbackUC usecase.FeedbackUsecase
}

func NewFeedbackController(feedbackUC usecase.FeedbackUsecase) *FeedbackController {
	return &FeedbackController{feedbackUC: feedbackUC}
}

// CreateFeedback
// @Summary Create feedback
// @Description Create feedback for a user after completing a match
// @Tags feedbacks
// @Accept json
// @Produce json
// @Param request body entity.CreateFeedbackRequest true "Feedback data"
// @Success 201 {object} entity.FeedbackResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/feedbacks [post]
func (fc *FeedbackController) CreateFeedback(c *gin.Context) {
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

	var req entity.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	resp, err := fc.feedbackUC.CreateFeedback(userIDUint, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("Feedback created successfully", resp))
}

// GetUserFeedbacks
// @Summary Get user feedback summary
// @Description Get aggregated feedback data for a specific user
// @Tags feedbacks
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.UserFeedbackSummary
// @Failure 404 {object} map[string]string
// @Router /api/v1/feedbacks/user/{id} [get]
func (fc *FeedbackController) GetUserFeedbacks(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}

	resp, err := fc.feedbackUC.GetUserFeedbacks(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("User feedbacks retrieved successfully", resp))
}

// GetBookingFeedbacks
// @Summary Get booking feedbacks
// @Description Get all feedbacks for a specific booking
// @Tags feedbacks
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} entity.BookingFeedbackList
// @Failure 404 {object} map[string]string
// @Router /api/v1/feedbacks/booking/{id} [get]
func (fc *FeedbackController) GetBookingFeedbacks(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid booking ID"))
		return
	}

	resp, err := fc.feedbackUC.GetBookingFeedbacks(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Booking feedbacks retrieved successfully", resp))
}
