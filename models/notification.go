package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"` // User receiving the notification
	User      User      `gorm:"foreignKey:UserID"`
	ActorID   uuid.UUID `gorm:"type:uuid;not null"` // User who performed the action
	Actor     User      `gorm:"foreignKey:ActorID"`
	PostID    *uuid.UUID `gorm:"type:uuid"` // Optional: if notification is related to a post
	Post      Post      `gorm:"foreignKey:PostID"`
	Type      string    `gorm:"type:varchar(50);not null"` // 'follow', 'like', 'comment', etc.
	Message   string    `gorm:"type:text;not null"` // Human-readable message
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	NotificationTypeFollow = "follow"
	NotificationTypeLike   = "like"
	NotificationTypeComment = "comment"
)
