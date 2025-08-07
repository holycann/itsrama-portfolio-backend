package middleware

import (
	"context"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/holycann/cultour-backend/internal/achievement/services"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/holycann/cultour-backend/pkg/response"
	"github.com/supabase-community/auth-go"
)

// Middleware handles JWT token authentication and validation
type Middleware struct {
	supabaseAuth     auth.Client
	jwks             *keyfunc.JWKS
	userBadgeService userServices.UserBadgeService
	badgeService     services.BadgeService
	logger           *logger.Logger
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
	logger *logger.Logger,
) *Middleware {
	return &Middleware{
		supabaseAuth:     supabaseAuth,
		jwks:             jwks,
		userBadgeService: userBadgeService,
		badgeService:     badgeService,
		logger:           logger,
	}
}

// VerifyJWT validates the JWT token from the Authorization header
func (m *Middleware) VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.handleAuthError(c, "Missing authorization token",
				errors.WithContext("authorization_header", "missing"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			m.handleAuthError(c, "Invalid token format",
				errors.WithContext("token_format", "invalid"))
			return
		}

		token, err := jwt.Parse(tokenString, m.jwks.Keyfunc)
		if err != nil || !token.Valid {
			m.handleAuthError(c, "Invalid token",
				errors.WithContext("token_validation", "failed"),
				errors.WithContext("error", err.Error()))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.handleAuthError(c, "Invalid token claims",
				errors.WithContext("token_claims", "invalid"))
			return
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		// Retrieve user badges (multiple) and map to badge names
		badgeNames, err := m.getUserBadges(c.Request.Context(), userID)
		if err != nil {
			m.logger.Error("Failed to retrieve user badges", "error", err)
		}

		// Fallback to JWT claim if no badge found
		if len(badgeNames) == 0 {
			badgeNames = m.extractBadgesFromClaims(claims)
		}

		// Set user context
		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)
		c.Set("badges", badgeNames)

		c.Next()
	}
}

// handleAuthError handles authentication errors with standardized response
func (m *Middleware) handleAuthError(c *gin.Context, message string, opts ...func(*errors.CustomError)) {
	errors.New(
		errors.ErrAuthentication,
		message,
		nil,
		opts...,
	)
	response.Unauthorized(c, "auth_error", message, "")
	c.Abort()
}

// getUserBadges retrieves user badges from the database
func (m *Middleware) getUserBadges(ctx context.Context, userID string) ([]string, error) {
	var badgeNames []string
	if userID == "" {
		return badgeNames, nil
	}

	userBadges, err := m.userBadgeService.GetUserBadgesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, ub := range userBadges {
		badge, err := m.badgeService.GetBadgeByID(ctx, ub.BadgeID.String())
		if err == nil && badge != nil && badge.Name != "" {
			badgeNames = append(badgeNames, strings.ToLower(badge.Name))
		}
	}

	return badgeNames, nil
}

// extractBadgesFromClaims extracts badges from JWT claims
func (m *Middleware) extractBadgesFromClaims(claims jwt.MapClaims) []string {
	var badgeNames []string

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

	return badgeNames
}

// RequireRoleOrBadge creates a middleware to check user roles OR badges
func (m *Middleware) RequireRoleOrBadge(allowedRoles string, allowedBadges string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		userBadgesVal, _ := c.Get("badges")

		// Convert userBadgesVal to []string
		userBadges := m.convertToStringSlice(userBadgesVal)

		// Check roles
		roleAllowed := m.checkRoles(userRole, allowedRoles)

		// Check badges
		badgeAllowed := m.checkBadges(userBadges, allowedBadges)

		if !roleAllowed && !badgeAllowed {
			m.handleForbiddenError(c, userRole, userBadges, allowedRoles, allowedBadges)
			return
		}

		c.Next()
	}
}

// checkRoles checks if the user's role matches any of the allowed roles
func (m *Middleware) checkRoles(userRole interface{}, allowedRoles string) bool {
	if allowedRoles == "" {
		return false
	}

	roles := strings.Split(allowedRoles, ",")
	for _, role := range roles {
		if userRole == strings.TrimSpace(role) {
			return true
		}
	}

	return false
}

// checkBadges checks if the user has any of the allowed badges
func (m *Middleware) checkBadges(userBadges []string, allowedBadges string) bool {
	if allowedBadges == "" || len(userBadges) == 0 {
		return false
	}

	allowedBadgeList := strings.Split(allowedBadges, ",")
	for _, allowedBadge := range allowedBadgeList {
		allowedBadge = strings.ToLower(strings.TrimSpace(allowedBadge))
		for _, userBadge := range userBadges {
			if userBadge == allowedBadge {
				return true
			}
		}
	}

	return false
}

// handleForbiddenError handles forbidden access with detailed error response
func (m *Middleware) handleForbiddenError(
	c *gin.Context,
	userRole interface{},
	userBadges []string,
	allowedRoles,
	allowedBadges string,
) {
	// Create custom error with detailed context
	errors.New(
		errors.ErrAuthorization,
		"Insufficient permissions",
		nil,
		errors.WithContext("required_roles", strings.Split(allowedRoles, ",")),
		errors.WithContext("required_badges", strings.Split(allowedBadges, ",")),
		errors.WithContext("user_role", userRole),
		errors.WithContext("user_badges", userBadges),
	)

	// Log the unauthorized access attempt
	m.logger.Info("Access denied",
		"roles", allowedRoles,
		"badges", allowedBadges,
		"user_role", userRole,
		"user_badges", userBadges,
	)

	// Send forbidden response
	response.Forbidden(c, "insufficient_permissions", "Insufficient permissions", "")
	c.Abort()
}

// convertToStringSlice converts badges context value to []string
func (m *Middleware) convertToStringSlice(userBadgesVal interface{}) []string {
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
	return userBadges
}

// GetUserFromContext retrieves user information from Gin context
func GetUserFromContext(c *gin.Context) (userID, email, role string, badges []string, err error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return "", "", "", nil, errors.New(
			errors.ErrAuthentication,
			"User ID not found in context",
			nil,
			errors.WithContext("context_key", "user_id"),
		)
	}

	emailVal, exists := c.Get("email")
	if !exists {
		return "", "", "", nil, errors.New(
			errors.ErrAuthentication,
			"Email not found in context",
			nil,
			errors.WithContext("context_key", "email"),
		)
	}

	roleVal, exists := c.Get("role")
	if !exists {
		return "", "", "", nil, errors.New(
			errors.ErrAuthentication,
			"Role not found in context",
			nil,
			errors.WithContext("context_key", "role"),
		)
	}

	badgesVal, _ := c.Get("badges")
	badgeList := convertToStringSlice(badgesVal)

	return userIDVal.(string), emailVal.(string), roleVal.(string), badgeList, nil
}

// convertToStringSlice is a standalone function to convert badges to string slice
func convertToStringSlice(userBadgesVal interface{}) []string {
	var badgeList []string
	switch v := userBadgesVal.(type) {
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
	return badgeList
}
