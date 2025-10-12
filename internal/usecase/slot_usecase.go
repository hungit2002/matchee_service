package usecase

import (
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"time"
)

// SlotUsecase defines the interface for availability slot business logic
type SlotUsecase interface {
	CreateSlot(courtID uint64, req *entity.CreateSlotRequest) (*entity.SlotResponse, error)
	UpdateSlot(slotID uint64, req *entity.UpdateSlotRequest) (*entity.SlotResponse, error)
	GetSlotByID(slotID uint64) (*entity.SlotResponse, error)
	DeleteSlot(slotID uint64) error
	GetSlotsByCourt(courtID uint64, req *entity.SlotListRequest) (*entity.SlotListResponse, error)
	GetAllSlots(page, limit int) (*entity.SlotListResponse, error)
}

type slotUsecase struct {
	slotRepo  repository.SlotRepository
	courtRepo repository.CourtRepository
}

// NewSlotUsecase creates a new instance of SlotUsecase
func NewSlotUsecase(slotRepo repository.SlotRepository, courtRepo repository.CourtRepository) SlotUsecase {
	return &slotUsecase{
		slotRepo:  slotRepo,
		courtRepo: courtRepo,
	}
}

func (s *slotUsecase) CreateSlot(courtID uint64, req *entity.CreateSlotRequest) (*entity.SlotResponse, error) {
	// Check if court exists
	_, err := s.courtRepo.GetCourtByID(courtID)
	if err != nil {
		return nil, errors.New("court not found")
	}

	// Validate time range
	if req.StartTime.After(req.EndTime) {
		return nil, errors.New("start time must be before end time")
	}

	// Check for overlapping slots
	existingSlots, _, err := s.slotRepo.GetSlotsByCourt(courtID, &req.StartTime, &req.EndTime, nil, 1, 100)
	if err != nil {
		return nil, err
	}

	for _, slot := range existingSlots {
		if req.StartTime.Before(slot.EndTime) && req.EndTime.After(slot.StartTime) {
			return nil, errors.New("slot overlaps with existing slot")
		}
	}

	slot := &entity.AvailabilitySlot{
		CourtID:     courtID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsAvailable: req.IsAvailable,
	}

	if err := s.slotRepo.CreateSlot(slot); err != nil {
		return nil, err
	}

	// Get the created slot with relationships
	createdSlot, err := s.slotRepo.GetSlotByID(slot.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToSlotResponse(createdSlot), nil
}

func (s *slotUsecase) UpdateSlot(slotID uint64, req *entity.UpdateSlotRequest) (*entity.SlotResponse, error) {
	// Check if slot exists
	_, err := s.slotRepo.GetSlotByID(slotID)
	if err != nil {
		return nil, errors.New("slot not found")
	}

	// Update fields if provided
	updateData := &entity.AvailabilitySlot{}
	if req.StartTime != nil {
		updateData.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		updateData.EndTime = *req.EndTime
	}
	if req.IsAvailable != nil {
		updateData.IsAvailable = *req.IsAvailable
	}

	// Validate time range if both times are provided
	if req.StartTime != nil && req.EndTime != nil {
		if req.StartTime.After(*req.EndTime) {
			return nil, errors.New("start time must be before end time")
		}
	}

	if err := s.slotRepo.UpdateSlot(slotID, updateData); err != nil {
		return nil, err
	}

	// Get the updated slot
	updatedSlot, err := s.slotRepo.GetSlotByID(slotID)
	if err != nil {
		return nil, err
	}

	return s.convertToSlotResponse(updatedSlot), nil
}

func (s *slotUsecase) GetSlotByID(slotID uint64) (*entity.SlotResponse, error) {
	slot, err := s.slotRepo.GetSlotByID(slotID)
	if err != nil {
		return nil, errors.New("slot not found")
	}

	return s.convertToSlotResponse(slot), nil
}

func (s *slotUsecase) DeleteSlot(slotID uint64) error {
	// Check if slot exists
	_, err := s.slotRepo.GetSlotByID(slotID)
	if err != nil {
		return errors.New("slot not found")
	}

	return s.slotRepo.DeleteSlot(slotID)
}

func (s *slotUsecase) GetSlotsByCourt(courtID uint64, req *entity.SlotListRequest) (*entity.SlotListResponse, error) {
	// Set default values
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	slots, total, err := s.slotRepo.GetSlotsByCourt(courtID, req.StartDate, req.EndDate, req.IsAvailable, page, limit)
	if err != nil {
		return nil, err
	}

	slotResponses := make([]*entity.SlotResponse, len(slots))
	for i, slot := range slots {
		slotResponses[i] = s.convertToSlotResponse(slot)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.SlotListResponse{
		Slots:      slotResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *slotUsecase) GetAllSlots(page, limit int) (*entity.SlotListResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	slots, total, err := s.slotRepo.GetAllSlots(page, limit)
	if err != nil {
		return nil, err
	}

	slotResponses := make([]*entity.SlotResponse, len(slots))
	for i, slot := range slots {
		slotResponses[i] = s.convertToSlotResponse(slot)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.SlotListResponse{
		Slots:      slotResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *slotUsecase) convertToSlotResponse(slot *entity.AvailabilitySlot) *entity.SlotResponse {
	response := &entity.SlotResponse{
		ID:          slot.ID,
		CourtID:     slot.CourtID,
		StartTime:   slot.StartTime.Format(time.RFC3339),
		EndTime:     slot.EndTime.Format(time.RFC3339),
		IsAvailable: slot.IsAvailable,
		CreatedAt:   slot.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   slot.UpdatedAt.Format(time.RFC3339),
	}

	if slot.Court != nil {
		response.Court = slot.Court
	}

	return response
}
