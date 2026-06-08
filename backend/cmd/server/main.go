package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"sales-module/internal/config"
	"sales-module/internal/database"
	"sales-module/internal/middleware"
	"sales-module/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set JWT secret
	middleware.JWTSecret = cfg.JWTSecret

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	// Setup Gin
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Register routes
	routes.Setup(r, db)

	// Start server
	log.Printf(" Sales Module API running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
