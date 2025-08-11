package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/users/handlers"
)

func RegisterUserRoutes(
	r *gin.Engine,
	userHandler *handlers.UserHandler,
	routerMiddleware *middleware.Middleware,
) {
	user := r.Group("/users")
	{
		// Create a new user account
		user.POST("", userHandler.CreateUser)

		// List users (requires authentication and admin privileges)
		user.GET("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userHandler.ListUsers)

		// Search users (requires authentication and admin privileges)
		user.GET("/search",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userHandler.SearchUsers)

		// Get user by ID (requires authentication)
		user.GET("/:id",
			routerMiddleware.VerifyJWT(),
			userHandler.GetUserByID)

		// Update user (requires authentication and admin privileges)
		user.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userHandler.UpdateUser)

		// Delete user (requires authentication and admin privileges)
		user.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userHandler.DeleteUser)
	}
}

func RegisterUserProfileRoutes(
	r *gin.Engine,
	userProfileHandler *handlers.UserProfileHandler,
	routerMiddleware *middleware.Middleware,
) {
	profile := r.Group("/users/profiles")
	{
		// Create a new user profile (requires authentication)
		profile.POST("",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.CreateUserProfile)

		// List user profiles (requires authentication and admin privileges)
		profile.GET("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userProfileHandler.ListUsersProfile)

		// Search user profiles (requires authentication and admin privileges)
		profile.GET("/search",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userProfileHandler.SearchUserProfile)

		// Get authenticated user's profile (requires authentication)
		profile.GET("/me",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.GetAuthenticatedUserProfile)

		// Get user profile by ID (requires authentication)
		profile.GET("/:id",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.GetUserProfileById)

		// Update user profile (requires authentication)
		profile.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.UpdateUserProfile)

		// Update user profile avatar (requires authentication)
		profile.PUT("/:id/avatar",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.UpdateUserAvatar)

		// Verify user identity (requires authentication)
		profile.PUT("/:id/verify",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.VerifyIdentity)

		// Delete user profile (requires authentication and admin privileges)
		profile.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userProfileHandler.DeleteUserProfile)
	}
}

func RegisterUserBadgeRoutes(
	r *gin.Engine,
	userBadgeHandler *handlers.UserBadgeHandler,
	routerMiddleware *middleware.Middleware,
) {
	badges := r.Group("/users/badges")
	{
		// Assign a badge to a user (admin or specific role only)
		badges.POST("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userBadgeHandler.AssignBadge)

		// Get user badges (requires authentication)
		badges.GET("",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.GetUserBadges)

		// Get user badges by specific user (requires authentication)
		badges.GET("/:user_id",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.GetUserBadgesByUser)

		// Count user badges (requires authentication)
		badges.GET("/count",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.CountUserBadges)

		// Remove a badge from a user (admin or specific role only)
		badges.DELETE("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", ""),
			userBadgeHandler.RemoveBadge)
	}
}
