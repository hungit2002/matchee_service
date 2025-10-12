package entity

import "time"

type Venue struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID      uint64     `gorm:"not null" json:"ownerId"`
	Name         string     `gorm:"size:150;not null" json:"name"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	Address      string     `gorm:"size:255;not null" json:"address"`
	Latitude     float64    `gorm:"type:decimal(10,7);not null" json:"latitude"`
	Longitude    float64    `gorm:"type:decimal(10,7);not null" json:"longitude"`
	Phone        string     `gorm:"size:20;not null" json:"phone"`
	Email        *string    `gorm:"size:100" json:"email,omitempty"`
	Website      *string    `gorm:"size:255" json:"website,omitempty"`
	PricePerHour float64    `gorm:"type:decimal(10,2);not null" json:"pricePerHour"`
	IsActive     bool       `gorm:"default:true" json:"isActive"`
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
