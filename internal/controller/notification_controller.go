package controller

import (
	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationController struct{ uc usecase.NotificationUsecase }

func NewNotificationController(uc usecase.NotificationUsecase) *NotificationController {
	return &NotificationController{uc: uc}
}

// GET /api/v1/notifications
func (ctl *NotificationController) GetMyNotifications(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("unauthorized"))
		return
	}
	userID, ok := uidVal.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("invalid userID"))
		return
	}
	onlyUnread := c.Query("onlyUnread") == "true"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	resp, err := ctl.uc.GetMyNotifications(userID, onlyUnread, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("OK", resp))
}

// PUT /api/v1/notifications/:id/read
func (ctl *NotificationController) MarkRead(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("unauthorized"))
		return
	}
	userID, ok := uidVal.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("invalid userID"))
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("invalid id"))
		return
	}
	if err := ctl.uc.MarkRead(userID, id); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Marked as read", nil))
}

// POST /api/v1/notifications
func (ctl *NotificationController) Create(c *gin.Context) {
	var req entity.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	resp, err := ctl.uc.CreateNotification(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, entity.CreatedResponse("Created", resp))
}
