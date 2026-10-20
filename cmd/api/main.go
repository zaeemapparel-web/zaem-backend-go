package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/config"
	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/logger"
	"github.com/zaem/backend-go/internal/middleware"
	"github.com/zaem/backend-go/internal/routes"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

func main() {
	// 1. Load config
	config.Load()

	// 2. Init logger
	logger.Init()

	// 3. Init validator
	customvalidator.Init()

	log.Info().Msg("🚀 Starting ZAEM Backend (Go + Gin)...")

	// 4. Connect database
	database.Connect()
	defer database.Close()

	// 5. Auto-migrate tables
	database.AutoMigrate()

	// 6. Gin setup
	if config.C.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 7. Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimitGeneral())

	// 8. Health check routes (Frontend/Backend connection test)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"brand":     "ZAEM",
			"tagline":   "Style. Redefined.",
			"status":    "API is running (Go)",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	router.GET("/api/health", func(c *gin.Context) {
		dbStatus := "ok"
		if err := database.Health(); err != nil {
			dbStatus = "error: " + err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"engine":    "Go + Gin",
			"uptime":    time.Since(startTime).String(),
			"services":  gin.H{"database": dbStatus},
			"timestamp": time.Now().UTC(),
		})
	})

	// 9. Register all routes
	routes.Register(router)

	// 10. HTTP Server
	srv := &http.Server{
		Addr:         ":" + config.C.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("📍 Server running on port %s", config.C.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("❌ Server failed")
		}
	}()

	// 11. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("❌ Forced shutdown")
	}

	log.Info().Msg("✅ Server exited cleanly")
}

var startTime = time.Now()