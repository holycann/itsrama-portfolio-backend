package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/handlers"
	"github.com/holycann/cultour-backend/internal/middleware"
)

// RegisterEventRoutes sets up routes for event-related operations
// @Description Configures HTTP routes for event management with authentication and authorization
// @Tags Events
func RegisterEventRoutes(
	r *gin.Engine,
	eventHandler *handlers.EventHandler,
	routeMiddleware *middleware.Middleware,
) {
	event := r.Group("/events")
	{
		// Create a new event
		// @Summary Create a new cultural event
		// @Description Allows authenticated users to add a new cultural event to the platform
		// @Tags Events
		event.POST("",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.CreateEvent,
		)

		// List events
		// @Summary Retrieve a list of events
		// @Description Fetches paginated list of events with optional filtering
		// @Tags Events
		event.GET("", eventHandler.ListEvents)

		// Search events
		// @Summary Search events by query
		// @Description Performs full-text search across event fields
		// @Tags Events
		event.GET("/search", eventHandler.SearchEvents)

		// Get trending events
		// @Summary Retrieve trending events
		// @Description Fetches most popular or recently viewed events
		// @Tags Events
		event.GET("/trending", eventHandler.GetTrendingEvents)

		// Get related events for a specific event
		// @Summary Retrieve related events
		// @Description Finds events similar to a specific event based on location and other criteria
		// @Tags Events
		event.GET("/:id/related", eventHandler.GetRelatedEvents)

		// Get a specific event by ID
		// @Summary Retrieve a specific event
		// @Description Fetches comprehensive details of an event by its unique identifier
		// @Tags Events
		event.GET("/:id", eventHandler.GetEventByID)

		// Update event views
		// @Summary Increment event view count
		// @Description Tracks and updates the number of times an event has been viewed
		// @Tags Events
		event.POST("/:id/views",
			routeMiddleware.VerifyJWT(),
			eventHandler.UpdateEventViews,
		)

		// Update an event
		// @Summary Update an existing event
		// @Description Allows event creator or administrator to modify event details
		// @Tags Events
		event.PUT("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.UpdateEvent,
		)

		// Delete an event
		// @Summary Delete an existing event
		// @Description Allows event creator or administrator to remove an event from the platform
		// @Tags Events
		event.DELETE("/:id",
			routeMiddleware.VerifyJWT(),
			routeMiddleware.RequireRoleOrBadge("admin", "warlok"),
			eventHandler.DeleteEvent,
		)
	}
}
