package usecase

import (
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"time"
)

type BookingUsecase interface {
	CreateBooking(req *entity.CreateBookingRequest) (*entity.BookingResponse, error)
	GetBookings(req *entity.BookingListRequest) (*entity.BookingListResponse, error)
	GetBookingByID(bookingID uint64) (*entity.BookingResponse, error)
	UpdateBookingStatus(bookingID uint64, req *entity.UpdateBookingStatusRequest) (*entity.BookingResponse, error)
	DeleteBooking(bookingID uint64) error
}

type bookingUsecase struct {
	bookingRepo repository.BookingRepository
	groupRepo   repository.MatchGroupRepository
}

func NewBookingUsecase(bookingRepo repository.BookingRepository, groupRepo repository.MatchGroupRepository) BookingUsecase {
	return &bookingUsecase{bookingRepo: bookingRepo, groupRepo: groupRepo}
}

func (u *bookingUsecase) CreateBooking(req *entity.CreateBookingRequest) (*entity.BookingResponse, error) {
	if req.StartTime == nil || req.EndTime == nil {
		return nil, errors.New("startTime and endTime are required")
	}
	if req.EndTime.Before(*req.StartTime) {
		return nil, errors.New("endTime must be after startTime")
	}
	b := &entity.Booking{
		MatchGroupID: req.MatchGroupID,
		VenueID:      req.VenueID,
		CourtID:      req.CourtID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		TotalPrice:   req.TotalPrice,
		Status:       "reserved",
	}
	if err := u.bookingRepo.CreateBooking(b); err != nil {
		return nil, err
	}
	created, err := u.bookingRepo.GetBookingByID(b.ID)
	if err != nil {
		return nil, err
	}
	return u.toBookingResponse(created), nil
}

func (u *bookingUsecase) GetBookings(req *entity.BookingListRequest) (*entity.BookingListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	list, total, err := u.bookingRepo.GetBookings(req.UserID, req.VenueID, req.Status, page, limit)
	if err != nil {
		return nil, err
	}
	res := make([]*entity.BookingResponse, len(list))
	for i, b := range list {
		res[i] = u.toBookingResponse(b)
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &entity.BookingListResponse{Bookings: res, Total: total, Page: page, Limit: limit, TotalPages: totalPages}, nil
}

func (u *bookingUsecase) GetBookingByID(bookingID uint64) (*entity.BookingResponse, error) {
	b, err := u.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	return u.toBookingResponse(b), nil
}

func (u *bookingUsecase) UpdateBookingStatus(bookingID uint64, req *entity.UpdateBookingStatusRequest) (*entity.BookingResponse, error) {
	update := &entity.Booking{Status: req.Status}
	if err := u.bookingRepo.UpdateBooking(bookingID, update); err != nil {
		return nil, err
	}
	b, err := u.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		return nil, err
	}
	return u.toBookingResponse(b), nil
}

func (u *bookingUsecase) DeleteBooking(bookingID uint64) error {
	return u.bookingRepo.DeleteBooking(bookingID)
}

func (u *bookingUsecase) toBookingResponse(b *entity.Booking) *entity.BookingResponse {
	resp := &entity.BookingResponse{
		ID:           b.ID,
		MatchGroupID: b.MatchGroupID,
		VenueID:      b.VenueID,
		CourtID:      b.CourtID,
		TotalPrice:   b.TotalPrice,
		Status:       b.Status,
		CreatedAt:    b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    b.UpdatedAt.Format(time.RFC3339),
	}
	if b.StartTime != nil {
		s := b.StartTime.Format(time.RFC3339)
		resp.StartTime = &s
	}
	if b.EndTime != nil {
		e := b.EndTime.Format(time.RFC3339)
		resp.EndTime = &e
	}
	if b.MatchGroup != nil {
		resp.MatchGroup = b.MatchGroup
	}
	if b.Venue != nil {
		resp.Venue = b.Venue
	}
	if b.Court != nil {
		resp.Court = b.Court
	}
	return resp
}
