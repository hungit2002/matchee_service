package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(n *entity.Notification) error
	ListByUser(userID uint64, onlyUnread bool, page, limit int) ([]*entity.Notification, int64, error)
	MarkRead(notificationID uint64) error
	GetByID(id uint64) (*entity.Notification, error)
}

type notificationRepository struct{ db *gorm.DB }

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(n *entity.Notification) error { return r.db.Create(n).Error }

func (r *notificationRepository) ListByUser(userID uint64, onlyUnread bool, page, limit int) ([]*entity.Notification, int64, error) {
	var list []*entity.Notification
	var total int64
	q := r.db.Model(&entity.Notification{}).Where("user_id = ?", userID)
	if onlyUnread {
		q = q.Where("is_read = ?", false)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *notificationRepository) MarkRead(notificationID uint64) error {
	return r.db.Model(&entity.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (r *notificationRepository) GetByID(id uint64) (*entity.Notification, error) {
	var n entity.Notification
	if err := r.db.Where("id = ?", id).First(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}
