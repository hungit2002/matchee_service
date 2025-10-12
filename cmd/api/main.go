package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"matchee/services/internal/config"
	"matchee/services/internal/router"
	"matchee/services/pkg/database"
	"matchee/services/pkg/logger"
	"matchee/services/pkg/rabbitmq"
	"matchee/services/pkg/redis"
)

func main() {
	logger := logger.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("load config: %v", err)
	}

	// Initialize infrastructure
	mysqlDB, err := database.NewMySQL(cfg)
	if err != nil {
		logger.Fatalf("mysql: %v", err)
	}

	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		logger.Fatalf("redis: %v", err)
	}

	amqpConn, amqpCh, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		logger.Fatalf("rabbitmq: %v", err)
	}

	// Build HTTP server
	e := router.BuildHTTPRouter(cfg, logger, mysqlDB, redisClient, amqpCh)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr(),
		Handler:      e,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Infof("HTTP server listening on %s", cfg.HTTPAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown: %v", err)
	}

	// Close resources
	if err := amqpCh.Close(); err != nil {
		logger.Warnf("amqp channel close: %v", err)
	}
	if err := amqpConn.Close(); err != nil {
		logger.Warnf("amqp connection close: %v", err)
	}
	if err := redisClient.Close(); err != nil {
		logger.Warnf("redis close: %v", err)
	}
	if sqlDB, err := mysqlDB.DB(); err == nil {
		_ = sqlDB.Close()
	}

	logger.Infof("server exited")
}
