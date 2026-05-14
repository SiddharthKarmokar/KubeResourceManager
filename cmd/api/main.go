package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SiddharthKarmokar/KubeResourceManager/internal/api/handlers"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/api/router"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/config"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/recommendation"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/engine/scoring"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/logger"
	"github.com/SiddharthKarmokar/KubeResourceManager/internal/services"

	_ "github.com/SiddharthKarmokar/KubeResourceManager/swagger" // register OpenAPI spec with swag
)

// @title Kubernetes Resource Optimization API
// @version 1.0
// @description Analyzes Kubernetes workload metrics and recommends safe CPU and memory requests using conservative heuristics, safety buffers, and guardrails.
// @termsOfService http://swagger.io/terms/

// @contact.name Siddharh Karmokar
// @contact.url https://siddkarmokar-portfolio.vercel.app/
// @contact.email siddkarmokar@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	logger.Init()
	log := slog.Default()

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	scorer := scoring.NewScorer()
	recommender := recommendation.NewRecommender(cfg.Recommendation, scorer)
	optimizer := services.NewOptimizerService(recommender)

	h := handlers.NewHandler(optimizer)
	r := router.Setup(h)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}
