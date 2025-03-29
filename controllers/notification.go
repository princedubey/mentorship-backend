package controllers

import (
	"mentorship-backend/config"
	"mentorship-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

// CreateNotification creates a new notification
func (nc *NotificationController) CreateNotification(notification *models.Notification) error {
	return config.GetDB().Create(notification).Error
}

// GetNotifications gets notifications for a user
func (nc *NotificationController) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var notifications []models.Notification
	if err := config.GetDB().
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("User").
		Preload("Actor").
		Preload("Post").
		Find(&notifications).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(200, notifications)
}

// MarkAsRead marks a notification as read
func (nc *NotificationController) MarkAsRead(c *gin.Context) {
	notifID := c.Param("id")
	if notifID == "" {
		c.JSON(400, gin.H{"error": "Notification ID is required"})
		return
	}

	var notification models.Notification
	if err := config.GetDB().First(&notification, notifID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Notification not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists || notification.UserID != userID.(uuid.UUID) {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.GetDB().Model(&notification).Update("is_read", true).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(200, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead marks all notifications as read
func (nc *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	if err := config.GetDB().Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(200, gin.H{"message": "All notifications marked as read"})
}
