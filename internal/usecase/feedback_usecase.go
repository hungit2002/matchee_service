package usecase

import (
	"errors"
	"time"

	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
)

type FeedbackUsecase interface {
	CreateFeedback(fromUserID uint64, req *entity.CreateFeedbackRequest) (*entity.FeedbackResponse, error)
	GetUserFeedbacks(userID uint64) (*entity.UserFeedbackSummary, error)
	GetBookingFeedbacks(bookingID uint64) (*entity.BookingFeedbackList, error)
}

type feedbackUsecase struct {
	feedbackRepo repository.FeedbackRepository
	bookingRepo  repository.BookingRepository
	userRepo     repository.UserRepository
}

func NewFeedbackUsecase(feedbackRepo repository.FeedbackRepository, bookingRepo repository.BookingRepository, userRepo repository.UserRepository) FeedbackUsecase {
	return &feedbackUsecase{
		feedbackRepo: feedbackRepo,
		bookingRepo:  bookingRepo,
		userRepo:     userRepo,
	}
}

func (f *feedbackUsecase) CreateFeedback(fromUserID uint64, req *entity.CreateFeedbackRequest) (*entity.FeedbackResponse, error) {
	// Verify booking exists
	booking, err := f.bookingRepo.GetBookingByID(req.BookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	// Verify target user exists
	_, err = f.userRepo.GetByID(nil, req.ToUserID)
	if err != nil {
		return nil, errors.New("target user not found")
	}

	// Check if feedback already exists for this booking and users
	existingFeedback, err := f.feedbackRepo.CheckExistingFeedback(req.BookingID, fromUserID, req.ToUserID)
	if err == nil && existingFeedback != nil {
		return nil, errors.New("feedback already exists for this booking and user")
	}

	// Validate that the booking is completed
	if booking.Status != "completed" {
		return nil, errors.New("can only provide feedback for completed bookings")
	}

	// Create feedback
	feedback := &entity.Feedback{
		BookingID:  req.BookingID,
		FromUserID: fromUserID,
		ToUserID:   req.ToUserID,
		Rating:     &req.Rating,
		Comment:    req.Comment,
	}

	if err := f.feedbackRepo.CreateFeedback(feedback); err != nil {
		return nil, err
	}

	// Get the created feedback with relationships
	createdFeedback, err := f.feedbackRepo.CheckExistingFeedback(req.BookingID, fromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}

	return f.toFeedbackResponse(createdFeedback), nil
}

func (f *feedbackUsecase) GetUserFeedbacks(userID uint64) (*entity.UserFeedbackSummary, error) {
	// Verify user exists
	_, err := f.userRepo.GetByID(nil, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	summary, err := f.feedbackRepo.GetUserFeedbackSummary(userID)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (f *feedbackUsecase) GetBookingFeedbacks(bookingID uint64) (*entity.BookingFeedbackList, error) {
	// Verify booking exists
	booking, err := f.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	feedbacks, err := f.feedbackRepo.GetFeedbacksByBooking(bookingID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	feedbackResponses := make([]*entity.FeedbackResponse, len(feedbacks))
	for i, fb := range feedbacks {
		feedbackResponses[i] = f.toFeedbackResponse(fb)
	}

	return &entity.BookingFeedbackList{
		BookingID: booking.ID,
		Feedbacks: feedbackResponses,
		Total:     int64(len(feedbackResponses)),
	}, nil
}

func (f *feedbackUsecase) toFeedbackResponse(feedback *entity.Feedback) *entity.FeedbackResponse {
	return &entity.FeedbackResponse{
		ID:         feedback.ID,
		BookingID:  feedback.BookingID,
		FromUserID: feedback.FromUserID,
		ToUserID:   feedback.ToUserID,
		Rating:     feedback.Rating,
		Comment:    feedback.Comment,
		CreatedAt:  feedback.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  feedback.UpdatedAt.Format(time.RFC3339),
		Booking:    feedback.Booking,
		FromUser:   feedback.FromUser,
		ToUser:     feedback.ToUser,
	}
}
