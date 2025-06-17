package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	auth "pixelbattle/internal/auth/service"
	storage "pixelbattle/internal/auth/storage/postgres"
	"pixelbattle/internal/config"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/metrics"
	"pixelbattle/internal/pixcelbattle/service"
	pgstorage "pixelbattle/internal/pixcelbattle/storage/postgres"
	"pixelbattle/internal/pixcelbattle/storage/redis"
	"pixelbattle/internal/s3"
	"pixelbattle/internal/server"
	"pixelbattle/internal/storage/postgres"
	jwtutil "pixelbattle/pkg/jwt"
	"pixelbattle/pkg/logger"
	"time"
)

func main() {
	ctx := context.Background()
	config := config.LoadConfig()
	log := logger.New(config.Environment)
	metrics := metrics.NewPrometheusMetrics()

	minioClient, err := s3.New(*config)
	if err != nil {
		log.Fatalf("Minio init failed: %v", err)
	}

	postgresDB, err := postgres.NewStorage(*config)
	if err != nil {
		log.Fatalf("Postgres init failed: %v", err)
	}
	defer postgresDB.Close()

	userStorage := storage.NewRepository(postgresDB)
	pixelHistoryStorage := pgstorage.NewRepository(postgresDB)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretkey"
	}
	jwtManager := jwtutil.NewManager(jwtSecret, 24*time.Hour)

	rds, err := redis.NewClient(ctx, *config)
	if err != nil {
		log.Fatalf("Redis init failed: %v", err)
	}
	defer rds.Close(ctx)

	br, err := broker.NewBroker(config.NATS.URL)
	if err != nil {
		log.Fatalf("NATS init failed: %v", err)
	}
	defer br.Close()

	pixelbattle := service.NewBattleService(rds, pixelHistoryStorage, br, log, metrics)
	auth := auth.NewService(userStorage, jwtManager, log, minioClient)
	router := server.InitRouter(pixelbattle, auth, log, metrics, jwtManager, minioClient, *config)

	go func() {
		fmt.Println("pprof listening on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()

	srv := server.New(config, router, log)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
