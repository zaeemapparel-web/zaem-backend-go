package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all environment variables
type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret    string
	JWTExpiresIn time.Duration

	// Frontend
	FrontendURL string

	// Email
	ResendAPIKey string
	EmailFrom    string
	EmailReplyTo string

	// AI
	GroqAPIKey string
}

// C is the global config instance (accessible everywhere)
var C *Config

// Load reads .env and populates the config
func Load() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	C = &Config{
		Port: getEnv("PORT", "5000"),
		Env:  getEnv("ENV", "development"),

		DatabaseURL: mustEnv("DATABASE_URL"),

		JWTSecret:    mustEnv("JWT_SECRET"),
		JWTExpiresIn: parseDuration(getEnv("JWT_EXPIRES_IN", "30d")),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		EmailFrom:    getEnv("EMAIL_FROM", ""),
		EmailReplyTo: getEnv("EMAIL_REPLY_TO", ""),

		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
	}

	log.Printf("✅ Config loaded (env: %s, port: %s)", C.Env, C.Port)
}

// getEnv returns env var or fallback if empty
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// mustEnv fails if env var is missing
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Missing required environment variable: %s", key)
	}
	return value
}

// parseDuration converts "30d" to time.Duration
func parseDuration(s string) time.Duration {
	// Handle "Nd" format (days)
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return 30 * 24 * time.Hour
		}
		return time.Duration(days) * 24 * time.Hour
	}

	// Standard duration (1h, 30m, 60s, etc.)
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * 24 * time.Hour
	}
	return d
}

// IsDevelopment returns true if env is development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if env is production
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}