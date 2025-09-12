package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/configs"
	"github.com/holycann/itsrama-portfolio-backend/internal/experience"
	"github.com/holycann/itsrama-portfolio-backend/internal/health"
	"github.com/holycann/itsrama-portfolio-backend/internal/middleware"
	"github.com/holycann/itsrama-portfolio-backend/internal/project"
	"github.com/holycann/itsrama-portfolio-backend/internal/routes"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
	"github.com/holycann/itsrama-portfolio-backend/pkg/response"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
)

// AppDependencies holds all the initialized dependencies
type AppDependencies struct {
	Config        *configs.Config
	Logger        *logger.Logger
	JWKS          *keyfunc.JWKS
	JWTMiddleware *middleware.Middleware
	Router        *gin.Engine

	// Supabase Dependencies
	SupabaseDefault *supabase.SupabaseClient
	SupabaseAuth    *supabase.SupabaseAuth
	SupabaseStorage *supabase.SupabaseStorage
}

type FeatureDependencies struct {
	// Health Dependencies
	HealthHandler *health.HealthHandler

	// Experience Dependencies
	ExperienceHandler    *experience.ExperienceHandler
	ExperienceService    *experience.ExperienceService
	ExperienceRepository *experience.ExperienceRepository

	// Project Dependencies
	ProjectHandler    *project.ProjectHandler
	ProjectService    *project.ProjectService
	ProjectRepository *project.ProjectRepository

	// Tech Stack Dependencies
	TechStackHandler    *tech_stack.TechStackHandler
	TechStackService    *tech_stack.TechStackService
	TechStackRepository *tech_stack.TechStackRepository
}

func main() {
	// Initialize application context
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		select {
		case <-ctx.Done():
		default:
		}
	}()

	// Initialize dependencies
	deps, err := initializeAppDependencies()
	if err != nil {
		fmt.Printf("Failed to initialize dependencies: %v\n", err)
		os.Exit(1)
	}
	defer cleanupAppDependencies(deps)

	// Initialize dependencies
	featureDeps, err := initializeFeatureDependencies(deps.SupabaseDefault, *deps.SupabaseStorage, deps.Config, deps.Logger)
	if err != nil {
		fmt.Printf("Failed to initialize dependencies: %v\n", err)
		os.Exit(1)
	}

	// Setup routes
	setupRoutes(deps, featureDeps)

	// Start server
	server := createHTTPServer(deps)

	// Graceful server startup and shutdown
	go startServer(server, deps.Logger, deps.Config)

	// Wait for shutdown signal
	waitForShutdown(server, deps.Logger, deps.Config)
}

// initializeDependencies sets up all application dependencies
func initializeAppDependencies() (*AppDependencies, error) {
	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logging
	appLogger := initializeLogger(cfg)

	// Initialize Supabase default schema client
	supabaseDefault, err := supabase.NewSupabaseClient(supabase.SupabaseClientConfig{
		ApiSecret: cfg.Supabase.ApiSecretKey,
		ProjectID: cfg.Supabase.ProjectID,
		Schema:    cfg.Database.Schema,
	})
	if err != nil {
		appLogger.Error("Failed to initialize Supabase client with default schema", "error", err)
		os.Exit(1)
	}

	// Initialize Supabase authentication
	supabaseAuth, err := supabase.NewSupabaseAuth(supabase.SupabaseAuthConfig{
		ApiKey:    cfg.Supabase.ApiSecretKey,
		ProjectID: cfg.Supabase.ProjectID,
	})
	if err != nil {
		appLogger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize Supabase authentication
	supabaseStorage, err := supabase.NewSupabaseStorage(supabase.StorageConfig{
		ProjectID:           cfg.Supabase.ProjectID,
		JwtApiSecret:        cfg.Supabase.JWTSecret,
		BucketID:            cfg.Supabase.StorageBucketID,
		DefaultFolder:       cfg.Supabase.DefaultStorageFolder,
		MaxFileSize:         cfg.Supabase.MaxFileSize,
		AllowedFileTypes:    cfg.Supabase.AllowedFileTypes,
		DefaultCacheControl: cfg.Supabase.CacheControl,
	})
	if err != nil {
		appLogger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize JWKS
	jwks := initializeJWKS(cfg, appLogger)

	// Allowed Emails For Backend Access
	allowedEmails := []string{
		"ehhramaa212@gmail.com",
		"muhamad.ramadhan.dev@gmail.com",
	}

	// Initialize JWT middleware
	jwtMiddleware := initializeJWTMiddleware(jwks, allowedEmails, appLogger)

	// Setup Gin router
	router := initializeRouter(appLogger, cfg)

	return &AppDependencies{
		Config:          cfg,
		Logger:          appLogger,
		SupabaseDefault: supabaseDefault,
		SupabaseAuth:    supabaseAuth,
		SupabaseStorage: supabaseStorage,
		JWKS:            jwks,
		JWTMiddleware:   jwtMiddleware,
		Router:          router,
	}, nil
}

func initializeFeatureDependencies(supabaseDefault *supabase.SupabaseClient, supabaseStorage supabase.SupabaseStorage, cfg *configs.Config, appLogger *logger.Logger) (*FeatureDependencies, error) {
	// Initialize health dependencies
	healthHandler := health.NewHealthHandler(supabaseDefault.GetClient())

	// Initialize experience dependencies
	experienceRepo := experience.NewExperienceRepository(supabaseDefault, supabaseStorage)
	experienceService := experience.NewExperienceService(experienceRepo, supabaseStorage)
	experienceHandler := experience.NewExperienceHandler(experienceService, appLogger)

	// Initialize project dependencies
	projectRepo := project.NewProjectRepository(supabaseDefault, supabaseStorage)
	projectService := project.NewProjectService(projectRepo, supabaseStorage)
	projectHandler := project.NewProjectHandler(projectService, appLogger)

	// Initialize tech stack dependencies
	techStackRepo := tech_stack.NewTechStackRepository(supabaseDefault)
	techStackService := tech_stack.NewTechStackService(techStackRepo)
	techStackHandler := tech_stack.NewTechStackHandler(techStackService, appLogger)

	return &FeatureDependencies{
		// Health Dependencies
		HealthHandler: healthHandler,

		// Experience Dependencies
		ExperienceHandler:    experienceHandler,
		ExperienceService:    &experienceService,
		ExperienceRepository: &experienceRepo,

		// Project Dependencies
		ProjectHandler:    projectHandler,
		ProjectService:    &projectService,
		ProjectRepository: &projectRepo,

		// Tech Stack Dependencies
		TechStackHandler:    techStackHandler,
		TechStackService:    &techStackService,
		TechStackRepository: &techStackRepo,
	}, nil
}

// cleanupDependencies performs cleanup for all initialized dependencies
func cleanupAppDependencies(deps *AppDependencies) {
	// Close logger
	if err := deps.Logger.Close(); err != nil {
		fmt.Printf("Error closing logger: %v\n", err)
	}
}

// setupRoutes configures all application routes
func setupRoutes(deps *AppDependencies, featureDeps *FeatureDependencies) {
	// Setup global error handler
	deps.Router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "route_not_found", "Endpoint not found", c.Request.URL.Path)
	})

	// Setup API routes
	v1Group := deps.Router.Group("/api/v1")
	{
		v1Group.GET("/", func(c *gin.Context) {
			apiInfo := map[string]string{
				"name":          "Itsrama Portfolio Backend API",
				"description":   "Comprehensive backend API for Itsrama Portfolio",
				"documentation": "https://github.com/holycann/itsrama-portfolio-backend",
				"status":        "operational",
				"version":       "1.0.0",
				"environment":   deps.Config.Environment,
			}

			response.Success(c, http.StatusOK, apiInfo, "API Info")
		})

		// Health check endpoint with comprehensive system checks
		v1Group.GET("/health", featureDeps.HealthHandler.GetHealthStatus)

		// Experience Routes
		routes.RegisterExperienceRoutes(
			v1Group,
			featureDeps.ExperienceHandler,
			deps.JWTMiddleware,
		)

		// Project Routes
		routes.RegisterProjectRoutes(
			v1Group,
			featureDeps.ProjectHandler,
			deps.JWTMiddleware,
		)

		// Tech Stack Routes
		routes.RegisterTechStackRoutes(
			v1Group,
			featureDeps.TechStackHandler,
			deps.JWTMiddleware,
		)
	}
}

// createHTTPServer creates and configures the HTTP server
func createHTTPServer(deps *AppDependencies) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", deps.Config.Server.Host, deps.Config.Server.Port),
		Handler:      deps.Router,
		ReadTimeout:  time.Duration(deps.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(deps.Config.Server.WriteTimeout) * time.Second,
	}
}

// startServer handles the server startup process
func startServer(server *http.Server, log *logger.Logger, cfg *configs.Config) {
	log.Info("Starting server",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
	)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Server startup failed", "error", err)
		os.Exit(1)
	}
}

// waitForShutdown handles graceful shutdown of the server
func waitForShutdown(server *http.Server, log *logger.Logger, cfg *configs.Config) {
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.Server.ShutdownTimeout)*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown failed", "error", err)
	}

	log.Info("Server exited")
}

// initializeLogger sets up the application logger
func initializeLogger(cfg *configs.Config) *logger.Logger {
	loggerConfig := logger.LoggerConfig{
		Path:        cfg.Logging.FilePath,
		Level:       logger.InfoLevel,
		Development: cfg.Environment == "development",
		MaxSize:     cfg.Logging.MaxSize,
		MaxBackups:  cfg.Logging.MaxBackups,
		MaxAge:      cfg.Logging.MaxAge,
		Compress:    cfg.Logging.Compress,
	}

	return logger.NewLogger(loggerConfig)
}

// initializeJWKS retrieves JWKS keys for JWT validation
func initializeJWKS(cfg *configs.Config, log *logger.Logger) *keyfunc.JWKS {
	jwksURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/.well-known/jwks.json", cfg.Supabase.ProjectID)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshUnknownKID: true,
		RefreshErrorHandler: func(err error) {
			log.Error("JWKS refresh error", "error", err)
		},
	})

	if err != nil {
		log.Error("Failed to retrieve JWKS keys", "error", err)
		os.Exit(1)
	}

	log.Info("JWKS keys initialized successfully")
	return jwks
}

// initializeJWTMiddleware creates JWT authentication middleware
func initializeJWTMiddleware(
	jwks *keyfunc.JWKS,
	allowedEmails []string,
	log *logger.Logger,
) *middleware.Middleware {
	return middleware.NewMiddleware(
		jwks,
		allowedEmails,
		log,
	)
}

// initializeRouter sets up the Gin router with global middleware
func initializeRouter(log *logger.Logger, cfg *configs.Config) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())

	// CORS Middleware
	// CORS Middleware
	if cfg.CORS.CORSEnabled {
		if cfg.Environment != "production" {
			cfg.CORS.MaxAge = 0
		}

		corsConfig := cors.Config{
			AllowOrigins:     cfg.CORS.AllowedOrigins,
			AllowMethods:     cfg.CORS.AllowedMethods,
			AllowHeaders:     cfg.CORS.AllowedHeaders,
			ExposeHeaders:    cfg.CORS.ExposedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
		}

		router.Use(cors.New(corsConfig))
	} else {
		router.Use(cors.New(cors.DefaultConfig()))
	}

	// Logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info("Request processed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
		)
	})

	return router
}
