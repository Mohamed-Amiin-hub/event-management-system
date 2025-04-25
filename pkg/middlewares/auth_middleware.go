package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the token and protects the routes by checking the token in the database.
func AuthMiddleware(tokenRepo repository.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// The token is usually in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Invalid Authorization format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Fetch the token from the database
		token, err := tokenRepo.FindByToken(tokenString)
		if err != nil {
			log.Printf("Token lookup failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if the token has expired
		if token.ExpiresAt.Before(time.Now()) {
			log.Printf("Token expired at: %v", token.ExpiresAt)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// If the token is valid, set the user ID in the request context
		c.Set("userID", token.UserID)

		// Proceed to the next handler in the chain
		c.Next()
	}
}
