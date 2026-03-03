package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hakathon-mvp/internal/adapters/kafka"
	"hakathon-mvp/internal/adapters/postgres"
	"hakathon-mvp/internal/adapters/redis"
	"hakathon-mvp/internal/domain/services"
	v1 "hakathon-mvp/internal/handlers/http/v1"
	"hakathon-mvp/internal/pkg/cache"
	"hakathon-mvp/internal/pkg/config"
	"hakathon-mvp/internal/pkg/database"
	"hakathon-mvp/internal/pkg/logger"
	"hakathon-mvp/internal/pkg/metrics"
	vld "hakathon-mvp/internal/pkg/validator"
	"hakathon-mvp/internal/usecases"

	redis2 "github.com/redis/go-redis/v9"
)

func main() {
	// conf
	cfg := config.Load()

	// logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// metrics
	metrics.Init(cfg.Metrics.Port)

	// psql
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Fatal(context.Background(), "failed to connect to database", "error", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(context.Background(), "failed to close database connection", "error", err)
		}
	}(db)

	// redis
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatal(context.Background(), "failed to connect to redis", "error", err)
	}
	defer func(redisClient *redis2.Client) {
		err := redisClient.Close()
		if err != nil {
			logger.Error(context.Background(), "failed to close redis client", "error", err)
		}
	}(redisClient)

	// kafka producer
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer func(producer *kafka.Producer) {
		err := producer.Close()
		if err != nil {
			logger.Error(context.Background(), "failed to close producer", "error", err)
		}
	}(producer)

	// reps and services
	citizenRepo := postgres.NewCitizenReportRepository(db)
	citizenCache := redis.NewCitizenReportCache(redisClient, cfg.Redis.TTL)
	validator := services.NewCitizenReportValidator()

	// usecases
	citizenUC := usecases.NewCitizenReportUseCase(citizenRepo, citizenCache, producer, (*vld.CitizenReportValidator)(validator))

	// http server
	router := v1.NewRouter(citizenUC, db, redisClient, cfg.Server.RateLimit)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info(context.Background(), "starting API server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(context.Background(), "failed to start server", "error", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "server forced to shutdown", "error", err)
	}

	logger.Info(context.Background(), "server exited")
}
