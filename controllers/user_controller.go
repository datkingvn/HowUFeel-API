package controllers

import (
	"HowUFeel-API-Prj/helpers"
	"HowUFeel-API-Prj/models"
	"HowUFeel-API-Prj/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Response helper
func handleErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser, err := services.RegisterUser(&user)
		if err != nil {
			handleErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": newUser})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			handleErrorResponse(c, http.StatusBadRequest, "Invalid input")
			return
		}

		user, accessToken, refreshToken, err := services.LoginUser(&loginRequest)
		if err != nil {
			handleErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":         user,
			"token":        accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			handleErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenClaims, ok := claims.(*helpers.Claims)
		if !ok || tokenClaims.Role != "ADMIN" {
			handleErrorResponse(c, http.StatusForbidden, "Forbidden: only admins can fetch users")
			return
		}

		users, err := services.GetAllUsers(tokenClaims.Role)
		if err != nil {
			handleErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			handleErrorResponse(c, http.StatusUnauthorized, "Unauthorized: missing token")
			return
		}

		tokenClaims, ok := claims.(*helpers.Claims)
		if !ok {
			handleErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID := c.Param("id")
		log.Println("Fetching user with ID:", userID)

		user, err := services.GetUserByID(userID, tokenClaims.UserID, tokenClaims.Role)
		if err != nil {
			handleErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
