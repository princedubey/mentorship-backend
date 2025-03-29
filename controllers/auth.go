package controllers

import (
	"mentorship-backend/config"
	"mentorship-backend/models"
	"mentorship-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type FirebaseAuthRequest struct {
	FirebaseToken string `json:"firebaseToken" binding:"required"`
}

// AuthenticateWithFirebase handles Firebase authentication and returns JWT tokens
func (ac *AuthController) AuthenticateWithFirebase(c *gin.Context) {
	var req FirebaseAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify Firebase token
	token, err := utils.VerifyFirebaseToken(c.Request.Context(), req.FirebaseToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Firebase token"})
		return
	}

	// Check if user exists in our database
	var user models.User
	result := config.GetDB().Where("firebase_uid = ?", token.UID).First(&user)

	if result.Error != nil {
		// Create new user if not exists
		user = models.User{
			FirebaseUID: token.UID,
			Name: func() string {
				if name, ok := token.Claims["name"].(string); ok {
					return name
				}
				// Use provider name as fallback
				if provider, ok := token.Claims["firebase_sign_in_provider"].(string); ok {
					return "User (" + provider + ")"
				}
				return "User"
			}(),
			Role: models.RoleUser,
		}

		// If Firebase user has a profile picture
		if picture, ok := token.Claims["picture"].(string); ok {
			user.AvatarURL = picture
		}

		// Set email if available
		if email, ok := token.Claims["email"].(string); ok {
			user.Email = email
		}

		// Set phone number if available
		if phone, ok := token.Claims["phone_number"].(string); ok {
			user.PhoneNumber = phone
		}

		if err := config.GetDB().Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Check if user is a mentor
	var mentorDetails models.MentorDetails
	isMentor := config.GetDB().Where("user_id = ?", user.ID).First(&mentorDetails).Error == nil

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":        user,
		"isMentor":    isMentor,
	})
}

// RefreshToken handles token refresh requests
func (ac *AuthController) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	accessToken, newRefreshToken, err := utils.RefreshTokenPair(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}
