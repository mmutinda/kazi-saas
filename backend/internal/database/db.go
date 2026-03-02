package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kazi/backend/internal/config"
	"kazi/backend/internal/models"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("Database connected")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.Admin{},
		&models.Job{},
		&models.ScrapeLog{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Add FULLTEXT index for search (MySQL specific)
	DB.Exec(`ALTER TABLE jobs ADD FULLTEXT INDEX ft_search (title, company, description)`)

	log.Println("Database migrated")
}

func SeedAdmin(cfg *config.Config, passwordHash string) {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)
	if count > 0 {
		return
	}

	admin := models.Admin{
		Email:        cfg.AdminEmail,
		PasswordHash: passwordHash,
		Name:         cfg.AdminName,
	}
	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin: %v", err)
		return
	}
	log.Printf("Admin seeded: %s", cfg.AdminEmail)
}
