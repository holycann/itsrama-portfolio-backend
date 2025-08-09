package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/configs"
	"github.com/holycann/cultour-backend/docs"
	_ "github.com/holycann/cultour-backend/docs"

	// Achievement module
	achievementHandlers "github.com/holycann/cultour-backend/internal/achievement/handlers"
	achievementRepositories "github.com/holycann/cultour-backend/internal/achievement/repositories"
	achievementServices "github.com/holycann/cultour-backend/internal/achievement/services"

	// Cultural module
	culturalHandlers "github.com/holycann/cultour-backend/internal/cultural/handlers"
	culturalRepositories "github.com/holycann/cultour-backend/internal/cultural/repositories"
	culturalServices "github.com/holycann/cultour-backend/internal/cultural/services"

	// Discussion module
	discussionHandlers "github.com/holycann/cultour-backend/internal/discussion/handlers"
	discussionRepositories "github.com/holycann/cultour-backend/internal/discussion/repositories"
	discussionServices "github.com/holycann/cultour-backend/internal/discussion/services"

	// Gemini module
	"github.com/holycann/cultour-backend/internal/gemini"

	// Place module
	placeHandlers "github.com/holycann/cultour-backend/internal/place/handlers"
	placeRepositories "github.com/holycann/cultour-backend/internal/place/repositories"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"

	// User module
	userHandlers "github.com/holycann/cultour-backend/internal/users/handlers"
	userRepositories "github.com/holycann/cultour-backend/internal/users/repositories"
	userServices "github.com/holycann/cultour-backend/internal/users/services"

	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/routes"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/holycann/cultour-backend/pkg/supabase"
	"github.com/supabase-community/auth-go"
	supabaseGoClient "github.com/supabase-community/supabase-go"
)

// @title           Cultour API
// @version         1.0
// @description     Backend API for Cultour Project
// @contact.name    Cultour Development Team

// @host            localhost:8181
// @schemes         http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Insert Your JWT Token. Do NOT include "Bearer " prefix. Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
// @bearerFormat JWT

func main() {
	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := initializeLogger(cfg)

	// Set mode based on environment
	setMode(cfg)

	// Initialize infrastructure components
	infrastructure := initializeInfrastructure(cfg, appLogger)

	// Initialize repositories
	repositories := initializeRepositories(infrastructure)

	// Initialize services
	services := initializeServices(repositories, infrastructure)

	// Initialize handlers
	handlers := initializeHandlers(services, appLogger, cfg)

	// Initialize middleware
	routeMiddleware := initializeMiddleware(infrastructure, services, appLogger)

	// Create router
	router := createRouter(appLogger)

	// Register routes
	registerApplicationRoutes(router, routeMiddleware, handlers)

	// Start server
	startServer(router, cfg, appLogger)
}

// Infrastructure holds all external services and clients
type Infrastructure struct {
	SupabaseClient  *supabaseGoClient.Client
	SupabaseAuth    auth.Client
	SupabaseStorage *supabase.SupabaseStorage
	JWKS            *keyfunc.JWKS
}

// Repositories holds all repository instances
type Repositories struct {
	BadgeRepository       achievementRepositories.BadgeRepository
	EventRepository       culturalRepositories.EventRepository
	ThreadRepository      discussionRepositories.ThreadRepository
	MessageRepository     discussionRepositories.MessageRepository
	ParticipantRepository discussionRepositories.ParticipantRepository
	ProvinceRepository    placeRepositories.ProvinceRepository
	CityRepository        placeRepositories.CityRepository
	LocationRepository    placeRepositories.LocationRepository
	UserRepository        userRepositories.UserRepository
	UserProfileRepository userRepositories.UserProfileRepository
	UserBadgeRepository   userRepositories.UserBadgeRepository
}

// Services holds all service instances
type Services struct {
	BadgeService       achievementServices.BadgeService
	EventService       culturalServices.EventService
	ThreadService      discussionServices.ThreadService
	MessageService     discussionServices.MessageService
	ProvinceService    placeServices.ProvinceService
	CityService        placeServices.CityService
	LocationService    placeServices.LocationService
	UserService        userServices.UserService
	UserProfileService userServices.UserProfileService
	UserBadgeService   userServices.UserBadgeService
}

// Handlers holds all handler instances
type Handlers struct {
	BadgeHandler       *achievementHandlers.BadgeHandler
	EventHandler       *culturalHandlers.EventHandler
	ThreadHandler      *discussionHandlers.ThreadHandler
	MessageHandler     *discussionHandlers.MessageHandler
	GeminiHandler      *gemini.GeminiHandler
	ProvinceHandler    *placeHandlers.ProvinceHandler
	CityHandler        *placeHandlers.CityHandler
	LocationHandler    *placeHandlers.LocationHandler
	UserHandler        *userHandlers.UserHandler
	UserProfileHandler *userHandlers.UserProfileHandler
	UserBadgeHandler   *userHandlers.UserBadgeHandler
}

// loadConfiguration loads application configuration
func loadConfiguration() (*configs.Config, error) {
	cfg, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configurations: %w", err)
	}
	return cfg, nil
}

// initializeLogger sets up application logging
func initializeLogger(cfg *configs.Config) *logger.Logger {
	loggerConfig := logger.LoggerConfig{
		Path: cfg.Logging.FilePath,
	}
	appLogger := logger.NewLogger(loggerConfig)
	appLogger.Info("Application initialization started...")
	return appLogger
}

// setMode configures Gin mode and Swagger info
func setMode(cfg *configs.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		docs.SwaggerInfo.Host = cfg.Server.ProductionDomain
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		gin.SetMode(gin.DebugMode)
		docs.SwaggerInfo.Schemes = []string{"http"}
	}
}

// initializeInfrastructure sets up external services and connections
func initializeInfrastructure(cfg *configs.Config, appLogger *logger.Logger) *Infrastructure {
	// Initialize Supabase client
	supabaseClientConfig := supabase.SupabaseClientConfig{
		ProjectID: cfg.Supabase.ProjectID,
		ApiSecret: cfg.Supabase.ApiSecretKey,
	}
	supabaseClientInstance, err := supabase.NewSupabaseClient(supabaseClientConfig)
	if err != nil {
		appLogger.Error("Failed to initialize Supabase client", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	// Initialize Supabase auth
	supabaseAuthConfig := supabase.SupabaseAuthConfig{
		ApiKey:    cfg.Supabase.ApiSecretKey,
		ProjectID: cfg.Supabase.ProjectID,
	}
	supabaseAuthInstance := supabase.NewSupabaseAuth(supabaseAuthConfig)

	// Initialize Supabase storage
	supabaseStorageConfig := supabase.StorageConfig{
		JwtApiSecret:  cfg.Supabase.JwtApiKeySecret,
		ProjectID:     cfg.Supabase.ProjectID,
		BucketID:      cfg.Supabase.StorageBucketID,
		DefaultFolder: cfg.Supabase.DefaultFolder,
	}
	supabaseStorageInstance := supabase.NewSupabaseStorage(supabaseStorageConfig)

	// Initialize JWKS
	jwksURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/.well-known/jwks.json", cfg.Supabase.ProjectID)
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshUnknownKID: true,
		RefreshErrorHandler: func(err error) {
			appLogger.Error("JWKS refresh error", map[string]interface{}{"error": err.Error()})
		},
	})
	if err != nil {
		appLogger.Error("Failed to retrieve JWKS keys", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	return &Infrastructure{
		SupabaseClient:  supabaseClientInstance.GetClient(),
		SupabaseAuth:    supabaseAuthInstance.GetClient(),
		SupabaseStorage: supabaseStorageInstance,
		JWKS:            jwks,
	}
}

// initializeRepositories creates all repository instances
func initializeRepositories(infra *Infrastructure) *Repositories {
	return &Repositories{
		BadgeRepository:       achievementRepositories.NewBadgeRepository(infra.SupabaseClient),
		EventRepository:       culturalRepositories.NewEventRepository(infra.SupabaseClient),
		ThreadRepository:      discussionRepositories.NewThreadRepository(infra.SupabaseClient),
		MessageRepository:     discussionRepositories.NewMessageRepository(infra.SupabaseClient),
		ParticipantRepository: discussionRepositories.NewParticipantRepository(infra.SupabaseClient),
		ProvinceRepository:    placeRepositories.NewProvinceRepository(infra.SupabaseClient),
		CityRepository:        placeRepositories.NewCityRepository(infra.SupabaseClient),
		LocationRepository:    placeRepositories.NewLocationRepository(infra.SupabaseClient),
		UserRepository:        userRepositories.NewUserRepository(infra.SupabaseAuth),
		UserProfileRepository: userRepositories.NewUserProfileRepository(infra.SupabaseClient),
		UserBadgeRepository:   userRepositories.NewUserBadgeRepository(infra.SupabaseClient),
	}
}

// initializeServices creates all service instances
func initializeServices(repos *Repositories, infra *Infrastructure) *Services {
	locationService := placeServices.NewLocationService(repos.LocationRepository)
	badgeService := achievementServices.NewBadgeService(repos.BadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(repos.UserBadgeRepository)
	participantServices := discussionServices.NewParticipantService(repos.ParticipantRepository)
	threadService := discussionServices.NewThreadService(repos.ThreadRepository, participantServices)

	return &Services{
		BadgeService: badgeService,
		EventService: culturalServices.NewEventService(
			repos.EventRepository,
			locationService,
			*infra.SupabaseStorage,
			threadService,
		),
		ThreadService:   threadService,
		MessageService:  discussionServices.NewMessageService(repos.MessageRepository),
		ProvinceService: placeServices.NewProvinceService(repos.ProvinceRepository),
		CityService:     placeServices.NewCityService(repos.CityRepository),
		LocationService: locationService,
		UserService:     userServices.NewUserService(repos.UserRepository),
		UserProfileService: userServices.NewUserProfileService(
			repos.UserProfileRepository,
			repos.UserRepository,
			repos.UserBadgeRepository, // Use repository directly to match function signature
			badgeService,
			infra.SupabaseStorage,
		),
		UserBadgeService: userBadgeService,
	}
}

// initializeHandlers creates all handler instances
func initializeHandlers(svcs *Services, appLogger *logger.Logger, cfg *configs.Config) *Handlers {
	// Initialize Gemini handler
	geminiHandler, err := gemini.NewGeminiHandler(
		cfg,
		appLogger,
		svcs.EventService,
		svcs.CityService,
		svcs.ProvinceService,
		svcs.LocationService,
		svcs.UserService,
		svcs.BadgeService,
		svcs.UserProfileService,
		svcs.UserBadgeService,
	)
	if err != nil {
		appLogger.Error("Failed to initialize Gemini Handler", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	return &Handlers{
		BadgeHandler:       achievementHandlers.NewBadgeHandler(svcs.BadgeService, appLogger),
		EventHandler:       culturalHandlers.NewEventHandler(svcs.EventService, appLogger),
		ThreadHandler:      discussionHandlers.NewThreadHandler(svcs.ThreadService, appLogger),
		MessageHandler:     discussionHandlers.NewMessageHandler(svcs.MessageService, appLogger),
		GeminiHandler:      geminiHandler,
		ProvinceHandler:    placeHandlers.NewProvinceHandler(svcs.ProvinceService, appLogger),
		CityHandler:        placeHandlers.NewCityHandler(svcs.CityService, appLogger),
		LocationHandler:    placeHandlers.NewLocationHandler(svcs.LocationService, appLogger),
		UserHandler:        userHandlers.NewUserHandler(svcs.UserService, appLogger),
		UserProfileHandler: userHandlers.NewUserProfileHandler(svcs.UserProfileService, appLogger),
		UserBadgeHandler:   userHandlers.NewUserBadgeHandler(svcs.UserBadgeService, appLogger),
	}
}

// initializeMiddleware creates middleware
func initializeMiddleware(infra *Infrastructure, svcs *Services, appLogger *logger.Logger) *middleware.Middleware {
	return middleware.NewMiddleware(infra.SupabaseAuth, infra.JWKS, svcs.UserBadgeService, svcs.BadgeService, appLogger)
}

// createRouter sets up Gin router with middleware
func createRouter(appLogger *logger.Logger) *gin.Engine {
	router := gin.New()

	// Disable automatic redirects
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Allow all origins during development
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight requests
	}))

	// Custom logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		appLogger.Info("Request processed",
			map[string]interface{}{
				"status":  c.Writer.Status(),
				"method":  c.Request.Method,
				"path":    c.Request.URL.Path,
				"latency": duration,
				"origin":  c.GetHeader("Origin"),
			},
		)
	})

	// Add a catch-all route for debugging
	router.NoRoute(func(c *gin.Context) {
		appLogger.Error("No Route Found",
			map[string]interface{}{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			},
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Route not found",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	})

	return router
}

// registerApplicationRoutes adds all application routes
func registerApplicationRoutes(router *gin.Engine, routeMiddleware *middleware.Middleware, handlers *Handlers) {
	// Setup basic routes and Swagger
	routes.SetupRouter(router)

	// Register all routes with their handlers
	routes.SetupGeminiRoutes(router, handlers.GeminiHandler, routeMiddleware)
	routes.RegisterEventRoutes(router, handlers.EventHandler, routeMiddleware)
	routes.RegisterLocationRoutes(router, handlers.LocationHandler, routeMiddleware)
	routes.RegisterUserRoutes(router, handlers.UserHandler, routeMiddleware)
	routes.RegisterUserProfileRoutes(router, handlers.UserProfileHandler, routeMiddleware)
	routes.RegisterUserBadgeRoutes(router, handlers.UserBadgeHandler, routeMiddleware)
	routes.RegisterBadgeRoutes(router, handlers.BadgeHandler, routeMiddleware)
	routes.RegisterCityRoutes(router, handlers.CityHandler, routeMiddleware)
	routes.RegisterProvinceRoutes(router, handlers.ProvinceHandler, routeMiddleware)
	routes.RegisterMessageRoutes(router, handlers.MessageHandler, routeMiddleware)
	routes.RegisterThreadRoutes(router, handlers.ThreadHandler, routeMiddleware)
}

// startServer runs the HTTP server with graceful shutdown
func startServer(router *gin.Engine, cfg *configs.Config, appLogger *logger.Logger) {
	// Configure server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Create server
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("Starting server", map[string]interface{}{"address": serverAddr})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed", map[string]interface{}{"error": err.Error()})
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown error", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	appLogger.Info("Server exited")
}
