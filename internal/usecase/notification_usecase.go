package usecase

import (
	"context"
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"matchee/services/internal/service"
	"time"
)

type NotificationUsecase interface {
	GetMyNotifications(userID uint64, onlyUnread bool, page, limit int) (*entity.NotificationListResponse, error)
	MarkRead(userID, notificationID uint64) error
	CreateNotification(req *entity.CreateNotificationRequest) (*entity.NotificationResponse, error)
}

type notificationUsecase struct {
	repo      repository.NotificationRepository
	publisher service.NotificationPublisher
}

func NewNotificationUsecase(repo repository.NotificationRepository, publisher service.NotificationPublisher) NotificationUsecase {
	return &notificationUsecase{repo: repo, publisher: publisher}
}

func (u *notificationUsecase) GetMyNotifications(userID uint64, onlyUnread bool, page, limit int) (*entity.NotificationListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	list, total, err := u.repo.ListByUser(userID, onlyUnread, page, limit)
	if err != nil {
		return nil, err
	}
	out := make([]*entity.NotificationResponse, len(list))
	for i, n := range list {
		out[i] = toNotificationResponse(n)
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &entity.NotificationListResponse{Notifications: out, Total: total, Page: page, Limit: limit, TotalPages: totalPages}, nil
}

func (u *notificationUsecase) MarkRead(userID, notificationID uint64) error {
	// ensure notification belongs to user
	n, err := u.repo.GetByID(notificationID)
	if err != nil {
		return err
	}
	if n.UserID != userID {
		return errors.New("forbidden")
	}
	return u.repo.MarkRead(notificationID)
}

func (u *notificationUsecase) CreateNotification(req *entity.CreateNotificationRequest) (*entity.NotificationResponse, error) {
	n := &entity.Notification{UserID: req.UserID, Title: req.Title, Message: req.Message}
	if err := u.repo.Create(n); err != nil {
		return nil, err
	}
	if u.publisher != nil {
		_ = u.publisher.Publish(context.Background(), "notifications", toNotificationResponse(n))
	}
	return toNotificationResponse(n), nil
}

func toNotificationResponse(n *entity.Notification) *entity.NotificationResponse {
	return &entity.NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
}
