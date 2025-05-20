package main

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"backend/scheduler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDatabase()

	// Run database migrations
	if err := config.RunMigrations(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	scheduler.StartWeatherScheduler()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/login", handlers.Login)
	r.GET("/validate", middleware.RequireAuth, handlers.Validate)
	r.POST("/locations", middleware.RequireAuth, handlers.SaveLocation)
	r.GET("/locations", middleware.RequireAuth, handlers.GetLocations)
	r.DELETE("/locations/:id", middleware.RequireAuth, handlers.DeleteLocation)
	r.PUT("/locations/:id", middleware.RequireAuth, handlers.UpdateLocationName)
	r.GET("/predictions", middleware.RequireAuth, handlers.GetPredictions)
	r.GET("/locations/:id/predictions", middleware.RequireAuth, handlers.GetPredictionsByLocation)

	if os.Getenv("JWT_KEY") == "" {
		os.Setenv("JWT_KEY", "SECRET")
	}

	r.Run(":8080")
}
