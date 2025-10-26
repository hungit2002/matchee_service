package entity

import "time"

// CreateMatchGroupRequest represents request to create a match group
type CreateMatchGroupRequest struct {
	MatchPostID   *uint64    `json:"matchPostId" binding:"omitempty"`
	VenueID       *uint64    `json:"venueId" binding:"omitempty"`
	CourtID       *uint64    `json:"courtId" binding:"omitempty"`
	ScheduledTime *time.Time `json:"scheduledTime" binding:"omitempty"`
	TotalPrice    *float64   `json:"totalPrice" binding:"omitempty,min=0"`
}

// UpdateMatchGroupRequest represents request to update a match group
// status: pending | confirmed | completed | cancelled
type UpdateMatchGroupRequest struct {
	VenueID       *uint64    `json:"venueId" binding:"omitempty"`
	CourtID       *uint64    `json:"courtId" binding:"omitempty"`
	ScheduledTime *time.Time `json:"scheduledTime" binding:"omitempty"`
	TotalPrice    *float64   `json:"totalPrice" binding:"omitempty,min=0"`
	Status        *string    `json:"status" binding:"omitempty,oneof=pending confirmed completed cancelled"`
}

// MatchGroupResponse represents response for a match group
type MatchGroupResponse struct {
	ID            uint64              `json:"id"`
	MatchPostID   *uint64             `json:"matchPostId,omitempty"`
	VenueID       *uint64             `json:"venueId,omitempty"`
	CourtID       *uint64             `json:"courtId,omitempty"`
	ScheduledTime *string             `json:"scheduledTime,omitempty"`
	TotalPrice    *float64            `json:"totalPrice,omitempty"`
	Status        string              `json:"status"`
	CreatedAt     string              `json:"createdAt"`
	UpdatedAt     string              `json:"updatedAt"`
	MatchPost     *MatchPost          `json:"matchPost,omitempty"`
	Venue         *Venue              `json:"venue,omitempty"`
	Court         *Court              `json:"court,omitempty"`
	Members       []*MatchGroupMember `json:"members,omitempty"`
}

// MatchGroupListResponse represents paginated list of groups
type MatchGroupListResponse struct {
	Groups     []*MatchGroupResponse `json:"groups"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"totalPages"`
}

// AddMemberRequest represents request to add a member to a group
type AddMemberRequest struct {
	UserID   uint64 `json:"userId" binding:"required"`
	IsLeader *bool  `json:"isLeader"`
}

// RemoveMemberRequest (from path params)
// empty - kept for consistency

type GroupListQuery struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}
