package usecase

import (
	"errors"
	"time"

	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
)

type PlayerUsecase interface {
	CreateOrUpdatePlayerProfile(userID uint64, req *entity.CreateUpdatePlayerProfileRequest) (*entity.PlayerProfileResponse, error)
	GetPlayerProfileByID(id uint64) (*entity.PlayerProfileResponse, error)
	GetPlayerProfileByUserID(userID uint64) (*entity.PlayerProfileResponse, error)
	GetPlayerSuggestions(userID uint64, req *entity.PlayerSuggestionRequest) (*entity.PlayerSuggestionResponse, error)
}

type playerUsecase struct {
	playerRepo repository.PlayerRepository
}

func NewPlayerUsecase(playerRepo repository.PlayerRepository) PlayerUsecase {
	return &playerUsecase{
		playerRepo: playerRepo,
	}
}

func (p *playerUsecase) CreateOrUpdatePlayerProfile(userID uint64, req *entity.CreateUpdatePlayerProfileRequest) (*entity.PlayerProfileResponse, error) {
	// Check if profile already exists
	existingProfile, err := p.playerRepo.GetPlayerProfileByUserID(userID)
	if err != nil && err.Error() != "record not found" {
		return nil, errors.New("failed to check existing profile")
	}

	var profile *entity.PlayerProfile
	if existingProfile != nil {
		// Update existing profile
		profile = existingProfile
		profile.Level = req.Level
		profile.Gender = req.Gender
		profile.PreferredLocation = req.PreferredLocation
		profile.Latitude = req.Latitude
		profile.Longitude = req.Longitude
		profile.Bio = req.Bio
		profile.UpdatedAt = time.Now()

		if err := p.playerRepo.UpdatePlayerProfile(profile); err != nil {
			return nil, errors.New("failed to update player profile")
		}
	} else {
		// Create new profile
		profile = &entity.PlayerProfile{
			UserID:            userID,
			Level:             req.Level,
			Gender:            req.Gender,
			PreferredLocation: req.PreferredLocation,
			Latitude:          req.Latitude,
			Longitude:         req.Longitude,
			Bio:               req.Bio,
		}

		if err := p.playerRepo.CreatePlayerProfile(profile); err != nil {
			return nil, errors.New("failed to create player profile")
		}
	}

	// Get the updated profile with relationships
	updatedProfile, err := p.playerRepo.GetPlayerProfileByID(profile.ID)
	if err != nil {
		return nil, errors.New("failed to get updated profile")
	}

	return p.convertToResponse(updatedProfile), nil
}

func (p *playerUsecase) GetPlayerProfileByID(id uint64) (*entity.PlayerProfileResponse, error) {
	profile, err := p.playerRepo.GetPlayerProfileByID(id)
	if err != nil {
		return nil, errors.New("player profile not found")
	}

	return p.convertToResponse(profile), nil
}

func (p *playerUsecase) GetPlayerProfileByUserID(userID uint64) (*entity.PlayerProfileResponse, error) {
	profile, err := p.playerRepo.GetPlayerProfileByUserID(userID)
	if err != nil {
		return nil, errors.New("player profile not found")
	}

	return p.convertToResponse(profile), nil
}

func (p *playerUsecase) GetPlayerSuggestions(userID uint64, req *entity.PlayerSuggestionRequest) (*entity.PlayerSuggestionResponse, error) {
	// Set default values
	radius := 10.0 // Default 10km radius
	limit := 20    // Default 20 suggestions

	if req.Radius != nil {
		radius = *req.Radius
	}
	if req.Limit != nil {
		limit = *req.Limit
	}

	var profiles []*entity.PlayerProfile
	var err error

	if req.Latitude != nil && req.Longitude != nil {
		// Get nearby players with optional level filter
		if req.Level != nil && *req.Level != "" {
			profiles, err = p.playerRepo.GetPlayerSuggestions(userID, req.Level, req.Latitude, req.Longitude, radius, limit)
		} else {
			profiles, err = p.playerRepo.GetNearbyPlayers(userID, *req.Latitude, *req.Longitude, radius, limit)
		}
	} else if req.Level != nil && *req.Level != "" {
		// Get players by level only
		profiles, err = p.playerRepo.GetPlayersByLevel(userID, *req.Level, limit)
	} else {
		// Get all players (with limit)
		profiles, err = p.playerRepo.GetPlayerSuggestions(userID, nil, nil, nil, 0, limit)
	}

	if err != nil {
		return nil, errors.New("failed to get player suggestions")
	}

	// Convert to response format
	var responseProfiles []entity.PlayerProfileResponse
	for _, profile := range profiles {
		responseProfiles = append(responseProfiles, *p.convertToResponse(profile))
	}

	return &entity.PlayerSuggestionResponse{
		Players: responseProfiles,
		Total:   len(responseProfiles),
		Page:    1, // For future pagination
		Limit:   limit,
	}, nil
}

// convertToResponse converts PlayerProfile entity to response format
func (p *playerUsecase) convertToResponse(profile *entity.PlayerProfile) *entity.PlayerProfileResponse {
	response := &entity.PlayerProfileResponse{
		ID:                profile.ID,
		UserID:            profile.UserID,
		Level:             profile.Level,
		Gender:            profile.Gender,
		PreferredLocation: profile.PreferredLocation,
		Latitude:          profile.Latitude,
		Longitude:         profile.Longitude,
		Bio:               profile.Bio,
		CreatedAt:         profile.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         profile.UpdatedAt.Format(time.RFC3339),
	}

	// Add distance if available
	if profile.Distance != nil {
		response.Distance = profile.Distance
	}

	// Add user information if available
	if profile.User != nil {
		response.User = profile.User
	}

	return response
}
