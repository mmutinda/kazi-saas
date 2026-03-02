package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string

	AdminEmail    string
	AdminPassword string
	AdminName     string

	Port string
}

func Load() *Config {
	// Overload forces .env values to win over existing shell env vars.
	// This prevents a stale DB_HOST=localhost:3306 in the shell from
	// corrupting the DSN into tcp(localhost:3306:3306).
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "kazi"),

		JWTSecret: getEnv("JWT_SECRET", "fallback-secret-key"),

		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@kazi.co.ke"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "Admin@123"),
		AdminName:     getEnv("ADMIN_NAME", "Kazi Admin"),

		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
