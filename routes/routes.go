package routes

import (
	"project/handlers"
	"project/middleware"

	"github.com/gin-contrib/cors" // CORS middleware for Gin
	"github.com/gin-gonic/gin"    // Gin web framework
)

// SetupRouter sets up the routes for the application.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS for all routes using the gin-contrib/cors middleware.
	r.Use(cors.Default())

	// Authentication routes.
	r.POST("/signup", handlers.SignUpHandler) // New signup route.
	r.POST("/login", handlers.LoginHandler)

	// Test route.
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	// You can remove or repurpose this route if /signup is your sign-up endpoint.
	r.POST("/users", handlers.CreateUser)

	// Protected route to get all users.
	r.GET("/users", middleware.AuthMiddleware(), handlers.GetAllUsers)

	// Additional routes.
	r.GET("/export-users/excel", handlers.ExportUsersToExcel)
	r.GET("/export-users/pdf", handlers.ExportUsersToPDF)
	r.GET("/users/between-dates", handlers.GetUsersBetweenDates)

	return r
}
