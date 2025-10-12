package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

// CourtRepository defines the interface for court operations
type CourtRepository interface {
	CreateCourt(court *entity.Court) error
	GetCourtByID(id uint64) (*entity.Court, error)
	UpdateCourt(id uint64, court *entity.Court) error
	DeleteCourt(id uint64) error
	GetCourtsByVenue(venueID uint64, page, limit int) ([]*entity.Court, int64, error)
	GetAllCourts(page, limit int) ([]*entity.Court, int64, error)
}

type courtRepository struct {
	db *gorm.DB
}

// NewCourtRepository creates a new instance of CourtRepository
func NewCourtRepository(db *gorm.DB) CourtRepository {
	return &courtRepository{db: db}
}

func (r *courtRepository) CreateCourt(court *entity.Court) error {
	return r.db.Create(court).Error
}

func (r *courtRepository) GetCourtByID(id uint64) (*entity.Court, error) {
	var court entity.Court
	err := r.db.Preload("Venue").Preload("AvailabilitySlots").Where("id = ?", id).First(&court).Error
	if err != nil {
		return nil, err
	}
	return &court, nil
}

func (r *courtRepository) UpdateCourt(id uint64, court *entity.Court) error {
	return r.db.Model(&entity.Court{}).Where("id = ?", id).Updates(court).Error
}

func (r *courtRepository) DeleteCourt(id uint64) error {
	return r.db.Delete(&entity.Court{}, id).Error
}

func (r *courtRepository) GetCourtsByVenue(venueID uint64, page, limit int) ([]*entity.Court, int64, error) {
	var courts []*entity.Court
	var total int64

	query := r.db.Model(&entity.Court{}).Where("venue_id = ?", venueID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Venue").Preload("AvailabilitySlots").
		Offset(offset).Limit(limit).
		Find(&courts).Error

	return courts, total, err
}

func (r *courtRepository) GetAllCourts(page, limit int) ([]*entity.Court, int64, error) {
	var courts []*entity.Court
	var total int64

	query := r.db.Model(&entity.Court{})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Venue").Preload("AvailabilitySlots").
		Offset(offset).Limit(limit).
		Find(&courts).Error

	return courts, total, err
}
