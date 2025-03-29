package main

import (
	"log"
	"mentorship-backend/config"
	"mentorship-backend/models"
	"mentorship-backend/routes"
	"mentorship-backend/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Firebase Admin SDK
	utils.InitFirebase()

	// Initialize database
	config.InitializeDatabase()

	// Initialize Cloudinary
	if err := utils.InitCloudinary(); err != nil {
		log.Fatal("Error initializing Cloudinary:", err)
	}

	// Auto-migrate all models
	config.GetDB().AutoMigrate(
		&models.User{},
		&models.MentorDetails{},
		&models.Post{},
		&models.Follow{},
		&models.Comment{},
		&models.Tag{},
		&models.Like{},
		&models.Notification{},
	)

	// Setup Gin router in release mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r)

	// Get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
