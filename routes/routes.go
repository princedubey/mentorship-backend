package routes

import (
	"mentorship-backend/controllers"
	"mentorship-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	userController := controllers.NewUserController()
	mentorController := controllers.NewMentorController()
	tagController := controllers.NewTagController()
	postController := controllers.NewPostController()
	commentController := controllers.NewCommentController()
	followController := controllers.NewFollowController()
	authController := controllers.NewAuthController()
	likeController := controllers.NewLikeController()
	notificationController := controllers.NewNotificationController()

	// Public routes
	public := r.Group("/api")
	{
		// Auth routes
		public.POST("/auth/firebase", authController.AuthenticateWithFirebase)
		public.POST("/auth/refresh", authController.RefreshToken)
		public.POST("/register", userController.RegisterUser)
		public.POST("/login", userController.LoginUser)

		// Public user routes
		public.GET("/users/:id", userController.GetUserProfile)
		public.GET("/users/:id/followers", followController.GetFollowers)
		public.GET("/users/:id/following", followController.GetFollowing)
		
		// Public mentor routes
		public.GET("/mentors", mentorController.ListMentors)
		public.GET("/mentors/:id", mentorController.GetMentorProfile)

		// Public tag routes
		public.GET("/tags", tagController.ListTags)

		// Public post routes
		public.GET("/posts", postController.ListPosts)
		public.GET("/posts/:id", postController.GetPost)
		public.GET("/posts/:id/comments", commentController.GetComments)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Protected user routes
		protected.GET("/profile", userController.GetProfile)
		protected.PUT("/profile", userController.UpdateProfile)
		protected.PUT("/profile/password", userController.ChangePassword)
		protected.GET("/profile/saved-posts", userController.GetSavedPosts)
		protected.POST("/profile/deactivate", userController.DeactivateAccount)
		
		// Follow routes
		protected.POST("/users/:id/follow", followController.FollowUser)
		protected.DELETE("/users/:id/follow", followController.UnfollowUser)
		
		// Protected mentor routes
		protected.POST("/mentor/profile", mentorController.CreateMentorProfile)
		protected.PUT("/mentor/availability", mentorController.UpdateAvailability)

		// Protected tag routes
		protected.POST("/tags", tagController.CreateTag)
		protected.POST("/user/tags", tagController.AddTagsToUser)
		protected.POST("/mentor/tags", tagController.AddTagsToMentor)

		// Protected notification routes
		protected.GET("/notifications", notificationController.GetNotifications)
		protected.PUT("/notifications/:id/read", notificationController.MarkAsRead)
		protected.PUT("/notifications/read-all", notificationController.MarkAllAsRead)

		// Protected post routes
		protected.POST("/posts", postController.CreatePost)
		protected.POST("/posts/:id/share", postController.SharePost)
		protected.POST("/posts/:id/save", postController.SavePost)
        protected.GET("/posts/:id/analytics", postController.GetPostAnalytics)
		protected.POST("/posts/:id/tags", postController.AddTagsToPost)
		protected.DELETE("/posts/:id", postController.DeletePost)
		protected.POST("/posts/:id/like", likeController.LikePost)
		protected.DELETE("/posts/:id/like", likeController.UnlikePost)
		protected.GET("/posts/:id/likes", likeController.GetPostLikes)

		// Protected comment routes
		protected.POST("/posts/:id/comments", commentController.CreateComment)
		protected.POST("/comments/:id/reply", commentController.ReplyToComment)
	}
}
