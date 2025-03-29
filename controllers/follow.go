package controllers

import (
	"fmt"
	"mentorship-backend/config"
	"mentorship-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowController struct{}

func NewFollowController() *FollowController {
	return &FollowController{}
}

// FollowUser handles following a user/mentor
func (fc *FollowController) FollowUser(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followingID := c.Param("id")
	followingUUID, err := uuid.Parse(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	var followingUser models.User
	if err := config.GetDB().First(&followingUser, "id = ?", followingUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to follow not found"})
		return
	}

	// Prevent self-following
	if followerID.(uuid.UUID) == followingUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	// Check if already following
	var existingFollow models.Follow
	result := config.GetDB().Where("follower_id = ? AND following_id = ?", followerID, followingUUID).First(&existingFollow)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already following this user"})
		return
	}

	follow := models.Follow{
		FollowerID:  followerID.(uuid.UUID),
		FollowingID: followingUUID,
	}

	tx := config.GetDB().Begin()
	if err := tx.Create(&follow).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	// Create notification for the followed user
	notification := &models.Notification{
		UserID:   followingUser.ID,
		ActorID:  followerID.(uuid.UUID),
		Type:     models.NotificationTypeFollow,
		Message:  fmt.Sprintf("%s started following you", followingUser.Name),
	}

	notificationController := NewNotificationController()
	if err := notificationController.CreateNotification(notification); err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to create notification"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser handles unfollowing a user/mentor
func (fc *FollowController) UnfollowUser(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followingID := c.Param("id")
	followingUUID, err := uuid.Parse(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result := config.GetDB().Where("follower_id = ? AND following_id = ?", followerID, followingUUID).Delete(&models.Follow{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not following this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

// GetFollowers gets all followers of a user
func (fc *FollowController) GetFollowers(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var follows []models.Follow
	if err := config.GetDB().Where("following_id = ?", userUUID).
		Preload("Follower").
		Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	// Extract followers from follows
	followers := make([]gin.H, len(follows))
	for i, follow := range follows {
		followers[i] = gin.H{
			"id":        follow.Follower.ID,
			"name":      follow.Follower.Name,
			"avatarURL": follow.Follower.AvatarURL,
			"followedAt": follow.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, followers)
}

// GetFollowing gets all users that a user is following
func (fc *FollowController) GetFollowing(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var follows []models.Follow
	if err := config.GetDB().Where("follower_id = ?", userUUID).
		Preload("Following").
		Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch following"})
		return
	}

	// Extract following users from follows
	following := make([]gin.H, len(follows))
	for i, follow := range follows {
		following[i] = gin.H{
			"id":        follow.Following.ID,
			"name":      follow.Following.Name,
			"avatarURL": follow.Following.AvatarURL,
			"followedAt": follow.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, following)
}
