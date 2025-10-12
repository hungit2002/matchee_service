package entity

import "time"

type User struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string     `gorm:"size:100;not null" json:"fullName"`
	Phone        string     `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Email        *string    `gorm:"uniqueIndex;size:100" json:"email,omitempty"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	PlayerProfile     *PlayerProfile     `gorm:"foreignKey:UserID" json:"playerProfile,omitempty"`
	Venues            []Venue            `gorm:"foreignKey:OwnerID" json:"venues,omitempty"`
	MatchPosts        []MatchPost        `gorm:"foreignKey:UserID" json:"matchPosts,omitempty"`
	MatchGroupMembers []MatchGroupMember `gorm:"foreignKey:UserID" json:"matchGroupMembers,omitempty"`
	Notifications     []Notification     `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	Payments          []Payment          `gorm:"foreignKey:PayerID" json:"payments,omitempty"`
	FeedbacksFrom     []Feedback         `gorm:"foreignKey:FromUserID" json:"feedbacksFrom,omitempty"`
	FeedbacksTo       []Feedback         `gorm:"foreignKey:ToUserID" json:"feedbacksTo,omitempty"`
	UserRoles         []UserRole         `gorm:"foreignKey:UserID" json:"userRoles,omitempty"`
}
