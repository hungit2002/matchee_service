package entity

import "time"

type Court struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VenueID     uint64     `gorm:"not null" json:"venueId"`
	Name        *string    `gorm:"size:100" json:"name,omitempty"`
	CourtNumber *int       `json:"courtNumber,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Venue             *Venue             `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	AvailabilitySlots []AvailabilitySlot `gorm:"foreignKey:CourtID" json:"availabilitySlots,omitempty"`
	MatchGroups       []MatchGroup       `gorm:"foreignKey:CourtID" json:"matchGroups,omitempty"`
	Bookings          []Booking          `gorm:"foreignKey:CourtID" json:"bookings,omitempty"`
}
