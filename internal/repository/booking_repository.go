package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(b *entity.Booking) error
	UpdateBooking(bookingID uint64, update *entity.Booking) error
	GetBookingByID(bookingID uint64) (*entity.Booking, error)
	GetBookings(userID, venueID *uint64, status *string, page, limit int) ([]*entity.Booking, int64, error)
	DeleteBooking(bookingID uint64) error
}

type bookingRepository struct{ db *gorm.DB }

func NewBookingRepository(db *gorm.DB) BookingRepository { return &bookingRepository{db: db} }

func (r *bookingRepository) CreateBooking(b *entity.Booking) error { return r.db.Create(b).Error }

func (r *bookingRepository) UpdateBooking(bookingID uint64, update *entity.Booking) error {
	return r.db.Model(&entity.Booking{}).Where("id = ?", bookingID).Updates(update).Error
}

func (r *bookingRepository) GetBookingByID(bookingID uint64) (*entity.Booking, error) {
	var b entity.Booking
	if err := r.db.Preload("MatchGroup").Preload("Venue").Preload("Court").Where("id = ?", bookingID).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *bookingRepository) GetBookings(userID, venueID *uint64, status *string, page, limit int) ([]*entity.Booking, int64, error) {
	var list []*entity.Booking
	var total int64
	query := r.db.Model(&entity.Booking{})
	if venueID != nil {
		query = query.Where("venue_id = ?", *venueID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userID != nil {
		// join through group members to filter by user
		query = query.Joins("LEFT JOIN match_groups mg ON mg.id = bookings.match_group_id").
			Joins("LEFT JOIN match_group_members mgm ON mgm.match_group_id = mg.id").
			Where("mgm.user_id = ?", *userID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Preload("MatchGroup").Preload("Venue").Preload("Court").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *bookingRepository) DeleteBooking(bookingID uint64) error {
	return r.db.Delete(&entity.Booking{}, bookingID).Error
}
