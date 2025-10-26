package usecase

import (
	"errors"
	"time"

	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
)

type AdminUsecase interface {
	GetRevenueStatistics(fromDate, toDate string) (*entity.RevenueStatisticsResponse, error)
	GetUsers(req *entity.AdminUserListRequest) (*entity.AdminUserListResponse, error)
	LockUser(userID uint64, req *entity.LockUserRequest) error
	ApproveVenue(venueID uint64, req *entity.VenueApprovalRequest) error
	GetFeedbacks(req *entity.AdminFeedbackListRequest) (*entity.AdminFeedbackListResponse, error)
}

type adminUsecase struct {
	adminRepo    repository.AdminRepository
	userRepo     repository.UserRepository
	venueRepo    repository.VenueRepository
	feedbackRepo repository.FeedbackRepository
}

func NewAdminUsecase(adminRepo repository.AdminRepository, userRepo repository.UserRepository, venueRepo repository.VenueRepository, feedbackRepo repository.FeedbackRepository) AdminUsecase {
	return &adminUsecase{
		adminRepo:    adminRepo,
		userRepo:     userRepo,
		venueRepo:    venueRepo,
		feedbackRepo: feedbackRepo,
	}
}

func (a *adminUsecase) GetRevenueStatistics(fromDate, toDate string) (*entity.RevenueStatisticsResponse, error) {
	// Validate date format if provided
	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			return nil, errors.New("invalid fromDate format, use YYYY-MM-DD")
		}
	}
	if toDate != "" {
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			return nil, errors.New("invalid toDate format, use YYYY-MM-DD")
		}
	}

	stats, totalRevenue, err := a.adminRepo.GetRevenueStatistics(fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Determine period description
	period := "All time"
	if fromDate != "" && toDate != "" {
		period = fromDate + " to " + toDate
	} else if fromDate != "" {
		period = "From " + fromDate
	} else if toDate != "" {
		period = "Until " + toDate
	}

	return &entity.RevenueStatisticsResponse{
		Period: period,
		Data:   stats,
		Total:  totalRevenue,
	}, nil
}

func (a *adminUsecase) GetUsers(req *entity.AdminUserListRequest) (*entity.AdminUserListResponse, error) {
	users, total, err := a.adminRepo.GetUsersWithStats(req)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.AdminUserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (a *adminUsecase) LockUser(userID uint64, req *entity.LockUserRequest) error {
	// Verify user exists
	_, err := a.userRepo.GetByID(nil, userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = a.adminRepo.LockUser(userID, req.Locked, req.Reason)
	if err != nil {
		return err
	}

	return nil
}

func (a *adminUsecase) ApproveVenue(venueID uint64, req *entity.VenueApprovalRequest) error {
	// Verify venue exists
	_, err := a.venueRepo.GetVenueByID(venueID)
	if err != nil {
		return errors.New("venue not found")
	}

	err = a.adminRepo.ApproveVenue(venueID, req.Approved, req.Reason)
	if err != nil {
		return err
	}

	return nil
}

func (a *adminUsecase) GetFeedbacks(req *entity.AdminFeedbackListRequest) (*entity.AdminFeedbackListResponse, error) {
	feedbacks, total, err := a.adminRepo.GetFeedbacksWithDetails(req)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.AdminFeedbackListResponse{
		Feedbacks:  feedbacks,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
