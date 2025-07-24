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
}

// NewJWTMiddleware creates a new JWT middleware instance
func NewMiddleware(supabaseAuth auth.Client, jwks *keyfunc.JWKS) *Middleware {
	return &Middleware{
		supabaseAuth: supabaseAuth,
		jwks:         jwks,
	}
}

// Verify validates the JWT token from the Authorization header
func (m *Middleware) VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization token", nil)
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(c, "Invalid token format", nil)
			c.Abort()
			return
		}

		// Verifikasi token pakai JWK dari Supabase
		token, err := jwt.Parse(tokenString, m.jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// Ambil data dari claims
		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		// Set user context untuk dipakai di handler selanjutnya
		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}

// RequireRole creates a middleware to check user roles
func (m *Middleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user role from context
		userRole, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "User role not found", nil)
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		roleAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			response.Forbidden(c, "Insufficient permissions", gin.H{
				"required_roles": allowedRoles,
				"user_role":      userRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext retrieves user information from Gin context
func GetUserFromContext(c *gin.Context) (userID, email, role string, err error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return "", "", "", errors.New("user_id not found in context")
	}

	emailVal, exists := c.Get("email")
	if !exists {
		return "", "", "", errors.New("email not found in context")
	}

	roleVal, exists := c.Get("role")
	if !exists {
		return "", "", "", errors.New("role not found in context")
	}

	return userIDVal.(string), emailVal.(string), roleVal.(string), nil
}
