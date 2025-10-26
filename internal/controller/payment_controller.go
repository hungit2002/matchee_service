package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentUC usecase.PaymentUsecase
}

func NewPaymentController(paymentUC usecase.PaymentUsecase) *PaymentController {
	return &PaymentController{paymentUC: paymentUC}
}

// CreatePayment
// @Summary Create payment request
// @Description Create a payment request for a booking
// @Tags payments
// @Accept json
// @Produce json
// @Param request body entity.CreatePaymentRequest true "Payment data"
// @Success 201 {object} entity.PaymentResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/payments [post]
func (pc *PaymentController) CreatePayment(c *gin.Context) {
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

	var req entity.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	resp, err := pc.paymentUC.CreatePayment(userIDUint, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, entity.CreatedResponse("Payment request created successfully", resp))
}

// GetPayments
// @Summary Get payment history
// @Description Get payment history for the current user
// @Tags payments
// @Produce json
// @Param status query string false "Payment status filter (pending, success, failed)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} entity.PaymentListResponse
// @Router /api/v1/payments [get]
func (pc *PaymentController) GetPayments(c *gin.Context) {
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

	req := &entity.PaymentListRequest{}
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

	resp, err := pc.paymentUC.GetPayments(userIDUint, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Payments retrieved successfully", resp))
}

// GetPaymentByID
// @Summary Get payment details
// @Description Get detailed payment information by ID
// @Tags payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} entity.PaymentResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/payments/{id} [get]
func (pc *PaymentController) GetPaymentByID(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid payment ID"))
		return
	}

	resp, err := pc.paymentUC.GetPaymentByID(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Payment retrieved successfully", resp))
}

// HandlePayPalWebhook
// @Summary PayPal webhook
// @Description Handle PayPal webhook notifications
// @Tags payments
// @Accept json
// @Produce json
// @Param request body entity.PayPalWebhookRequest true "PayPal webhook data"
// @Success 200 {object} map[string]string
// @Router /api/v1/payments/webhook [post]
func (pc *PaymentController) HandlePayPalWebhook(c *gin.Context) {
	var req entity.PayPalWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	if err := pc.paymentUC.HandlePayPalWebhook(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Webhook processed successfully", nil))
}
