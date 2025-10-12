package repository

import (
	"math"

	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	// Player profile methods
	CreatePlayerProfile(profile *entity.PlayerProfile) error
	GetPlayerProfileByID(id uint64) (*entity.PlayerProfile, error)
	GetPlayerProfileByUserID(userID uint64) (*entity.PlayerProfile, error)
	UpdatePlayerProfile(profile *entity.PlayerProfile) error
	DeletePlayerProfile(id uint64) error

	// Player suggestion methods
	GetPlayerSuggestions(userID uint64, level *string, latitude, longitude *float64, radius float64, limit int) ([]*entity.PlayerProfile, error)
	GetNearbyPlayers(userID uint64, latitude, longitude float64, radius float64, limit int) ([]*entity.PlayerProfile, error)
	GetPlayersByLevel(userID uint64, level string, limit int) ([]*entity.PlayerProfile, error)
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

// Player profile methods
func (r *playerRepository) CreatePlayerProfile(profile *entity.PlayerProfile) error {
	return r.db.Create(profile).Error
}

func (r *playerRepository) GetPlayerProfileByID(id uint64) (*entity.PlayerProfile, error) {
	var profile entity.PlayerProfile
	err := r.db.Preload("User.UserRoles.Role").Where("id = ?", id).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *playerRepository) GetPlayerProfileByUserID(userID uint64) (*entity.PlayerProfile, error) {
	var profile entity.PlayerProfile
	err := r.db.Preload("User.UserRoles.Role").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *playerRepository) UpdatePlayerProfile(profile *entity.PlayerProfile) error {
	return r.db.Save(profile).Error
}

func (r *playerRepository) DeletePlayerProfile(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&entity.PlayerProfile{}).Error
}

// Player suggestion methods
func (r *playerRepository) GetPlayerSuggestions(userID uint64, level *string, latitude, longitude *float64, radius float64, limit int) ([]*entity.PlayerProfile, error) {
	var profiles []*entity.PlayerProfile

	query := r.db.Preload("User.UserRoles.Role").
		Where("user_id != ?", userID) // Exclude current user

	// Filter by level if provided
	if level != nil && *level != "" {
		query = query.Where("level = ?", *level)
	}

	// Filter by location if provided
	if latitude != nil && longitude != nil {
		// Use Haversine formula for distance calculation
		// This is a simplified version - in production, consider using PostGIS or similar
		query = query.Where("latitude IS NOT NULL AND longitude IS NOT NULL")
	}

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	// Filter by radius if location is provided
	if latitude != nil && longitude != nil && radius > 0 {
		var filteredProfiles []*entity.PlayerProfile
		for _, profile := range profiles {
			if profile.Latitude != nil && profile.Longitude != nil {
				distance := r.calculateDistance(*latitude, *longitude, *profile.Latitude, *profile.Longitude)
				if distance <= radius {
					profile.Distance = &distance // Add distance to profile
					filteredProfiles = append(filteredProfiles, profile)
				}
			}
		}
		return filteredProfiles, nil
	}

	return profiles, nil
}

func (r *playerRepository) GetNearbyPlayers(userID uint64, latitude, longitude float64, radius float64, limit int) ([]*entity.PlayerProfile, error) {
	var profiles []*entity.PlayerProfile

	query := r.db.Preload("User.UserRoles.Role").
		Where("user_id != ?", userID).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	// Filter by radius
	var filteredProfiles []*entity.PlayerProfile
	for _, profile := range profiles {
		if profile.Latitude != nil && profile.Longitude != nil {
			distance := r.calculateDistance(latitude, longitude, *profile.Latitude, *profile.Longitude)
			if distance <= radius {
				profile.Distance = &distance
				filteredProfiles = append(filteredProfiles, profile)
			}
		}
	}

	return filteredProfiles, nil
}

func (r *playerRepository) GetPlayersByLevel(userID uint64, level string, limit int) ([]*entity.PlayerProfile, error) {
	var profiles []*entity.PlayerProfile

	query := r.db.Preload("User.UserRoles.Role").
		Where("user_id != ?", userID).
		Where("level = ?", level)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&profiles).Error
	return profiles, err
}

// calculateDistance calculates the distance between two points using Haversine formula
func (r *playerRepository) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
