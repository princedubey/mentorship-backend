package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `gorm:"uniqueIndex;not null"`
	Category  string    `gorm:"type:varchar(50)"` // e.g., "skill", "interest", "topic"
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Many-to-Many relationships
	Users   []User          `gorm:"many2many:user_tags;"`
	Mentors []MentorDetails `gorm:"many2many:mentor_tags;"`
	Posts   []Post          `gorm:"many2many:post_tags;"`
}
