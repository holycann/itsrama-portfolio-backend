package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/holycann/cultour-backend/internal/response"
	"github.com/supabase-community/auth-go"
)

// Middleware handles JWT token authentication and validation
type Middleware struct {
	supabaseAuth auth.Client
	jwks         *keyfunc.JWKS
}

// UserContext represents the authenticated user's context
type UserContext struct {
	ID    string
	Email string
	Role  string
	Badge string
}

// NewMiddleware creates a new JWT middleware instance
func NewMiddleware(supabaseAuth auth.Client, jwks *keyfunc.JWKS) *Middleware {
	return &Middleware{
		supabaseAuth: supabaseAuth,
		jwks:         jwks,
	}
}

// VerifyJWT validates the JWT token from the Authorization header
func (m *Middleware) VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization token", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(c, "Invalid token format", nil)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, m.jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)
		badge, _ := claims["badge"].(string)

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)
		c.Set("badge", badge)

		c.Next()
	}
}

// RequireRole creates a middleware to check user roles
func (m *Middleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// userRole, exists := c.Get("role")
		// if !exists {
		// 	response.Forbidden(c, "User role not found", nil)
		// 	c.Abort()
		// 	return
		// }

		// roleAllowed := true
		// for _, role := range allowedRoles {
		// 	if userRole == role {
		// 		roleAllowed = true
		// 		break
		// 	}
		// }

		// userRole = "admin"

		// if !roleAllowed {
		// 	response.Forbidden(c, "Insufficient permissions", gin.H{
		// 		"required_roles": allowedRoles,
		// 		"user_role":      userRole,
		// 	})
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}

// RequireRoleOrBadge creates a middleware to check user roles or badges
func (m *Middleware) RequireRoleOrBadge(allowedRoles string, allowedBadges string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		userBadge, _ := c.Get("badge")

		roles := strings.Split(allowedRoles, ",")
		badges := strings.Split(allowedBadges, ",")

		roleAllowed := true
		badgeAllowed := true

		for _, role := range roles {
			if userRole == strings.TrimSpace(role) {
				roleAllowed = true
				break
			}
		}
		for _, badge := range badges {
			if userBadge == strings.TrimSpace(badge) {
				badgeAllowed = true
				break
			}
		}

		if !roleAllowed && !badgeAllowed {
			response.Forbidden(c, "Insufficient permissions", gin.H{
				"required_roles":  roles,
				"required_badges": badges,
				"user_role":       userRole,
				"user_badge":      userBadge,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext retrieves user information from Gin context
func GetUserFromContext(c *gin.Context) (userID, email, role, badge string, err error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return "", "", "", "", errors.New("user_id not found in context")
	}

	emailVal, exists := c.Get("email")
	if !exists {
		return "", "", "", "", errors.New("email not found in context")
	}

	roleVal, exists := c.Get("role")
	if !exists {
		return "", "", "", "", errors.New("role not found in context")
	}

	badgeVal, _ := c.Get("badge")

	return userIDVal.(string), emailVal.(string), roleVal.(string), badgeVal.(string), nil
}
