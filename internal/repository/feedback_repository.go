package repository

import (
	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type FeedbackRepository interface {
	CreateFeedback(f *entity.Feedback) error
	GetFeedbacksByUser(userID uint64) ([]*entity.Feedback, error)
	GetFeedbacksByBooking(bookingID uint64) ([]*entity.Feedback, error)
	GetUserFeedbackSummary(userID uint64) (*entity.UserFeedbackSummary, error)
	CheckExistingFeedback(bookingID, fromUserID, toUserID uint64) (*entity.Feedback, error)
}

type feedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) CreateFeedback(f *entity.Feedback) error {
	return r.db.Create(f).Error
}

func (r *feedbackRepository) GetFeedbacksByUser(userID uint64) ([]*entity.Feedback, error) {
	var feedbacks []*entity.Feedback
	err := r.db.Preload("FromUser").Preload("ToUser").Preload("Booking").
		Where("to_user_id = ?", userID).
		Order("created_at DESC").
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) GetFeedbacksByBooking(bookingID uint64) ([]*entity.Feedback, error) {
	var feedbacks []*entity.Feedback
	err := r.db.Preload("FromUser").Preload("ToUser").Preload("Booking").
		Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) GetUserFeedbackSummary(userID uint64) (*entity.UserFeedbackSummary, error) {
	var summary entity.UserFeedbackSummary
	summary.UserID = userID
	summary.RatingCounts = make(map[int]int64)

	// Get total ratings count
	err := r.db.Model(&entity.Feedback{}).
		Where("to_user_id = ? AND rating IS NOT NULL", userID).
		Count(&summary.TotalRatings).Error
	if err != nil {
		return nil, err
	}

	if summary.TotalRatings > 0 {
		// Calculate average rating
		var avgRating float64
		err = r.db.Model(&entity.Feedback{}).
			Where("to_user_id = ? AND rating IS NOT NULL", userID).
			Select("AVG(rating)").
			Scan(&avgRating).Error
		if err != nil {
			return nil, err
		}
		summary.AverageRating = avgRating

		// Get rating counts for each star (1-5)
		for i := 1; i <= 5; i++ {
			var count int64
			err = r.db.Model(&entity.Feedback{}).
				Where("to_user_id = ? AND rating = ?", userID, i).
				Count(&count).Error
			if err != nil {
				return nil, err
			}
			summary.RatingCounts[i] = count
		}
	}

	// Get recent feedbacks (last 10)
	var recentFeedbacks []*entity.Feedback
	err = r.db.Preload("FromUser").Preload("Booking").
		Where("to_user_id = ?", userID).
		Order("created_at DESC").
		Limit(10).
		Find(&recentFeedbacks).Error
	if err != nil {
		return nil, err
	}

	// Convert to response format
	summary.RecentFeedbacks = make([]*entity.FeedbackResponse, len(recentFeedbacks))
	for i, fb := range recentFeedbacks {
		summary.RecentFeedbacks[i] = r.toFeedbackResponse(fb)
	}

	return &summary, nil
}

func (r *feedbackRepository) CheckExistingFeedback(bookingID, fromUserID, toUserID uint64) (*entity.Feedback, error) {
	var feedback entity.Feedback
	err := r.db.Where("booking_id = ? AND from_user_id = ? AND to_user_id = ?",
		bookingID, fromUserID, toUserID).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *feedbackRepository) toFeedbackResponse(f *entity.Feedback) *entity.FeedbackResponse {
	return &entity.FeedbackResponse{
		ID:         f.ID,
		BookingID:  f.BookingID,
		FromUserID: f.FromUserID,
		ToUserID:   f.ToUserID,
		Rating:     f.Rating,
		Comment:    f.Comment,
		CreatedAt:  f.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  f.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Booking:    f.Booking,
		FromUser:   f.FromUser,
		ToUser:     f.ToUser,
	}
}
