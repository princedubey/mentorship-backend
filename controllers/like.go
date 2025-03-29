package controllers

import (
	"fmt"
	"mentorship-backend/config"
	"mentorship-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeController struct{}

func NewLikeController() *LikeController {
	return &LikeController{}
}

// LikePost handles liking a post
func (lc *LikeController) LikePost(c *gin.Context) {
	var like models.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if user has already liked this post
	var existingLike models.Like
	if err := config.GetDB().Where("post_id = ? AND user_id = ?", like.PostID, like.UserID).First(&existingLike).Error; err == nil {
		c.JSON(400, gin.H{"error": "Post already liked"})
		return
	}

	tx := config.GetDB().Begin()
	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to like post"})
		return
	}

	// Update post likes count
	if err := tx.Model(&models.Post{}).Where("id = ?", like.PostID).Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to update likes count"})
		return
	}

	// Get post and user details
	var post models.Post
	if err := tx.First(&post, like.PostID).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to fetch post"})
		return
	}

	var user models.User
	if err := tx.First(&user, like.UserID).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to fetch user"})
		return
	}

	notificationController := NewNotificationController()
	// Create notification for the post owner
	postOwnerID := post.UserID
	if postOwnerID != like.UserID {
		notification := &models.Notification{
			UserID:   postOwnerID,
			ActorID:  like.UserID,
			PostID:   &post.ID,
			Type:     models.NotificationTypeLike,
			Message:  fmt.Sprintf("%s liked your post", user.Name),
		}

		if err := notificationController.CreateNotification(notification); err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to create notification"})
			return
		}
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "Post liked successfully"})
}

// UnlikePost handles unliking a post
func (lc *LikeController) UnlikePost(c *gin.Context) {
	postId := c.Param("id")
	userId := c.GetHeader("User-Id") // Get from auth middleware

	if postId == "" || userId == "" {
		c.JSON(400, gin.H{"error": "Post ID and User ID are required"})
		return
	}

	// Convert IDs to UUID
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find and delete the like
	if err := config.GetDB().Where("post_id = ? AND user_id = ?", postUUID, userUUID).Delete(&models.Like{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to unlike post"})
		return
	}

	// Update post likes count
	if err := config.GetDB().Model(&models.Post{}).Where("id = ?", postUUID).Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update likes count"})
		return
	}

	c.JSON(200, gin.H{"message": "Post unliked successfully"})
}

// GetPostLikes gets all users who liked a post
func (lc *LikeController) GetPostLikes(c *gin.Context) {
	postId := c.Param("id")
	if postId == "" {
		c.JSON(400, gin.H{"error": "Post ID is required"})
		return
	}

	// Convert ID to UUID
	postUUID, err := uuid.Parse(postId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	var likes []models.Like
	if err := config.GetDB().Preload("User").Where("post_id = ?", postUUID).Find(&likes).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch likes"})
		return
	}

	c.JSON(200, likes)
}
