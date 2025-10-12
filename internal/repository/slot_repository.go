package repository

import (
	"matchee/services/internal/entity"
	"time"

	"gorm.io/gorm"
)

// SlotRepository defines the interface for availability slot operations
type SlotRepository interface {
	CreateSlot(slot *entity.AvailabilitySlot) error
	GetSlotByID(id uint64) (*entity.AvailabilitySlot, error)
	UpdateSlot(id uint64, slot *entity.AvailabilitySlot) error
	DeleteSlot(id uint64) error
	GetSlotsByCourt(courtID uint64, startDate, endDate *time.Time, isActive *bool, page, limit int) ([]*entity.AvailabilitySlot, int64, error)
	GetAllSlots(page, limit int) ([]*entity.AvailabilitySlot, int64, error)
}

type slotRepository struct {
	db *gorm.DB
}

// NewSlotRepository creates a new instance of SlotRepository
func NewSlotRepository(db *gorm.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) CreateSlot(slot *entity.AvailabilitySlot) error {
	return r.db.Create(slot).Error
}

func (r *slotRepository) GetSlotByID(id uint64) (*entity.AvailabilitySlot, error) {
	var slot entity.AvailabilitySlot
	err := r.db.Preload("Court").Where("id = ?", id).First(&slot).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *slotRepository) UpdateSlot(id uint64, slot *entity.AvailabilitySlot) error {
	return r.db.Model(&entity.AvailabilitySlot{}).Where("id = ?", id).Updates(slot).Error
}

func (r *slotRepository) DeleteSlot(id uint64) error {
	return r.db.Delete(&entity.AvailabilitySlot{}, id).Error
}

func (r *slotRepository) GetSlotsByCourt(courtID uint64, startDate, endDate *time.Time, isActive *bool, page, limit int) ([]*entity.AvailabilitySlot, int64, error) {
	var slots []*entity.AvailabilitySlot
	var total int64

	query := r.db.Model(&entity.AvailabilitySlot{}).Where("court_id = ?", courtID)

	// Filter by date range
	if startDate != nil {
		query = query.Where("start_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("end_time <= ?", *endDate)
	}

	// Filter by active status
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Court").
		Offset(offset).Limit(limit).
		Order("start_time ASC").
		Find(&slots).Error

	return slots, total, err
}

func (r *slotRepository) GetAllSlots(page, limit int) ([]*entity.AvailabilitySlot, int64, error) {
	var slots []*entity.AvailabilitySlot
	var total int64

	query := r.db.Model(&entity.AvailabilitySlot{})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Court").
		Offset(offset).Limit(limit).
		Order("start_time ASC").
		Find(&slots).Error

	return slots, total, err
}
