package routes

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/configs"
	eventRepo "github.com/holycann/cultour-backend/internal/cultural/repositories"
	eventSvc "github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/gemini"
	"github.com/holycann/cultour-backend/internal/logger"
	placeRepo "github.com/holycann/cultour-backend/internal/place/repositories"
	placeSvc "github.com/holycann/cultour-backend/internal/place/services"
	sup "github.com/holycann/cultour-backend/internal/supabase"
	userRepo "github.com/holycann/cultour-backend/internal/users/repositories"
	userSvc "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/supabase-go"
)

func SetupGeminiRoutes(router *gin.Engine, cfg *configs.Config, supabaseClient *supabase.Client, supabaseAuth auth.Client, supabaseStorage *sup.SupabaseStorage, appLogger *logger.Logger) {
	// Initialize repositories
	eventRepository := eventRepo.NewEventRepository(supabaseClient)
	cityRepository := placeRepo.NewCityRepository(supabaseClient)
	locationRepository := placeRepo.NewLocationRepository(supabaseClient)
	userProfileRepository := userRepo.NewUserProfileRepository(supabaseClient)
	userRepository := userRepo.NewUserRepository(supabaseAuth)

	// Initialize services
	cityService := placeSvc.NewCityService(cityRepository)
	locationService := placeSvc.NewLocationService(locationRepository)
	eventService := eventSvc.NewEventService(eventRepository, locationService, supabaseStorage)
	userProfileService := userSvc.NewUserProfileService(userProfileRepository, userRepository, supabaseStorage)

	// Initialize Gemini Handler
	geminiHandler, err := gemini.NewGeminiHandler(
		cfg,
		eventService,
		cityService,
		userProfileService,
	)
	if err != nil {
		appLogger.Error("Failed to initialize Gemini Handler", slog.Any("error", err))
		os.Exit(1)
	}

	aiGroup := router.Group("/ai")
	{
		aiGroup.POST("/chat/session", geminiHandler.CreateChatSession)
		aiGroup.POST("/chat/:sessionID/message", geminiHandler.SendMessage)
		aiGroup.GET("/event/:eventID/description", geminiHandler.GenerateEventDescription)
	}
}
