package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

// MatchGroupRepository defines DB operations for match groups
type MatchGroupRepository interface {
	CreateGroup(group *entity.MatchGroup) error
	UpdateGroup(groupID uint64, update *entity.MatchGroup) error
	GetGroupByID(groupID uint64) (*entity.MatchGroup, error)
	GetGroupsByUser(userID uint64, page, limit int) ([]*entity.MatchGroup, int64, error)
}

// MatchGroupMemberRepository defines DB operations for group members
type MatchGroupMemberRepository interface {
	AddMember(member *entity.MatchGroupMember) error
	RemoveMember(groupID uint64, userID uint64) error
	CountMembers(groupID uint64) (int64, error)
	IsMember(groupID uint64, userID uint64) (bool, error)
}

type matchGroupRepository struct{ db *gorm.DB }

type matchGroupMemberRepository struct{ db *gorm.DB }

func NewMatchGroupRepository(db *gorm.DB) MatchGroupRepository {
	return &matchGroupRepository{db: db}
}

func NewMatchGroupMemberRepository(db *gorm.DB) MatchGroupMemberRepository {
	return &matchGroupMemberRepository{db: db}
}

func (r *matchGroupRepository) CreateGroup(group *entity.MatchGroup) error {
	return r.db.Create(group).Error
}

func (r *matchGroupRepository) UpdateGroup(groupID uint64, update *entity.MatchGroup) error {
	return r.db.Model(&entity.MatchGroup{}).Where("id = ?", groupID).Updates(update).Error
}

func (r *matchGroupRepository) GetGroupByID(groupID uint64) (*entity.MatchGroup, error) {
	var group entity.MatchGroup
	if err := r.db.Preload("MatchPost").Preload("Venue").Preload("Court").Preload("Members").Where("id = ?", groupID).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *matchGroupRepository) GetGroupsByUser(userID uint64, page, limit int) ([]*entity.MatchGroup, int64, error) {
	var groups []*entity.MatchGroup
	var total int64

	query := r.db.Model(&entity.MatchGroup{}).Joins("JOIN match_group_members mgm ON mgm.match_group_id = match_groups.id").Where("mgm.user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Preload("MatchPost").Preload("Venue").Preload("Court").Preload("Members").
		Offset(offset).Limit(limit).Order("match_groups.created_at DESC").Find(&groups).Error
	return groups, total, err
}

func (r *matchGroupMemberRepository) AddMember(member *entity.MatchGroupMember) error {
	return r.db.Create(member).Error
}

func (r *matchGroupMemberRepository) RemoveMember(groupID uint64, userID uint64) error {
	return r.db.Where("match_group_id = ? AND user_id = ?", groupID, userID).Delete(&entity.MatchGroupMember{}).Error
}

func (r *matchGroupMemberRepository) CountMembers(groupID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&entity.MatchGroupMember{}).Where("match_group_id = ?", groupID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *matchGroupMemberRepository) IsMember(groupID uint64, userID uint64) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.MatchGroupMember{}).Where("match_group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
