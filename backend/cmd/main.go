package main

import (
	"log"
	"net/http"
	"strings"
	"video-conference-sdk/backend/api"
	"video-conference-sdk/backend/config"
	"video-conference-sdk/backend/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Initialize database
	db.InitPostgres()
	
	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     strings.Split(cfg.AllowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "video-conference-sdk",
		})
	})

	// Public routes
	r.POST("/auth/signup", api.SignupHandler)
	r.POST("/auth/login", api.LoginHandler)

	// Protected routes
	protected := r.Group("/")
	protected.Use(api.AuthMiddleware())
	{
		protected.POST("/room/create", api.CreateRoomHandler)
		protected.POST("/room/join", api.JoinQueueHandler)
	}

	// WebSocket signaling (no auth for simplicity, but can be added)
	r.GET("/ws/signaling", api.SignalingHandler)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}