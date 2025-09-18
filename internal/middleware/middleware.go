package middleware

import (
	"slices"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/holycann/itsrama-portfolio-backend/internal/response"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/logger"
)

// Middleware handles JWT token authentication and validation
type Middleware struct {
	jwks          *keyfunc.JWKS
	allowedEmails []string
	logger        *logger.Logger
}

// NewMiddleware creates a new JWT middleware instance
func NewMiddleware(
	jwks *keyfunc.JWKS,
	allowedEmails []string,
	logger *logger.Logger,
) *Middleware {
	return &Middleware{
		jwks:          jwks,
		allowedEmails: allowedEmails,
		logger:        logger,
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

		if isAllowed := slices.Contains(m.allowedEmails, email); !isAllowed {
			m.handleAuthError(c, "Unauthorized email",
				errors.WithContext("email", email))
			return
		}

		// Set user context
		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}

// handleAuthError handles authentication errors with standardized response
func (m *Middleware) handleAuthError(c *gin.Context, message string, opts ...func(*errors.CustomError)) {
	err := errors.New(
		errors.ErrAuthentication,
		message,
		nil,
		opts...,
	)
	if err != nil {
		// Log the authentication error
		m.logger.Error("Authentication error", err)
	}
	response.Unauthorized(c, "auth_error", message, "")
	c.Abort()
}
