package middleware

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/supabase-community/auth-go"
)

// Middleware handles JWT token authentication and validation
type Middleware struct {
	supabaseAuth     auth.Client
	jwks             *keyfunc.JWKS
	userBadgeService userServices.UserBadgeService
	badgeService     services.BadgeService
}

// UserContext represents the authenticated user's context
type UserContext struct {
	ID    string
	Email string
	Role  string
	Badge string
}

// NewMiddleware creates a new JWT middleware instance
func NewMiddleware(
	supabaseAuth auth.Client,
	jwks *keyfunc.JWKS,
	userBadgeService userServices.UserBadgeService,
	badgeService services.BadgeService,
) *Middleware {
	return &Middleware{
		supabaseAuth:     supabaseAuth,
		jwks:             jwks,
		userBadgeService: userBadgeService,
		badgeService:     badgeService,
	}
}

// VerifyJWT validates the JWT token from the Authorization header
func (m *Middleware) VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			details, _ := json.Marshal(map[string]interface{}{
				"authorization_header": "missing",
			})
			response.Unauthorized(c, "Missing authorization token", string(details), "")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			details, _ := json.Marshal(map[string]interface{}{
				"token_format": "invalid",
			})
			response.Unauthorized(c, "Invalid token format", string(details), "")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, m.jwks.Keyfunc)
		if err != nil || !token.Valid {
			details, _ := json.Marshal(map[string]interface{}{
				"token_validation": "failed",
				"error":            err.Error(),
			})
			response.Unauthorized(c, "Invalid token", string(details), "")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			details, _ := json.Marshal(map[string]interface{}{
				"token_claims": "invalid",
			})
			response.Unauthorized(c, "Invalid token claims", string(details), "")
			c.Abort()
			return
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		// Retrieve user badges (multiple) and map to badge names
		var badgeNames []string
		if userID != "" {
			userBadges, err := m.userBadgeService.GetUserBadgesByUser(c.Request.Context(), userID)
			if err == nil && len(userBadges) > 0 {
				for _, ub := range userBadges {
					// Get badge name from badge service
					badge, err := m.badgeService.GetBadgeByID(c.Request.Context(), ub.BadgeID.String())
					if err == nil && badge != nil && badge.Name != "" {
						badgeNames = append(badgeNames, strings.ToLower(badge.Name))
					}
				}
			}
		}

		// Fallback to JWT claim if no badge found
		if len(badgeNames) == 0 {
			// Try to get "badge" claim as string or []interface{}
			if badgeVal, ok := claims["badge"]; ok {
				switch v := badgeVal.(type) {
				case string:
					if v != "" {
						badgeNames = append(badgeNames, strings.ToLower(v))
					}
				case []interface{}:
					for _, b := range v {
						if s, ok := b.(string); ok && s != "" {
							badgeNames = append(badgeNames, strings.ToLower(s))
						}
					}
				}
			}
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)
		c.Set("badges", badgeNames)

		c.Next()
	}
}

// RequireRoleOrBadge creates a middleware to check user roles OR badges
func (m *Middleware) RequireRoleOrBadge(allowedRoles string, allowedBadges string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		userBadgesVal, _ := c.Get("badges")

		// Convert userBadgesVal to []string
		var userBadges []string
		switch v := userBadgesVal.(type) {
		case []string:
			for _, badge := range v {
				userBadges = append(userBadges, strings.ToLower(badge))
			}
		case []interface{}:
			for _, b := range v {
				if s, ok := b.(string); ok {
					userBadges = append(userBadges, strings.ToLower(s))
				}
			}
		case string:
			if v != "" {
				userBadges = []string{strings.ToLower(v)}
			}
		}

		// Check roles
		roleAllowed := false
		if allowedRoles != "" {
			roles := strings.Split(allowedRoles, ",")
			for _, role := range roles {
				if userRole == strings.TrimSpace(role) {
					roleAllowed = true
					break
				}
			}
		}

		// Check badges (by badge name)
		badgeAllowed := false
		if allowedBadges != "" && len(userBadges) > 0 {
			allowedBadgeList := strings.Split(allowedBadges, ",")
			for _, allowedBadge := range allowedBadgeList {
				allowedBadge = strings.ToLower(strings.TrimSpace(allowedBadge))
				for _, userBadge := range userBadges {
					if userBadge == allowedBadge {
						badgeAllowed = true
						break
					}
				}
				if badgeAllowed {
					break
				}
			}
		}

		if !roleAllowed && !badgeAllowed {
			// Convert map to JSON string for error details
			details, _ := json.Marshal(map[string]interface{}{
				"required_roles":  strings.Split(allowedRoles, ","),
				"required_badges": strings.Split(allowedBadges, ","),
				"user_role":       userRole,
				"user_badges":     userBadges,
			})
			logger.DefaultLogger().Info("DETAILS:", string(details))
			response.Forbidden(c, "Insufficient permissions", string(details), "")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext retrieves user information from Gin context
func GetUserFromContext(c *gin.Context) (userID, email, role string, badges []string, err error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return "", "", "", nil, errors.New("user_id not found in context")
	}

	emailVal, exists := c.Get("email")
	if !exists {
		return "", "", "", nil, errors.New("email not found in context")
	}

	roleVal, exists := c.Get("role")
	if !exists {
		return "", "", "", nil, errors.New("role not found in context")
	}

	badgesVal, _ := c.Get("badges")
	var badgeList []string
	switch v := badgesVal.(type) {
	case []string:
		for _, badge := range v {
			badgeList = append(badgeList, strings.ToLower(badge))
		}
	case []interface{}:
		for _, b := range v {
			if s, ok := b.(string); ok {
				badgeList = append(badgeList, strings.ToLower(s))
			}
		}
	case string:
		if v != "" {
			badgeList = []string{strings.ToLower(v)}
		}
	}

	return userIDVal.(string), emailVal.(string), roleVal.(string), badgeList, nil
}
