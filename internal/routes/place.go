package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/place/handlers"
	"github.com/holycann/cultour-backend/internal/place/repositories"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/supabase-community/supabase-go"
)

func RegisterCityRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
) {
	cityRepository := repositories.NewCityRepository(supabaseClient, *repositories.DefaultCityConfig())
	cityService := services.NewCityService(cityRepository)
	cityHandler := handlers.NewCityHandler(cityService, logger)

	city := r.Group("/cities")
	{
		city.POST("/", cityHandler.CreateCity)
		city.GET("/", cityHandler.ListCities)
		city.GET("/search", cityHandler.SearchCity)
		city.PUT("/:id", cityHandler.UpdateCity)
		city.DELETE("/:id", cityHandler.DeleteCity)
	}
}

func RegisterLocationRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
) {
	locationRepository := repositories.NewLocationRepository(supabaseClient, *repositories.DefaultLocationConfig())
	locationService := services.NewLocationService(locationRepository)
	locationHandler := handlers.NewLocationHandler(locationService, logger)

	location := r.Group("/locations")
	{
		location.POST("/", locationHandler.CreateLocation)
		location.GET("/", locationHandler.ListLocation)
		location.GET("/search", locationHandler.SearchLocations)
		location.PUT("/:id", locationHandler.UpdateLocation)
		location.DELETE("/:id", locationHandler.DeleteLocation)
	}
}
