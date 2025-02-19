package main

import (
	"HowUFeel-API-Prj/configs"
	"HowUFeel-API-Prj/helpers"
	"HowUFeel-API-Prj/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	secret := configs.GenerateRandomKey()
	helpers.SetJWTSecretKey(secret)

	r := gin.Default()
	
	routes.UserRoutes(r)

	// Start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		log.Println("PORT environment variable not found, defaulting to 8080")
		port = "8080"
	}
	r.Run(":" + port)
	log.Println("Sever running on port: ", port)
}
