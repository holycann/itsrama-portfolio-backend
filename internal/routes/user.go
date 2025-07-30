package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/supabase"
	"github.com/holycann/cultour-backend/internal/users/handlers"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/internal/users/services"
)

func RegisterUserRoutes(
	r *gin.Engine,
	supabaseAuth *supabase.SupabaseAuth,
	routerMiddleware *middleware.Middleware,
) {
	userRepository := repositories.NewUserRepository(supabaseAuth.GetClient())
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	user := r.Group("/users")
	{
		user.POST("",
			routerMiddleware.VerifyJWT(),
			userHandler.CreateUser)
		user.GET("",
			routerMiddleware.VerifyJWT(),
			userHandler.ListUsers)
		user.GET("/search",
			routerMiddleware.VerifyJWT(),
			userHandler.SearchUsers)
		user.GET("/:id",
			routerMiddleware.VerifyJWT(),
			userHandler.GetUserByID) // detail by id
		user.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			userHandler.UpdateUser)
		user.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			userHandler.DeleteUser)
	}
}

func RegisterUserProfileRoutes(
	r *gin.Engine,
	supabaseClient *supabase.SupabaseClient,
	supabaseAuth *supabase.SupabaseAuth,
	supabaseStorage *supabase.SupabaseStorage,
	routerMiddleware *middleware.Middleware,
) {
	userRepository := repositories.NewUserRepository(supabaseAuth.GetClient())
	userProfileRepository := repositories.NewUserProfileRepository(supabaseClient.GetClient())
	userProfileService := services.NewUserProfileService(userProfileRepository, userRepository, supabaseStorage)
	userProfileHandler := handlers.NewUserProfileHandler(userProfileService)

	profile := r.Group("/profile")
	{
		profile.POST("",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.CreateUserProfile)
		profile.GET("",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.ListUsersProfile)
		profile.GET("/search",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.SearchUserProfile)
		profile.GET("/me",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.GetAuthenticatedUserProfile)
		profile.GET("/:id",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.GetUserProfileById)
		profile.PUT("/:id",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.UpdateUserProfile)
		profile.DELETE("/:id",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.DeleteUserProfile)
	}
}

func RegisterUserBadgeRoutes(
	r *gin.Engine,
	supabaseClient *supabase.SupabaseClient,
	routerMiddleware *middleware.Middleware,
	appLogger *logger.Logger,
) {
	userBadgeRepository := repositories.NewUserBadgeRepository(supabaseClient.GetClient())
	userBadgeService := services.NewUserBadgeService(userBadgeRepository)
	userBadgeHandler := handlers.NewUserBadgeHandler(userBadgeService, appLogger)

	badges := r.Group("/users/badges")
	{
		badges.POST("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", "warlok"),
			userBadgeHandler.AssignBadge,
		)
		badges.GET("",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.GetUserBadges,
		)
		badges.GET("/search",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.GetUserBadges, // Reuse GetUserBadges with query parameter
		)
		badges.DELETE("",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", "warlok"),
			userBadgeHandler.RemoveBadge,
		)
		badges.GET("/count",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.CountUserBadges,
		)
	}
}
