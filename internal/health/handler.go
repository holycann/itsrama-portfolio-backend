package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/internal/base"
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

// GetHealthStatus godoc
// @Summary      Check Health Status
// @Description  Performs a health check on the service and its dependencies
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthStatus  "Service is healthy"
// @Success      206  {object}  HealthStatus  "Service is degraded"
// @Failure      503  {object}  HealthStatus  "Service is unavailable"
// @Failure      500  {object}  HealthStatus  "Internal server error"
// @Router       /health [get]
func (h *HealthHandler) GetHealthStatus(c *gin.Context) {
	// Perform health check
	healthStatus := CheckHealth(h.supabaseClient)

	// Determine HTTP status code based on health status
	var statusCode int
	switch healthStatus.Status {
	case "ok":
		statusCode = http.StatusOK
	case "degraded":
		statusCode = http.StatusPartialContent
	case "critical":
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusInternalServerError
	}

	// Respond with health status
	c.JSON(statusCode, healthStatus)
}
