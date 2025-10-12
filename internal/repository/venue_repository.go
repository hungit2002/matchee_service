package repository

import (
	"math"

	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

// VenueRepository defines the interface for venue operations
type VenueRepository interface {
	CreateVenue(venue *entity.Venue) error
	GetVenueByID(id uint64) (*entity.Venue, error)
	UpdateVenue(id uint64, venue *entity.Venue) error
	DeleteVenue(id uint64) error
	GetVenuesByOwner(ownerID uint64, page, limit int) ([]*entity.Venue, int64, error)
	GetVenuesByLocation(latitude, longitude, radius float64, minPrice, maxPrice *float64, isActive *bool, page, limit int) ([]*entity.Venue, int64, error)
	GetAllVenues(page, limit int) ([]*entity.Venue, int64, error)
}

type venueRepository struct {
	db *gorm.DB
}

// NewVenueRepository creates a new instance of VenueRepository
func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) CreateVenue(venue *entity.Venue) error {
	return r.db.Create(venue).Error
}

func (r *venueRepository) GetVenueByID(id uint64) (*entity.Venue, error) {
	var venue entity.Venue
	err := r.db.Preload("Owner").Preload("Courts").Where("id = ?", id).First(&venue).Error
	if err != nil {
		return nil, err
	}
	return &venue, nil
}

func (r *venueRepository) UpdateVenue(id uint64, venue *entity.Venue) error {
	return r.db.Model(&entity.Venue{}).Where("id = ?", id).Updates(venue).Error
}

func (r *venueRepository) DeleteVenue(id uint64) error {
	return r.db.Delete(&entity.Venue{}, id).Error
}

func (r *venueRepository) GetVenuesByOwner(ownerID uint64, page, limit int) ([]*entity.Venue, int64, error) {
	var venues []*entity.Venue
	var total int64

	query := r.db.Model(&entity.Venue{}).Where("owner_id = ?", ownerID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Owner").Preload("Courts").
		Offset(offset).Limit(limit).
		Find(&venues).Error

	return venues, total, err
}

func (r *venueRepository) GetVenuesByLocation(latitude, longitude, radius float64, minPrice, maxPrice *float64, isActive *bool, page, limit int) ([]*entity.Venue, int64, error) {
	var venues []*entity.Venue
	var total int64

	query := r.db.Model(&entity.Venue{})

	// Filter by location if provided
	if latitude != 0 && longitude != 0 && radius > 0 {
		query = query.Where("latitude IS NOT NULL AND longitude IS NOT NULL")
	}

	// Filter by price range
	if minPrice != nil {
		query = query.Where("price_per_hour >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price_per_hour <= ?", *maxPrice)
	}

	// Filter by active status
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Owner").Preload("Courts").
		Offset(offset).Limit(limit).
		Find(&venues).Error

	if err != nil {
		return nil, 0, err
	}

	// Filter by radius if location is provided
	if latitude != 0 && longitude != 0 && radius > 0 {
		var filteredVenues []*entity.Venue
		for _, venue := range venues {
			if venue.Latitude != 0 && venue.Longitude != 0 {
				distance := r.calculateDistance(latitude, longitude, venue.Latitude, venue.Longitude)
				if distance <= radius {
					filteredVenues = append(filteredVenues, venue)
				}
			}
		}
		venues = filteredVenues
		total = int64(len(filteredVenues))
	}

	return venues, total, err
}

func (r *venueRepository) GetAllVenues(page, limit int) ([]*entity.Venue, int64, error) {
	var venues []*entity.Venue
	var total int64

	query := r.db.Model(&entity.Venue{})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Owner").Preload("Courts").
		Offset(offset).Limit(limit).
		Find(&venues).Error

	return venues, total, err
}

func (r *venueRepository) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
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
