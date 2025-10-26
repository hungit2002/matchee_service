package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminUC usecase.AdminUsecase
}

func NewAdminController(adminUC usecase.AdminUsecase) *AdminController {
	return &AdminController{adminUC: adminUC}
}

// GetRevenueStatistics
// @Summary Get revenue statistics
// @Description Get revenue statistics for admin dashboard
// @Tags admin
// @Produce json
// @Param fromDate query string false "Start date (YYYY-MM-DD)"
// @Param toDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} entity.RevenueStatisticsResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/statistics/revenue [get]
func (ac *AdminController) GetRevenueStatistics(c *gin.Context) {
	fromDate := c.Query("fromDate")
	toDate := c.Query("toDate")

	resp, err := ac.adminUC.GetRevenueStatistics(fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Revenue statistics retrieved successfully", resp))
}

// GetUsers
// @Summary Get users list
// @Description Get paginated list of users for admin management
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Filter by status (active, locked)"
// @Param role query string false "Filter by role (user, venue_owner, admin)"
// @Param search query string false "Search by email or name"
// @Success 200 {object} entity.AdminUserListResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/users [get]
func (ac *AdminController) GetUsers(c *gin.Context) {
	req := &entity.AdminUserListRequest{}

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
	if status := c.Query("status"); status != "" {
		req.Status = &status
	}
	if role := c.Query("role"); role != "" {
		req.Role = &role
	}
	if search := c.Query("search"); search != "" {
		req.Search = &search
	}

	resp, err := ac.adminUC.GetUsers(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Users retrieved successfully", resp))
}

// LockUser
// @Summary Lock or unlock user
// @Description Lock or unlock a user account
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body entity.LockUserRequest true "Lock request data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/users/{id}/lock [put]
func (ac *AdminController) LockUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}

	var req entity.LockUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	err = ac.adminUC.LockUser(userID, &req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	action := "unlocked"
	if req.Locked {
		action = "locked"
	}

	c.JSON(http.StatusOK, entity.OKResponse("User "+action+" successfully", nil))
}

// ApproveVenue
// @Summary Approve or reject venue
// @Description Approve or reject a venue for operation
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Venue ID"
// @Param request body entity.VenueApprovalRequest true "Approval request data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/admin/venues/{id}/approve [put]
func (ac *AdminController) ApproveVenue(c *gin.Context) {
	venueIDStr := c.Param("id")
	venueID, err := strconv.ParseUint(venueIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid venue ID"))
		return
	}

	var req entity.VenueApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	err = ac.adminUC.ApproveVenue(venueID, &req)
	if err != nil {
		if err.Error() == "venue not found" {
			c.JSON(http.StatusNotFound, entity.NotFoundResponse(err.Error()))
		} else {
			c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		}
		return
	}

	action := "rejected"
	if req.Approved {
		action = "approved"
	}

	c.JSON(http.StatusOK, entity.OKResponse("Venue "+action+" successfully", nil))
}

// GetFeedbacks
// @Summary Get feedbacks list
// @Description Get paginated list of feedbacks for admin review
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param rating query int false "Filter by rating (1-5)"
// @Param fromDate query string false "Start date (YYYY-MM-DD)"
// @Param toDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} entity.AdminFeedbackListResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/feedbacks [get]
func (ac *AdminController) GetFeedbacks(c *gin.Context) {
	req := &entity.AdminFeedbackListRequest{}

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
	if ratingStr := c.Query("rating"); ratingStr != "" {
		if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
			req.Rating = &rating
		}
	}
	if fromDate := c.Query("fromDate"); fromDate != "" {
		req.FromDate = &fromDate
	}
	if toDate := c.Query("toDate"); toDate != "" {
		req.ToDate = &toDate
	}

	resp, err := ac.adminUC.GetFeedbacks(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.OKResponse("Feedbacks retrieved successfully", resp))
}
