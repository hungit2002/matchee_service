package entity

import "time"

type Booking struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchGroupID *uint64    `json:"matchGroupId,omitempty"`
	VenueID      *uint64    `json:"venueId,omitempty"`
	CourtID      *uint64    `json:"courtId,omitempty"`
	StartTime    *time.Time `json:"startTime,omitempty"`
	EndTime      *time.Time `json:"endTime,omitempty"`
	TotalPrice   *float64   `gorm:"type:decimal(10,2)" json:"totalPrice,omitempty"`
	Status       string     `gorm:"type:enum('reserved','paid','cancelled','completed');default:'reserved'" json:"status"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	MatchGroup *MatchGroup `gorm:"foreignKey:MatchGroupID" json:"matchGroup,omitempty"`
	Venue      *Venue      `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	Court      *Court      `gorm:"foreignKey:CourtID" json:"court,omitempty"`
	Payments   []Payment   `gorm:"foreignKey:BookingID" json:"payments,omitempty"`
	Feedbacks  []Feedback  `gorm:"foreignKey:BookingID" json:"feedbacks,omitempty"`
}
