package main

import (
	"log"
	"os"

	"goldprice-api/handlers"
	"goldprice-api/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Firebase
	if err := services.InitFirebase("firebase-credentials.json"); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer services.CloseFirebase()

	// Create Gin router
	router := gin.Default()

	// CORS Configuration - Public access for now
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	router.Use(cors.New(config))

	// Routes
	router.GET("/", handlers.HealthCheck)
	router.POST("/api/upload", handlers.UploadCSV)
	router.GET("/api/prices", handlers.GetAllPrices)
	router.GET("/api/prices/latest", handlers.GetLatestPrices)
	router.GET("/api/prices/:date", handlers.GetPriceByDate)
	router.GET("/api/metadata", handlers.GetMetadata)
	router.DELETE("/api/prices/clear", handlers.ClearAllData)

	// Get port from environment variable (Cloud Run sets PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for Cloud Run
	}

	// Start server
	log.Printf("🚀 Gold Price API starting on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
