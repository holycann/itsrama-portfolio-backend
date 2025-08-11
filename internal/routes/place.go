package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/place/handlers"
)

// RegisterCityRoutes sets up all routes related to city management
func RegisterCityRoutes(
	r *gin.Engine,
	cityHandler *handlers.CityHandler,
	routeMiddleware *middleware.Middleware,
) {
	city := r.Group("/cities")
	{
		// Create a new city (Admin only)
		city.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", ""),
			cityHandler.CreateCity,
		)

		// Update an existing city (Admin only)
		city.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", ""),
			cityHandler.UpdateCity,
		)

		// Delete a city (Admin only)
		city.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", ""),
			cityHandler.DeleteCity,
		)

		// List cities (Public)
		city.GET("", cityHandler.ListCities)

		// Search cities with advanced filtering (Public)
		city.GET("/search", cityHandler.SearchCities)

		// Get city details by ID (Public)
		city.GET("/:id", cityHandler.GetCityByID)
	}
}

// RegisterProvinceRoutes sets up all routes related to province management
func RegisterProvinceRoutes(
	r *gin.Engine,
	provinceHandler *handlers.ProvinceHandler,
	routeMiddleware *middleware.Middleware,
) {
	province := r.Group("/provinces")
	{
		// Create a new province (Admin or Warlok)
		province.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			provinceHandler.CreateProvince,
		)

		// Update an existing province (Admin or Warlok)
		province.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			provinceHandler.UpdateProvince,
		)

		// Delete a province (Admin or Warlok)
		province.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			provinceHandler.DeleteProvince,
		)

		// List provinces (Public)
		province.GET("", provinceHandler.ListProvinces)

		// Search provinces with advanced filtering (Public)
		province.GET("/search", provinceHandler.SearchProvinces)

		// Get province details by ID (Public)
		province.GET("/:id", provinceHandler.GetProvinceByID)
	}
}

// RegisterLocationRoutes sets up all routes related to location management
func RegisterLocationRoutes(
	r *gin.Engine,
	locationHandler *handlers.LocationHandler,
	routeMiddleware *middleware.Middleware,
) {
	location := r.Group("/locations")
	{
		// Create a new location (Admin or Warlok)
		location.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			locationHandler.CreateLocation,
		)

		// List locations (Public)
		location.GET("", locationHandler.ListLocations)

		// Search locations with advanced filtering (Public)
		location.GET("/search", locationHandler.SearchLocations)

		// Get location details by ID (Public)
		location.GET("/:id", locationHandler.GetLocationByID)

		// Update an existing location (Admin or Warlok)
		location.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			locationHandler.UpdateLocation,
		)

		// Delete a location (Admin or Warlok)
		location.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			locationHandler.DeleteLocation,
		)
	}
}
