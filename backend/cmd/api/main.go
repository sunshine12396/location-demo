package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/example/location-demo/internal/config"
	"github.com/example/location-demo/internal/external/google"
	"github.com/example/location-demo/internal/location"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Connect to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying SQL DB: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	// 3. Wire dependencies (Dependency Injection)
	repo := location.NewPostgresRepository(db)

	// Create the external provider using Google Places API
	extApi := google.NewGooglePlacesClient(cfg.ExternalAPIKey)

	svc := location.NewService(repo, extApi, cfg.LocationSyncDays)
	handler := location.NewHandler(svc)

	// 4. Setup Gin router
	r := gin.Default()

	// CORS middleware for frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register location routes
	handler.RegisterRoutes(r)

	// 5. Start server
	addr := ":" + cfg.ServerPort
	log.Printf("🚀 Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
