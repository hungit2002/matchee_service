package usecase

import (
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"time"
)

// CourtUsecase defines the interface for court business logic
type CourtUsecase interface {
	CreateCourt(venueID uint64, req *entity.CreateCourtRequest) (*entity.CourtResponse, error)
	UpdateCourt(courtID uint64, req *entity.UpdateCourtRequest) (*entity.CourtResponse, error)
	GetCourtByID(courtID uint64) (*entity.CourtResponse, error)
	DeleteCourt(courtID uint64) error
	GetCourtsByVenue(venueID uint64, page, limit int) ([]*entity.CourtResponse, int64, error)
	GetAllCourts(page, limit int) ([]*entity.CourtResponse, int64, error)
}

type courtUsecase struct {
	courtRepo repository.CourtRepository
	venueRepo repository.VenueRepository
}

// NewCourtUsecase creates a new instance of CourtUsecase
func NewCourtUsecase(courtRepo repository.CourtRepository, venueRepo repository.VenueRepository) CourtUsecase {
	return &courtUsecase{
		courtRepo: courtRepo,
		venueRepo: venueRepo,
	}
}

func (c *courtUsecase) CreateCourt(venueID uint64, req *entity.CreateCourtRequest) (*entity.CourtResponse, error) {
	// Check if venue exists
	_, err := c.venueRepo.GetVenueByID(venueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	court := &entity.Court{
		VenueID:     venueID,
		Name:        req.Name,
		Description: req.Description,
		CourtType:   req.CourtType,
		Surface:     req.Surface,
		IsActive:    req.IsActive,
	}

	if err := c.courtRepo.CreateCourt(court); err != nil {
		return nil, err
	}

	// Get the created court with relationships
	createdCourt, err := c.courtRepo.GetCourtByID(court.ID)
	if err != nil {
		return nil, err
	}

	return c.convertToCourtResponse(createdCourt), nil
}

func (c *courtUsecase) UpdateCourt(courtID uint64, req *entity.UpdateCourtRequest) (*entity.CourtResponse, error) {
	// Check if court exists
	_, err := c.courtRepo.GetCourtByID(courtID)
	if err != nil {
		return nil, errors.New("court not found")
	}

	// Update fields if provided
	updateData := &entity.Court{}
	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Description != nil {
		updateData.Description = req.Description
	}
	if req.CourtType != nil {
		updateData.CourtType = *req.CourtType
	}
	if req.Surface != nil {
		updateData.Surface = *req.Surface
	}
	if req.IsActive != nil {
		updateData.IsActive = *req.IsActive
	}

	if err := c.courtRepo.UpdateCourt(courtID, updateData); err != nil {
		return nil, err
	}

	// Get the updated court
	updatedCourt, err := c.courtRepo.GetCourtByID(courtID)
	if err != nil {
		return nil, err
	}

	return c.convertToCourtResponse(updatedCourt), nil
}

func (c *courtUsecase) GetCourtByID(courtID uint64) (*entity.CourtResponse, error) {
	court, err := c.courtRepo.GetCourtByID(courtID)
	if err != nil {
		return nil, errors.New("court not found")
	}

	return c.convertToCourtResponse(court), nil
}

func (c *courtUsecase) DeleteCourt(courtID uint64) error {
	// Check if court exists
	_, err := c.courtRepo.GetCourtByID(courtID)
	if err != nil {
		return errors.New("court not found")
	}

	return c.courtRepo.DeleteCourt(courtID)
}

func (c *courtUsecase) GetCourtsByVenue(venueID uint64, page, limit int) ([]*entity.CourtResponse, int64, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	courts, total, err := c.courtRepo.GetCourtsByVenue(venueID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	courtResponses := make([]*entity.CourtResponse, len(courts))
	for i, court := range courts {
		courtResponses[i] = c.convertToCourtResponse(court)
	}

	return courtResponses, total, nil
}

func (c *courtUsecase) GetAllCourts(page, limit int) ([]*entity.CourtResponse, int64, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	courts, total, err := c.courtRepo.GetAllCourts(page, limit)
	if err != nil {
		return nil, 0, err
	}

	courtResponses := make([]*entity.CourtResponse, len(courts))
	for i, court := range courts {
		courtResponses[i] = c.convertToCourtResponse(court)
	}

	return courtResponses, total, nil
}

func (c *courtUsecase) convertToCourtResponse(court *entity.Court) *entity.CourtResponse {
	response := &entity.CourtResponse{
		ID:          court.ID,
		VenueID:     court.VenueID,
		Name:        court.Name,
		Description: court.Description,
		CourtType:   court.CourtType,
		Surface:     court.Surface,
		IsActive:    court.IsActive,
		CreatedAt:   court.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   court.UpdatedAt.Format(time.RFC3339),
	}

	if court.Venue != nil {
		response.Venue = court.Venue
	}

	if len(court.AvailabilitySlots) > 0 {
		slots := make([]*entity.AvailabilitySlot, len(court.AvailabilitySlots))
		for i, slot := range court.AvailabilitySlots {
			slots[i] = &slot
		}
		response.Slots = slots
	}

	return response
}
