package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
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
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/middleware"
	"github.com/holycann/cultour-backend/internal/routes"
	"github.com/holycann/cultour-backend/internal/supabase"
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

	// Initialize Supabase client
	supabaseClient, err := supabase.NewSupabaseClient(supabase.SupabaseClientConfig{
		ProjectID: cfg.Supabase.ProjectID,
		ApiSecret: cfg.Supabase.ApiSecretKey,
	})
	if err != nil {
		appLogger.Error("Failed to initialize Supabase client", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize Supabase auth
	supabaseAuth := supabase.NewSupabaseAuth(supabase.SupabaseAuthConfig{
		ApiKey:    cfg.Supabase.ApiSecretKey,
		ProjectID: cfg.Supabase.ProjectID,
	})

	// Initialize Supabase storage
	supabaseStorage := supabase.NewSupabaseStorage(supabase.SupabaseStorageConfig{
		JwtApiSecret:  cfg.Supabase.JwtApiKeySecret,
		ProjectID:     cfg.Supabase.ProjectID,
		BucketID:      cfg.Supabase.StorageBucketID,
		DefaultFolder: cfg.Supabase.DefaultFolder,
	})

	jwks, err := keyfunc.Get(fmt.Sprintf("https://%s.supabase.co/auth/v1/.well-known/jwks.json", cfg.Supabase.ProjectID), keyfunc.Options{
		RefreshUnknownKID: true,
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc: %s", err.Error())
		},
	})
	if err != nil {
		appLogger.Error("Failed to retrieve JWKS keys", slog.Any("error", err))
		os.Exit(1)
	}

	routeMiddleware := middleware.NewMiddleware(supabaseAuth.GetClient(), jwks)

	// Create router
	router := createRouter(appLogger)

	// Register routes
	registerApplicationRoutes(router, supabaseClient, supabaseAuth, supabaseStorage, routeMiddleware, appLogger, cfg)

	// Start server
	startServer(router, cfg, appLogger)
}

// loadConfiguration loads application configuration
// @Summary Load application configuration
// @Description Reads and validates the application configuration
// @Produce json
// @Success 200 {object} configs.Config
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
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Duration("latency", duration),
			slog.String("origin", c.GetHeader("Origin")),
		)
	})

	// Add a catch-all route for debugging
	router.NoRoute(func(c *gin.Context) {
		appLogger.Error("No Route Found",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
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
func registerApplicationRoutes(
	router *gin.Engine,
	supabaseClient *supabase.SupabaseClient,
	supabaseAuth *supabase.SupabaseAuth,
	supabaseStorage *supabase.SupabaseStorage,
	routeMiddleware *middleware.Middleware,
	appLogger *logger.Logger,
	config *configs.Config,
) {
	// Setup Gemini routes
	routes.SetupRouter(router)
	routes.SetupGeminiRoutes(router, config, supabaseClient.GetClient(), supabaseAuth.GetClient(), supabaseStorage, appLogger)
	routes.RegisterEventRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
	routes.RegisterLocationRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
	routes.RegisterUserRoutes(router, supabaseAuth, routeMiddleware)
	routes.RegisterUserProfileRoutes(router, supabaseClient, supabaseAuth, supabaseStorage, routeMiddleware)
	routes.RegisterUserBadgeRoutes(router, supabaseClient, routeMiddleware, appLogger)
	routes.RegisterBadgeRoutes(router, supabaseClient.GetClient(), routeMiddleware, appLogger)
	routes.RegisterCityRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
	routes.RegisterProvinceRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
	routes.RegisterLocalStoryRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
	routes.RegisterMessageRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
	routes.RegisterThreadRoutes(router, appLogger, supabaseClient.GetClient(), routeMiddleware)
}

// startServer runs the HTTP server with graceful shutdown
func startServer(
	router *gin.Engine,
	cfg *configs.Config,
	appLogger *logger.Logger,
) {
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
		appLogger.Info("Starting server", slog.String("address", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed", slog.Any("error", err))
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
		appLogger.Error("Server shutdown error", slog.Any("error", err))
		os.Exit(1)
	}

	appLogger.Info("Server exited")
}
