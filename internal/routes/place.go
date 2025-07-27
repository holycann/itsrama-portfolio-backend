package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/place/handlers"
	"github.com/holycann/cultour-backend/internal/place/repositories"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/supabase-community/supabase-go"
)

func RegisterCityRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
	routeMiddleware ...*middleware.Middleware,
) {
	cityRepository := repositories.NewCityRepository(supabaseClient)
	cityService := services.NewCityService(cityRepository)
	cityHandler := handlers.NewCityHandler(cityService, logger)

	var mw *middleware.Middleware
	if len(routeMiddleware) > 0 {
		mw = routeMiddleware[0]
	}

	city := r.Group("/cities")
	{
		if mw != nil {
			city.POST("",
				mw.VerifyJWT(),
				mw.RequireRoleOrBadge("admin", "warlok"),
				cityHandler.CreateCity,
			)
			city.PUT("/:id",
				mw.VerifyJWT(),
				mw.RequireRoleOrBadge("admin", "warlok"),
				cityHandler.UpdateCity,
			)
			city.DELETE("/:id",
				mw.VerifyJWT(),
				mw.RequireRoleOrBadge("admin", "warlok"),
				cityHandler.DeleteCity,
			)
		} else {
			city.POST("", cityHandler.CreateCity)
			city.PUT("/:id", cityHandler.UpdateCity)
			city.DELETE("/:id", cityHandler.DeleteCity)
		}
		city.GET("", cityHandler.ListCities)
		city.GET("/search", cityHandler.SearchCities)
		city.GET("/:id", cityHandler.GetCityByID)
	}
}

func RegisterProvinceRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
	routeMiddleware ...*middleware.Middleware,
) {
	provinceRepository := repositories.NewProvinceRepository(supabaseClient)
	provinceService := services.NewProvinceService(provinceRepository)
	provinceHandler := handlers.NewProvinceHandler(provinceService, logger)

	var mw *middleware.Middleware
	if len(routeMiddleware) > 0 {
		mw = routeMiddleware[0]
	}

	province := r.Group("/provinces")
	{
		if mw != nil {
			province.POST("",
				mw.VerifyJWT(),
				mw.RequireRoleOrBadge("admin", "warlok"),
				provinceHandler.CreateProvince,
			)
			province.PUT("/:id",
				mw.VerifyJWT(),
				mw.RequireRoleOrBadge("admin", "warlok"),
				provinceHandler.UpdateProvince,
			)
			province.DELETE("/:id",
				mw.VerifyJWT(),
				mw.RequireRoleOrBadge("admin", "warlok"),
				provinceHandler.DeleteProvince,
			)
		} else {
			province.POST("", provinceHandler.CreateProvince)
			province.PUT("/:id", provinceHandler.UpdateProvince)
			province.DELETE("/:id", provinceHandler.DeleteProvince)
		}
		province.GET("", provinceHandler.ListProvinces)
		province.GET("/search", provinceHandler.SearchProvinces)
		province.GET("/:id", provinceHandler.GetProvinceByID)
	}
}

func RegisterLocationRoutes(
	r *gin.Engine,
	logger *logger.Logger,
	supabaseClient *supabase.Client,
	routeMiddleware *middleware.Middleware,
) {
	locationRepository := repositories.NewLocationRepository(supabaseClient)
	locationService := services.NewLocationService(locationRepository)
	locationHandler := handlers.NewLocationHandler(locationService, logger)

	location := r.Group("/locations")
	{
		location.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			locationHandler.CreateLocation,
		)
		location.GET("", locationHandler.ListLocations)
		location.GET("/search", locationHandler.SearchLocations)
		location.GET("/:id", locationHandler.GetLocationByID)
		location.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			locationHandler.UpdateLocation,
		)
		location.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			locationHandler.DeleteLocation,
		)
	}
}
