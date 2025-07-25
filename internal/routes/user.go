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
		user.POST("/",
			routerMiddleware.VerifyJWT(),
			userHandler.CreateUser)
		user.GET("/",
			routerMiddleware.VerifyJWT(),
			userHandler.ListUsers)
		user.GET("/search",
			routerMiddleware.VerifyJWT(),
			userHandler.SearchUser)
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
	table string,
	routerMiddleware *middleware.Middleware,
) {
	userProfileRepository := repositories.NewUserProfileRepository(supabaseClient.GetClient(), table)
	userProfileService := services.NewUserProfileService(userProfileRepository)
	userProfileHandler := handlers.NewUserProfileHandler(userProfileService)

	profile := r.Group("/profile")
	{
		profile.POST("/",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.CreateUserProfile)
		profile.GET("/",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.ListUsersProfile)
		profile.GET("/search",
			routerMiddleware.VerifyJWT(),
			userProfileHandler.ListUsersProfile)
	}
}

func RegisterUserBadgeRoutes(
	r *gin.Engine,
	supabaseClient *supabase.SupabaseClient,
	routerMiddleware *middleware.Middleware,
	appLogger *logger.Logger,
) {
	userBadgeRepository := repositories.NewUserBadgeRepository(supabaseClient.GetClient(), *repositories.DefaultUserBadgeConfig())
	userBadgeService := services.NewUserBadgeService(userBadgeRepository)
	userBadgeHandler := handlers.NewUserBadgeHandler(userBadgeService, appLogger)

	badges := r.Group("/users/badges")
	{
		badges.POST("/",
			routerMiddleware.VerifyJWT(),
			routerMiddleware.RequireRoleOrBadge("admin", "warlok"),
			userBadgeHandler.AssignBadge,
		)
		badges.GET("/",
			routerMiddleware.VerifyJWT(),
			userBadgeHandler.GetUserBadges,
		)
		badges.DELETE("/",
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
