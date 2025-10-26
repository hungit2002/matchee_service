package controller

import (
	"net/http"
	"strconv"

	"matchee/services/internal/entity"
	"matchee/services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MatchGroupController struct {
	groupUC usecase.MatchGroupUsecase
}

func NewMatchGroupController(groupUC usecase.MatchGroupUsecase) *MatchGroupController {
	return &MatchGroupController{groupUC: groupUC}
}

// CreateMatchGroup
// @Summary Create match group
// @Description Create a match group (system can create automatically)
// @Tags match-groups
// @Accept json
// @Produce json
// @Param request body entity.CreateMatchGroupRequest true "Group data"
// @Success 201 {object} entity.MatchGroupResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/match-groups [post]
func (gc *MatchGroupController) CreateMatchGroup(c *gin.Context) {
	var req entity.CreateMatchGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	resp, err := gc.groupUC.CreateGroup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, entity.CreatedResponse("Match group created successfully", resp))
}

// GetMyMatchGroups
// @Summary Get my match groups
// @Description Get match groups that the current user is a member of
// @Tags match-groups
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} entity.MatchGroupListResponse
// @Router /api/v1/match-groups [get]
func (gc *MatchGroupController) GetMyMatchGroups(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, entity.UnauthorizedResponse("User not authenticated"))
		return
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}
	page, limit := 1, 20
	if pageStr := c.Query("page"); pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil {
			page = v
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}
	resp, err := gc.groupUC.GetGroupsByUser(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Match groups retrieved successfully", resp))
}

// AddMember
// @Summary Add member to group
// @Description Add a user to a match group
// @Tags match-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body entity.AddMemberRequest true "Member data"
// @Success 200 {object} entity.MatchGroupResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/match-groups/{id}/members [post]
func (gc *MatchGroupController) AddMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid group ID"))
		return
	}
	var req entity.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	resp, err := gc.groupUC.AddMember(groupID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Member added to match group successfully", resp))
}

// RemoveMember
// @Summary Remove member from group
// @Description A member leaves a match group
// @Tags match-groups
// @Param id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/match-groups/{id}/members/{user_id} [delete]
func (gc *MatchGroupController) RemoveMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid group ID"))
		return
	}
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid user ID"))
		return
	}
	if err := gc.groupUC.RemoveMember(groupID, userID); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Member removed from match group successfully", nil))
}

// UpdateMatchGroup
// @Summary Update match group
// @Description Confirm / Cancel match group, or update schedule/venue/court/price
// @Tags match-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body entity.UpdateMatchGroupRequest true "Update data"
// @Success 200 {object} entity.MatchGroupResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/match-groups/{id} [put]
func (gc *MatchGroupController) UpdateMatchGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse("Invalid group ID"))
		return
	}
	var req entity.UpdateMatchGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	resp, err := gc.groupUC.UpdateGroup(groupID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.BadRequestResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, entity.OKResponse("Match group updated successfully", resp))
}
