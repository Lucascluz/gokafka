package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucas/gokafka/api-gateway/internal/cache"
	"github.com/lucas/gokafka/shared/auth"
)

type AuthMiddleware struct {
	jwtBlacklist *cache.TokenBlacklist
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		jwtBlacklist: cache.NewTokenBlacklist(),
	}
}

func (am *AuthMiddleware) AuthMiddleware(blackListCheck bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		log.Printf("Extracted token: %s", tokenString)

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is blacklisted
		if blackListCheck {
			if am.jwtBlacklist.IsTokenBlacklisted(claims.ID) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Revoked token"})
				c.Abort()
				return
			}
		}

		log.Printf("Token validated successfully for user: %s", claims.UserID)

		// Set user info in context for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token_id", claims.ID)
		if claims.ExpiresAt != nil {
			c.Set("token_exp", claims.ExpiresAt.Unix())
		}

		c.Next()
	}
}

func (am *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if userRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
