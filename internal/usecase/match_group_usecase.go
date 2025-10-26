package usecase

import (
	"context"
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"time"
)

// MatchGroupUsecase defines business logic for match groups
type MatchGroupUsecase interface {
	CreateGroup(req *entity.CreateMatchGroupRequest) (*entity.MatchGroupResponse, error)
	GetGroupsByUser(userID uint64, page, limit int) (*entity.MatchGroupListResponse, error)
	AddMember(groupID uint64, req *entity.AddMemberRequest) (*entity.MatchGroupResponse, error)
	RemoveMember(groupID uint64, userID uint64) error
	UpdateGroup(groupID uint64, req *entity.UpdateMatchGroupRequest) (*entity.MatchGroupResponse, error)
}

type matchGroupUsecase struct {
	groupRepo  repository.MatchGroupRepository
	memberRepo repository.MatchGroupMemberRepository
	userRepo   repository.UserRepository
}

func NewMatchGroupUsecase(groupRepo repository.MatchGroupRepository, memberRepo repository.MatchGroupMemberRepository, userRepo repository.UserRepository) MatchGroupUsecase {
	return &matchGroupUsecase{groupRepo: groupRepo, memberRepo: memberRepo, userRepo: userRepo}
}

func (u *matchGroupUsecase) CreateGroup(req *entity.CreateMatchGroupRequest) (*entity.MatchGroupResponse, error) {
	group := &entity.MatchGroup{
		MatchPostID:   req.MatchPostID,
		VenueID:       req.VenueID,
		CourtID:       req.CourtID,
		ScheduledTime: req.ScheduledTime,
		TotalPrice:    req.TotalPrice,
		Status:        "pending",
	}
	if err := u.groupRepo.CreateGroup(group); err != nil {
		return nil, err
	}
	created, err := u.groupRepo.GetGroupByID(group.ID)
	if err != nil {
		return nil, err
	}
	return u.toGroupResponse(created), nil
}

func (u *matchGroupUsecase) GetGroupsByUser(userID uint64, page, limit int) (*entity.MatchGroupListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	groups, total, err := u.groupRepo.GetGroupsByUser(userID, page, limit)
	if err != nil {
		return nil, err
	}
	resp := make([]*entity.MatchGroupResponse, len(groups))
	for i, g := range groups {
		resp[i] = u.toGroupResponse(g)
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &entity.MatchGroupListResponse{Groups: resp, Total: total, Page: page, Limit: limit, TotalPages: totalPages}, nil
}

func (u *matchGroupUsecase) AddMember(groupID uint64, req *entity.AddMemberRequest) (*entity.MatchGroupResponse, error) {
	// ensure group exists
	_, err := u.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, errors.New("group not found")
	}
	// ensure user exists
	if _, err := u.userRepo.GetByID(context.Background(), req.UserID); err != nil {
		return nil, errors.New("user not found")
	}
	// ensure not duplicated
	isMember, err := u.memberRepo.IsMember(groupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, errors.New("user already in group")
	}
	member := &entity.MatchGroupMember{MatchGroupID: groupID, UserID: req.UserID}
	if req.IsLeader != nil {
		member.IsLeader = *req.IsLeader
	}
	if err := u.memberRepo.AddMember(member); err != nil {
		return nil, err
	}
	updated, err := u.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	return u.toGroupResponse(updated), nil
}

func (u *matchGroupUsecase) RemoveMember(groupID uint64, userID uint64) error {
	// ensure member exists
	isMember, err := u.memberRepo.IsMember(groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("user not in group")
	}
	return u.memberRepo.RemoveMember(groupID, userID)
}

func (u *matchGroupUsecase) UpdateGroup(groupID uint64, req *entity.UpdateMatchGroupRequest) (*entity.MatchGroupResponse, error) {
	// ensure group exists
	if _, err := u.groupRepo.GetGroupByID(groupID); err != nil {
		return nil, errors.New("group not found")
	}
	update := &entity.MatchGroup{}
	if req.VenueID != nil {
		update.VenueID = req.VenueID
	}
	if req.CourtID != nil {
		update.CourtID = req.CourtID
	}
	if req.ScheduledTime != nil {
		update.ScheduledTime = req.ScheduledTime
	}
	if req.TotalPrice != nil {
		update.TotalPrice = req.TotalPrice
	}
	if req.Status != nil {
		update.Status = *req.Status
	}
	if err := u.groupRepo.UpdateGroup(groupID, update); err != nil {
		return nil, err
	}
	g, err := u.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	return u.toGroupResponse(g), nil
}

func (u *matchGroupUsecase) toGroupResponse(g *entity.MatchGroup) *entity.MatchGroupResponse {
	resp := &entity.MatchGroupResponse{
		ID:          g.ID,
		MatchPostID: g.MatchPostID,
		VenueID:     g.VenueID,
		CourtID:     g.CourtID,
		TotalPrice:  g.TotalPrice,
		Status:      g.Status,
		CreatedAt:   g.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   g.UpdatedAt.Format(time.RFC3339),
	}
	if g.ScheduledTime != nil {
		t := g.ScheduledTime.Format(time.RFC3339)
		resp.ScheduledTime = &t
	}
	if g.MatchPost != nil {
		resp.MatchPost = g.MatchPost
	}
	if g.Venue != nil {
		resp.Venue = g.Venue
	}
	if g.Court != nil {
		resp.Court = g.Court
	}
	if len(g.Members) > 0 {
		members := make([]*entity.MatchGroupMember, len(g.Members))
		for i, m := range g.Members {
			members[i] = &m
		}
		resp.Members = members
	}
	return resp
}
