package usecase

import (
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"time"
)

// VenueUsecase defines the interface for venue business logic
type VenueUsecase interface {
	CreateVenue(ownerID uint64, req *entity.CreateVenueRequest) (*entity.VenueResponse, error)
	UpdateVenue(venueID uint64, req *entity.UpdateVenueRequest) (*entity.VenueResponse, error)
	GetVenueByID(venueID uint64) (*entity.VenueResponse, error)
	DeleteVenue(venueID uint64) error
	GetVenuesByOwner(ownerID uint64, page, limit int) (*entity.VenueListResponse, error)
	GetVenuesByLocation(req *entity.VenueListRequest) (*entity.VenueListResponse, error)
	GetAllVenues(page, limit int) (*entity.VenueListResponse, error)
}

type venueUsecase struct {
	venueRepo repository.VenueRepository
}

// NewVenueUsecase creates a new instance of VenueUsecase
func NewVenueUsecase(venueRepo repository.VenueRepository) VenueUsecase {
	return &venueUsecase{
		venueRepo: venueRepo,
	}
}

func (v *venueUsecase) CreateVenue(ownerID uint64, req *entity.CreateVenueRequest) (*entity.VenueResponse, error) {
	venue := &entity.Venue{
		OwnerID:      ownerID,
		Name:         req.Name,
		Description:  req.Description,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Phone:        req.Phone,
		Email:        req.Email,
		Website:      req.Website,
		PricePerHour: req.PricePerHour,
		IsActive:     req.IsActive,
	}

	if err := v.venueRepo.CreateVenue(venue); err != nil {
		return nil, err
	}

	// Get the created venue with relationships
	createdVenue, err := v.venueRepo.GetVenueByID(venue.ID)
	if err != nil {
		return nil, err
	}

	return v.convertToVenueResponse(createdVenue), nil
}

func (v *venueUsecase) UpdateVenue(venueID uint64, req *entity.UpdateVenueRequest) (*entity.VenueResponse, error) {
	// Check if venue exists
	_, err := v.venueRepo.GetVenueByID(venueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	// Update fields if provided
	updateData := &entity.Venue{}
	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Description != nil {
		updateData.Description = req.Description
	}
	if req.Address != nil {
		updateData.Address = *req.Address
	}
	if req.Latitude != nil {
		updateData.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		updateData.Longitude = *req.Longitude
	}
	if req.Phone != nil {
		updateData.Phone = *req.Phone
	}
	if req.Email != nil {
		updateData.Email = req.Email
	}
	if req.Website != nil {
		updateData.Website = req.Website
	}
	if req.PricePerHour != nil {
		updateData.PricePerHour = *req.PricePerHour
	}
	if req.IsActive != nil {
		updateData.IsActive = *req.IsActive
	}

	if err := v.venueRepo.UpdateVenue(venueID, updateData); err != nil {
		return nil, err
	}

	// Get the updated venue
	updatedVenue, err := v.venueRepo.GetVenueByID(venueID)
	if err != nil {
		return nil, err
	}

	return v.convertToVenueResponse(updatedVenue), nil
}

func (v *venueUsecase) GetVenueByID(venueID uint64) (*entity.VenueResponse, error) {
	venue, err := v.venueRepo.GetVenueByID(venueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	return v.convertToVenueResponse(venue), nil
}

func (v *venueUsecase) DeleteVenue(venueID uint64) error {
	// Check if venue exists
	_, err := v.venueRepo.GetVenueByID(venueID)
	if err != nil {
		return errors.New("venue not found")
	}

	return v.venueRepo.DeleteVenue(venueID)
}

func (v *venueUsecase) GetVenuesByOwner(ownerID uint64, page, limit int) (*entity.VenueListResponse, error) {
	venues, total, err := v.venueRepo.GetVenuesByOwner(ownerID, page, limit)
	if err != nil {
		return nil, err
	}

	venueResponses := make([]*entity.VenueResponse, len(venues))
	for i, venue := range venues {
		venueResponses[i] = v.convertToVenueResponse(venue)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.VenueListResponse{
		Venues:     venueResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (v *venueUsecase) GetVenuesByLocation(req *entity.VenueListRequest) (*entity.VenueListResponse, error) {
	// Set default values
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	var venues []*entity.Venue
	var total int64
	var err error

	if req.Latitude != nil && req.Longitude != nil && req.Radius != nil {
		venues, total, err = v.venueRepo.GetVenuesByLocation(*req.Latitude, *req.Longitude, *req.Radius, req.MinPrice, req.MaxPrice, req.IsActive, page, limit)
	} else {
		venues, total, err = v.venueRepo.GetAllVenues(page, limit)
	}

	if err != nil {
		return nil, err
	}

	venueResponses := make([]*entity.VenueResponse, len(venues))
	for i, venue := range venues {
		venueResponses[i] = v.convertToVenueResponse(venue)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.VenueListResponse{
		Venues:     venueResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (v *venueUsecase) GetAllVenues(page, limit int) (*entity.VenueListResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	venues, total, err := v.venueRepo.GetAllVenues(page, limit)
	if err != nil {
		return nil, err
	}

	venueResponses := make([]*entity.VenueResponse, len(venues))
	for i, venue := range venues {
		venueResponses[i] = v.convertToVenueResponse(venue)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.VenueListResponse{
		Venues:     venueResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (v *venueUsecase) convertToVenueResponse(venue *entity.Venue) *entity.VenueResponse {
	response := &entity.VenueResponse{
		ID:           venue.ID,
		OwnerID:      venue.OwnerID,
		Name:         venue.Name,
		Description:  venue.Description,
		Address:      venue.Address,
		Latitude:     venue.Latitude,
		Longitude:    venue.Longitude,
		Phone:        venue.Phone,
		Email:        venue.Email,
		Website:      venue.Website,
		PricePerHour: venue.PricePerHour,
		IsActive:     venue.IsActive,
		CreatedAt:    venue.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    venue.UpdatedAt.Format(time.RFC3339),
	}

	if venue.Owner != nil {
		response.Owner = venue.Owner
	}

	if len(venue.Courts) > 0 {
		courts := make([]*entity.Court, len(venue.Courts))
		for i, court := range venue.Courts {
			courts[i] = &court
		}
		response.Courts = courts
	}

	return response
}
