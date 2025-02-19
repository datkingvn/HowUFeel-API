package middlewares

import (
	"HowUFeel-API-Prj/helpers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Remove Bearer <JWT_TOKEN>
		token := strings.TrimPrefix(authHeader, "Bearer")
		token = strings.TrimSpace(token)

		claims, err := helpers.VerifyToken(token)

		if err != nil {
			log.Println("Token validation error: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
