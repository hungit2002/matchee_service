package entity

import "time"

type MatchGroup struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchPostID   *uint64    `json:"matchPostId,omitempty"`
	VenueID       *uint64    `json:"venueId,omitempty"`
	CourtID       *uint64    `json:"courtId,omitempty"`
	ScheduledTime *time.Time `json:"scheduledTime,omitempty"`
	TotalPrice    *float64   `gorm:"type:decimal(10,2)" json:"totalPrice,omitempty"`
	Status        string     `gorm:"type:enum('pending','confirmed','completed','cancelled');default:'pending'" json:"status"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	MatchPost *MatchPost         `gorm:"foreignKey:MatchPostID" json:"matchPost,omitempty"`
	Venue     *Venue             `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	Court     *Court             `gorm:"foreignKey:CourtID" json:"court,omitempty"`
	Members   []MatchGroupMember `gorm:"foreignKey:MatchGroupID" json:"members,omitempty"`
	Bookings  []Booking          `gorm:"foreignKey:MatchGroupID" json:"bookings,omitempty"`
}
