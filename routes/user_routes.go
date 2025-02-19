package routes

import (
	"HowUFeel-API-Prj/controllers"
	"HowUFeel-API-Prj/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/v1/users")

	userGroup.POST("/register", controllers.Register())
	userGroup.POST("/login", controllers.Login())

	protected := userGroup.Group("/")

	protected.Use(middlewares.Authenticate())

	{
		protected.GET("/", controllers.GetUsers())
		protected.GET("/:id", controllers.GetUser())
	}
}
