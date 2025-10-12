package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

// MatchPostRepository defines the interface for match post data operations
type MatchPostRepository interface {
	CreateMatchPost(matchPost *entity.MatchPost) error
	UpdateMatchPost(matchPostID uint64, matchPost *entity.MatchPost) error
	GetMatchPostByID(matchPostID uint64) (*entity.MatchPost, error)
	DeleteMatchPost(matchPostID uint64) error
	GetMatchPostsByUser(userID uint64, page, limit int) ([]*entity.MatchPost, int64, error)
	GetMatchPostsByLocation(latitude, longitude, radius float64, desiredLevel, status *string, page, limit int) ([]*entity.MatchPost, int64, error)
	GetMatchPostsByVenue(venueID uint64, status *string, page, limit int) ([]*entity.MatchPost, int64, error)
	GetAllMatchPosts(status *string, page, limit int) ([]*entity.MatchPost, int64, error)
	GetSuggestedMatchPosts(userID uint64, latitude, longitude, radius *float64, desiredLevel *string, limit int) ([]*entity.MatchPost, int64, error)
}

type matchPostRepository struct {
	db *gorm.DB
}

// NewMatchPostRepository creates a new instance of MatchPostRepository
func NewMatchPostRepository(db *gorm.DB) MatchPostRepository {
	return &matchPostRepository{db: db}
}

func (r *matchPostRepository) CreateMatchPost(matchPost *entity.MatchPost) error {
	return r.db.Create(matchPost).Error
}

func (r *matchPostRepository) UpdateMatchPost(matchPostID uint64, matchPost *entity.MatchPost) error {
	return r.db.Model(&entity.MatchPost{}).Where("id = ?", matchPostID).Updates(matchPost).Error
}

func (r *matchPostRepository) GetMatchPostByID(matchPostID uint64) (*entity.MatchPost, error) {
	var matchPost entity.MatchPost
	err := r.db.Preload("User").Preload("Venue").Preload("MatchGroups").Where("id = ?", matchPostID).First(&matchPost).Error
	if err != nil {
		return nil, err
	}
	return &matchPost, nil
}

func (r *matchPostRepository) DeleteMatchPost(matchPostID uint64) error {
	return r.db.Delete(&entity.MatchPost{}, matchPostID).Error
}

func (r *matchPostRepository) GetMatchPostsByUser(userID uint64, page, limit int) ([]*entity.MatchPost, int64, error) {
	var matchPosts []*entity.MatchPost
	var total int64

	query := r.db.Model(&entity.MatchPost{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("User").Preload("Venue").Preload("MatchGroups").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&matchPosts).Error

	return matchPosts, total, err
}

func (r *matchPostRepository) GetMatchPostsByLocation(latitude, longitude, radius float64, desiredLevel, status *string, page, limit int) ([]*entity.MatchPost, int64, error) {
	var matchPosts []*entity.MatchPost
	var total int64

	query := r.db.Model(&entity.MatchPost{}).Where("latitude IS NOT NULL AND longitude IS NOT NULL")

	// Add location filter using Haversine formula
	query = query.Where("(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) <= ?",
		latitude, longitude, latitude, radius)

	// Add filters
	if desiredLevel != nil {
		query = query.Where("desired_level = ?", *desiredLevel)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("User").Preload("Venue").Preload("MatchGroups").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&matchPosts).Error

	return matchPosts, total, err
}

func (r *matchPostRepository) GetMatchPostsByVenue(venueID uint64, status *string, page, limit int) ([]*entity.MatchPost, int64, error) {
	var matchPosts []*entity.MatchPost
	var total int64

	query := r.db.Model(&entity.MatchPost{}).Where("venue_id = ?", venueID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("User").Preload("Venue").Preload("MatchGroups").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&matchPosts).Error

	return matchPosts, total, err
}

func (r *matchPostRepository) GetAllMatchPosts(status *string, page, limit int) ([]*entity.MatchPost, int64, error) {
	var matchPosts []*entity.MatchPost
	var total int64

	query := r.db.Model(&entity.MatchPost{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("User").Preload("Venue").Preload("MatchGroups").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&matchPosts).Error

	return matchPosts, total, err
}

func (r *matchPostRepository) GetSuggestedMatchPosts(userID uint64, latitude, longitude, radius *float64, desiredLevel *string, limit int) ([]*entity.MatchPost, int64, error) {
	var matchPosts []*entity.MatchPost
	var total int64

	query := r.db.Model(&entity.MatchPost{}).Where("user_id != ? AND status = 'open'", userID)

	// Add location filter if provided
	if latitude != nil && longitude != nil && radius != nil {
		query = query.Where("latitude IS NOT NULL AND longitude IS NOT NULL")
		query = query.Where("(6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) <= ?",
			*latitude, *longitude, *latitude, *radius)
	}

	// Add desired level filter
	if desiredLevel != nil {
		query = query.Where("desired_level = ?", *desiredLevel)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get limited results
	err := query.Preload("User").Preload("Venue").Preload("MatchGroups").
		Limit(limit).Order("created_at DESC").Find(&matchPosts).Error

	return matchPosts, total, err
}
