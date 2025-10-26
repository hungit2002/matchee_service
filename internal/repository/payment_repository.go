package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(p *entity.Payment) error
	UpdatePayment(paymentID uint64, update *entity.Payment) error
	GetPaymentByID(paymentID uint64) (*entity.Payment, error)
	GetPaymentsByPayer(payerID uint64, status *string, page, limit int) ([]*entity.Payment, int64, error)
	GetPaymentByBooking(bookingID uint64) (*entity.Payment, error)
}

type paymentRepository struct{ db *gorm.DB }

func NewPaymentRepository(db *gorm.DB) PaymentRepository { return &paymentRepository{db: db} }

func (r *paymentRepository) CreatePayment(p *entity.Payment) error { return r.db.Create(p).Error }

func (r *paymentRepository) UpdatePayment(paymentID uint64, update *entity.Payment) error {
	return r.db.Model(&entity.Payment{}).Where("id = ?", paymentID).Updates(update).Error
}

func (r *paymentRepository) GetPaymentByID(paymentID uint64) (*entity.Payment, error) {
	var p entity.Payment
	if err := r.db.Preload("Booking").Preload("Payer").Where("id = ?", paymentID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) GetPaymentsByPayer(payerID uint64, status *string, page, limit int) ([]*entity.Payment, int64, error) {
	var list []*entity.Payment
	var total int64
	query := r.db.Model(&entity.Payment{}).Where("payer_id = ?", payerID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Preload("Booking").Preload("Payer").Offset(offset).Limit(limit).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *paymentRepository) GetPaymentByBooking(bookingID uint64) (*entity.Payment, error) {
	var p entity.Payment
	if err := r.db.Preload("Booking").Preload("Payer").Where("booking_id = ?", bookingID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
