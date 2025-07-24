package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/supabase"
	"github.com/holycann/cultour-backend/internal/users/handlers"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/internal/users/services"
)

func RegisterUserRoutes(
	r *gin.Engine,
	supabaseAuth *supabase.SupabaseAuth,
) {
	userRepository := repositories.NewUserRepository(supabaseAuth.GetClient())
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	user := r.Group("/users")
	{
		user.POST("/", userHandler.CreateUser)
		user.GET("/", userHandler.ListUsers)
		user.GET("/search", userHandler.SearchUser)
		user.PUT("/:id", userHandler.UpdateUser)
		user.DELETE("/:id", userHandler.DeleteUser)
	}
}

func RegisterUserProfileRoutes(
	r *gin.Engine,
	supabaseClient *supabase.SupabaseClient,
	table string,
) {
	userProfileRepository := repositories.NewUserProfileRepository(supabaseClient.GetClient(), table)
	userProfileService := services.NewUserProfileService(userProfileRepository)
	userProfileHandler := handlers.NewUserProfileHandler(userProfileService)

	profile := r.Group("/profile")
	{
		profile.POST("/", userProfileHandler.CreateUserProfile)
		profile.GET("/", userProfileHandler.ListUsersProfile)
		profile.GET("/search", userProfileHandler.ListUsersProfile)
		profile.PUT("/:id", userProfileHandler.UpdateUserProfile)
		profile.DELETE("/:id", userProfileHandler.DeleteUserProfile)
	}
}
