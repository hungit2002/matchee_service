package entity

import "time"

type Court struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VenueID     uint64     `gorm:"not null" json:"venueId"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CourtType   string     `gorm:"type:enum('indoor','outdoor');default:'outdoor'" json:"courtType"`
	Surface     string     `gorm:"type:enum('hard','clay','grass','synthetic');default:'hard'" json:"surface"`
	IsActive    bool       `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Venue             *Venue             `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	AvailabilitySlots []AvailabilitySlot `gorm:"foreignKey:CourtID" json:"availabilitySlots,omitempty"`
	MatchGroups       []MatchGroup       `gorm:"foreignKey:CourtID" json:"matchGroups,omitempty"`
	Bookings          []Booking          `gorm:"foreignKey:CourtID" json:"bookings,omitempty"`
}
