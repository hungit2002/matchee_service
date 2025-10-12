package usecase

import (
	"context"
	"errors"
	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
	"time"
)

// MatchPostUsecase defines the interface for match post business logic
type MatchPostUsecase interface {
	CreateMatchPost(userID uint64, req *entity.CreateMatchPostRequest) (*entity.MatchPostResponse, error)
	UpdateMatchPost(matchPostID uint64, userID uint64, req *entity.UpdateMatchPostRequest) (*entity.MatchPostResponse, error)
	GetMatchPostByID(matchPostID uint64) (*entity.MatchPostResponse, error)
	DeleteMatchPost(matchPostID uint64, userID uint64) error
	GetMatchPostsByUser(userID uint64, page, limit int) (*entity.MatchPostListResponse, error)
	GetMatchPostsByLocation(req *entity.MatchPostListRequest) (*entity.MatchPostListResponse, error)
	GetMatchPostsByVenue(venueID uint64, req *entity.MatchPostListRequest) (*entity.MatchPostListResponse, error)
	GetAllMatchPosts(req *entity.MatchPostListRequest) (*entity.MatchPostListResponse, error)
	GetSuggestedMatchPosts(req *entity.MatchPostSuggestRequest) (*entity.MatchPostSuggestResponse, error)
}

type matchPostUsecase struct {
	matchPostRepo repository.MatchPostRepository
	userRepo      repository.UserRepository
	venueRepo     repository.VenueRepository
}

// NewMatchPostUsecase creates a new instance of MatchPostUsecase
func NewMatchPostUsecase(matchPostRepo repository.MatchPostRepository, userRepo repository.UserRepository, venueRepo repository.VenueRepository) MatchPostUsecase {
	return &matchPostUsecase{
		matchPostRepo: matchPostRepo,
		userRepo:      userRepo,
		venueRepo:     venueRepo,
	}
}

func (m *matchPostUsecase) CreateMatchPost(userID uint64, req *entity.CreateMatchPostRequest) (*entity.MatchPostResponse, error) {
	// Check if user exists
	_, err := m.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if venue exists (if provided)
	if req.VenueID != nil {
		_, err := m.venueRepo.GetVenueByID(*req.VenueID)
		if err != nil {
			return nil, errors.New("venue not found")
		}
	}

	matchPost := &entity.MatchPost{
		UserID:         userID,
		DesiredLevel:   req.DesiredLevel,
		Location:       req.Location,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		VenueID:        req.VenueID,
		DesiredTime:    req.DesiredTime,
		PricePerPerson: req.PricePerPerson,
		MaxPlayers:     req.MaxPlayers,
		Note:           req.Note,
		Status:         "open",
	}

	if err := m.matchPostRepo.CreateMatchPost(matchPost); err != nil {
		return nil, err
	}

	// Get the created match post with relationships
	createdMatchPost, err := m.matchPostRepo.GetMatchPostByID(matchPost.ID)
	if err != nil {
		return nil, err
	}

	return m.convertToMatchPostResponse(createdMatchPost), nil
}

func (m *matchPostUsecase) UpdateMatchPost(matchPostID uint64, userID uint64, req *entity.UpdateMatchPostRequest) (*entity.MatchPostResponse, error) {
	// Check if match post exists and belongs to user
	existingMatchPost, err := m.matchPostRepo.GetMatchPostByID(matchPostID)
	if err != nil {
		return nil, errors.New("match post not found")
	}

	if existingMatchPost.UserID != userID {
		return nil, errors.New("unauthorized to update this match post")
	}

	// Update fields if provided
	updateData := &entity.MatchPost{}
	if req.DesiredLevel != nil {
		updateData.DesiredLevel = req.DesiredLevel
	}
	if req.Location != nil {
		updateData.Location = req.Location
	}
	if req.Latitude != nil {
		updateData.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		updateData.Longitude = req.Longitude
	}
	if req.VenueID != nil {
		// Check if venue exists
		_, err := m.venueRepo.GetVenueByID(*req.VenueID)
		if err != nil {
			return nil, errors.New("venue not found")
		}
		updateData.VenueID = req.VenueID
	}
	if req.DesiredTime != nil {
		updateData.DesiredTime = req.DesiredTime
	}
	if req.PricePerPerson != nil {
		updateData.PricePerPerson = req.PricePerPerson
	}
	if req.MaxPlayers != nil {
		updateData.MaxPlayers = *req.MaxPlayers
	}
	if req.Note != nil {
		updateData.Note = req.Note
	}
	if req.Status != nil {
		updateData.Status = *req.Status
	}

	if err := m.matchPostRepo.UpdateMatchPost(matchPostID, updateData); err != nil {
		return nil, err
	}

	// Get the updated match post
	updatedMatchPost, err := m.matchPostRepo.GetMatchPostByID(matchPostID)
	if err != nil {
		return nil, err
	}

	return m.convertToMatchPostResponse(updatedMatchPost), nil
}

func (m *matchPostUsecase) GetMatchPostByID(matchPostID uint64) (*entity.MatchPostResponse, error) {
	matchPost, err := m.matchPostRepo.GetMatchPostByID(matchPostID)
	if err != nil {
		return nil, errors.New("match post not found")
	}

	return m.convertToMatchPostResponse(matchPost), nil
}

func (m *matchPostUsecase) DeleteMatchPost(matchPostID uint64, userID uint64) error {
	// Check if match post exists and belongs to user
	existingMatchPost, err := m.matchPostRepo.GetMatchPostByID(matchPostID)
	if err != nil {
		return errors.New("match post not found")
	}

	if existingMatchPost.UserID != userID {
		return errors.New("unauthorized to delete this match post")
	}

	return m.matchPostRepo.DeleteMatchPost(matchPostID)
}

func (m *matchPostUsecase) GetMatchPostsByUser(userID uint64, page, limit int) (*entity.MatchPostListResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	matchPosts, total, err := m.matchPostRepo.GetMatchPostsByUser(userID, page, limit)
	if err != nil {
		return nil, err
	}

	matchPostResponses := make([]*entity.MatchPostResponse, len(matchPosts))
	for i, matchPost := range matchPosts {
		matchPostResponses[i] = m.convertToMatchPostResponse(matchPost)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.MatchPostListResponse{
		MatchPosts: matchPostResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (m *matchPostUsecase) GetMatchPostsByLocation(req *entity.MatchPostListRequest) (*entity.MatchPostListResponse, error) {
	// Set default values
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	var matchPosts []*entity.MatchPost
	var total int64
	var err error

	if req.Latitude != nil && req.Longitude != nil && req.Radius != nil {
		matchPosts, total, err = m.matchPostRepo.GetMatchPostsByLocation(*req.Latitude, *req.Longitude, *req.Radius, req.DesiredLevel, req.Status, page, limit)
	} else {
		matchPosts, total, err = m.matchPostRepo.GetAllMatchPosts(req.Status, page, limit)
	}

	if err != nil {
		return nil, err
	}

	matchPostResponses := make([]*entity.MatchPostResponse, len(matchPosts))
	for i, matchPost := range matchPosts {
		matchPostResponses[i] = m.convertToMatchPostResponse(matchPost)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.MatchPostListResponse{
		MatchPosts: matchPostResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (m *matchPostUsecase) GetMatchPostsByVenue(venueID uint64, req *entity.MatchPostListRequest) (*entity.MatchPostListResponse, error) {
	// Set default values
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	matchPosts, total, err := m.matchPostRepo.GetMatchPostsByVenue(venueID, req.Status, page, limit)
	if err != nil {
		return nil, err
	}

	matchPostResponses := make([]*entity.MatchPostResponse, len(matchPosts))
	for i, matchPost := range matchPosts {
		matchPostResponses[i] = m.convertToMatchPostResponse(matchPost)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.MatchPostListResponse{
		MatchPosts: matchPostResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (m *matchPostUsecase) GetAllMatchPosts(req *entity.MatchPostListRequest) (*entity.MatchPostListResponse, error) {
	// Set default values
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	matchPosts, total, err := m.matchPostRepo.GetAllMatchPosts(req.Status, page, limit)
	if err != nil {
		return nil, err
	}

	matchPostResponses := make([]*entity.MatchPostResponse, len(matchPosts))
	for i, matchPost := range matchPosts {
		matchPostResponses[i] = m.convertToMatchPostResponse(matchPost)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &entity.MatchPostListResponse{
		MatchPosts: matchPostResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (m *matchPostUsecase) GetSuggestedMatchPosts(req *entity.MatchPostSuggestRequest) (*entity.MatchPostSuggestResponse, error) {
	// Set default limit
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	matchPosts, total, err := m.matchPostRepo.GetSuggestedMatchPosts(req.UserID, req.Latitude, req.Longitude, req.Radius, req.DesiredLevel, limit)
	if err != nil {
		return nil, err
	}

	matchPostResponses := make([]*entity.MatchPostResponse, len(matchPosts))
	for i, matchPost := range matchPosts {
		matchPostResponses[i] = m.convertToMatchPostResponse(matchPost)
	}

	return &entity.MatchPostSuggestResponse{
		Suggestions: matchPostResponses,
		Total:       total,
	}, nil
}

func (m *matchPostUsecase) convertToMatchPostResponse(matchPost *entity.MatchPost) *entity.MatchPostResponse {
	response := &entity.MatchPostResponse{
		ID:             matchPost.ID,
		UserID:         matchPost.UserID,
		DesiredLevel:   matchPost.DesiredLevel,
		Location:       matchPost.Location,
		Latitude:       matchPost.Latitude,
		Longitude:      matchPost.Longitude,
		VenueID:        matchPost.VenueID,
		PricePerPerson: matchPost.PricePerPerson,
		MaxPlayers:     matchPost.MaxPlayers,
		Note:           matchPost.Note,
		Status:         matchPost.Status,
		CreatedAt:      matchPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      matchPost.UpdatedAt.Format(time.RFC3339),
	}

	if matchPost.DesiredTime != nil {
		desiredTimeStr := matchPost.DesiredTime.Format(time.RFC3339)
		response.DesiredTime = &desiredTimeStr
	}

	if matchPost.User != nil {
		response.User = matchPost.User
	}

	if matchPost.Venue != nil {
		response.Venue = matchPost.Venue
	}

	if len(matchPost.MatchGroups) > 0 {
		groups := make([]*entity.MatchGroup, len(matchPost.MatchGroups))
		for i, group := range matchPost.MatchGroups {
			groups[i] = &group
		}
		response.MatchGroups = groups
	}

	return response
}
