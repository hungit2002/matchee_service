package entity

import "time"

type Venue struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID      uint64     `gorm:"not null" json:"ownerId"`
	Name         string     `gorm:"size:150;not null" json:"name"`
	Address      *string    `gorm:"size:255" json:"address,omitempty"`
	Latitude     *float64   `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude    *float64   `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	ContactPhone *string    `gorm:"size:20" json:"contactPhone,omitempty"`
	BasePrice    *float64   `gorm:"type:decimal(10,2)" json:"basePrice,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Owner       *User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Courts      []Court      `gorm:"foreignKey:VenueID" json:"courts,omitempty"`
	MatchPosts  []MatchPost  `gorm:"foreignKey:VenueID" json:"matchPosts,omitempty"`
	MatchGroups []MatchGroup `gorm:"foreignKey:VenueID" json:"matchGroups,omitempty"`
	Bookings    []Booking    `gorm:"foreignKey:VenueID" json:"bookings,omitempty"`
}
