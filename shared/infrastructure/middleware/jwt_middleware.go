package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	secretKey []byte
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ProfileID string `json:"profileId"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(secretKey string) *JWTMiddleware {
	return &JWTMiddleware{
		secretKey: []byte(secretKey),
	}
}

// AuthMiddleware validates JWT token from Authorization header
func (jm *JWTMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		claims, err := jm.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set claims in context for use in handlers
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("profileID", claims.ProfileID)

		log.Printf("Authenticated user: %s (UserID: %s)", claims.Email, claims.UserID)
		c.Next()
	}
}

// ValidateToken validates and parses JWT token
func (jm *JWTMiddleware) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jm.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// GetUserIDFromContext extracts userID from gin context
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.New("userID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("userID is not a string")
	}

	return userIDStr, nil
}

// GetProfileIDFromContext extracts profileID from gin context
func GetProfileIDFromContext(c *gin.Context) (string, error) {
	profileID, exists := c.Get("profileID")
	if !exists {
		return "", errors.New("profileID not found in context")
	}

	profileIDStr, ok := profileID.(string)
	if !ok {
		return "", errors.New("profileID is not a string")
	}

	return profileIDStr, nil
}

// GetRoleFromContext extracts role from gin context
func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role not found in context")
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", errors.New("role is not a string")
	}

	return roleStr, nil
}

// RequireRole creates middleware that checks if user has required role
func (jm *JWTMiddleware) RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, required := range requiredRoles {
			if role == required {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
