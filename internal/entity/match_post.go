package entity

import "time"

type MatchPost struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64     `gorm:"not null" json:"userId"`
	DesiredLevel   *string    `gorm:"type:enum('beginner','average','good','pro')" json:"desiredLevel,omitempty"`
	Location       *string    `gorm:"size:255" json:"location,omitempty"`
	Latitude       *float64   `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude      *float64   `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	VenueID        *uint64    `json:"venueId,omitempty"`
	DesiredTime    *time.Time `json:"desiredTime,omitempty"`
	PricePerPerson *float64   `gorm:"type:decimal(10,2)" json:"pricePerPerson,omitempty"`
	MaxPlayers     int        `gorm:"default:4" json:"maxPlayers"`
	Note           *string    `gorm:"type:text" json:"note,omitempty"`
	Status         string     `gorm:"type:enum('open','matched','cancelled','done');default:'open'" json:"status"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt      *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Venue       *Venue       `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	MatchGroups []MatchGroup `gorm:"foreignKey:MatchPostID" json:"matchGroups,omitempty"`
}
