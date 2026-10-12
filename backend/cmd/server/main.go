package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"kazi/backend/internal/config"
	"kazi/backend/internal/database"
	"kazi/backend/internal/handlers"
	"kazi/backend/internal/middleware"
	"kazi/backend/internal/scraper"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate()

	// Hash and seed admin password
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}
	database.SeedAdmin(cfg, string(hash))

	scraper.StartScheduler()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World"})
	})

	// Public routes
	jobsHandler := &handlers.JobsHandler{}
	r.GET("/api/jobs", jobsHandler.ListJobs)
	r.GET("/api/jobs/:id", jobsHandler.GetJob)

	// Admin routes
	adminHandler := &handlers.AdminHandler{JWTSecret: cfg.JWTSecret}
	r.POST("/api/admin/login", adminHandler.Login)

	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		admin.GET("/jobs", adminHandler.ListJobs)
		admin.POST("/jobs", adminHandler.CreateJob)
		admin.PUT("/jobs/approve-all", adminHandler.ApproveAllPending)
		admin.GET("/jobs/:id", adminHandler.GetJob)
		admin.PUT("/jobs/:id", adminHandler.UpdateJob)
		admin.PUT("/jobs/:id/approve", adminHandler.ApproveJob)
		admin.PUT("/jobs/:id/reject", adminHandler.RejectJob)
		admin.DELETE("/jobs/:id", adminHandler.DeleteJob)
		admin.GET("/stats", adminHandler.Stats)
		admin.POST("/scrape", adminHandler.TriggerScrape)
		admin.GET("/scrape/logs", adminHandler.ScrapeLogs)
	}

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
