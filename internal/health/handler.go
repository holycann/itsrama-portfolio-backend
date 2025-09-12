package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/pkg/base"
	supabaseClient "github.com/supabase-community/supabase-go"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	base.BaseHandler
	supabaseClient *supabaseClient.Client
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(supabaseClient *supabaseClient.Client) *HealthHandler {
	return &HealthHandler{
		supabaseClient: supabaseClient,
	}
}

// GetHealthStatus handles the health check endpoint
func (h *HealthHandler) GetHealthStatus(c *gin.Context) {
	// Perform health check
	healthStatus := CheckHealth(h.supabaseClient)

	// Determine HTTP status code based on health status
	statusCode := http.StatusOK
	if healthStatus.Status == "degraded" {
		statusCode = http.StatusPartialContent
	} else if healthStatus.Status == "critical" {
		statusCode = http.StatusServiceUnavailable
	}

	// Respond with health status
	c.JSON(statusCode, healthStatus)
}
