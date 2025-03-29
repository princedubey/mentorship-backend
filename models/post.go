package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostAnalytics struct {
	Views        int `json:"views" gorm:"default:0"`
	Shares       int `json:"shares" gorm:"default:0"`
	SavedCount   int `json:"savedCount" gorm:"default:0"`
	CommentCount int `json:"commentCount" gorm:"default:0"`
	Likes        int `json:"likes" gorm:"default:0"`
	Comments     int `json:"comments" gorm:"default:0"`
}

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Content   string    `gorm:"type:text"`
	MediaURLs []string  `gorm:"type:text[]"`
	IsPrivate bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Analytics
	Analytics PostAnalytics `gorm:"embedded"`

	// Sharing
	OriginalPostID *uuid.UUID `gorm:"type:uuid"` // If this is a shared post
	OriginalPost   *Post      `gorm:"foreignKey:OriginalPostID"`
	SharedPosts    []Post     `gorm:"foreignKey:OriginalPostID"` // Posts that shared this post

	// Comments
	Comments []Comment `gorm:"foreignKey:PostID"`

	// Tags for categorizing posts (topics, skills, etc.)
	Tags []Tag `gorm:"many2many:post_tags;"`

	// Users who saved this post
	SavedBy []User `gorm:"many2many:user_saved_posts;"`

	// Likes
	Likes []Like `gorm:"foreignKey:PostID"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// GetLikesCount returns the number of likes for a post
func (p *Post) GetLikesCount(db *gorm.DB) (int64, error) {
	var count int64
	if err := db.Model(&Like{}).Where("post_id = ?", p.ID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentsCount returns the number of comments for a post
func (p *Post) GetCommentsCount(db *gorm.DB) (int64, error) {
	var count int64
	if err := db.Model(&Comment{}).Where("post_id = ?", p.ID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
