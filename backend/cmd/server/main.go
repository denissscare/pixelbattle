package main

import (
	"context"
	"pixelbattle/internal/config"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/internal/pixcelbattle/storage/redis"
	"pixelbattle/internal/server"
	"pixelbattle/pkg/logger"
)

func main() {
	ctx := context.Background()
	config := config.LoadConfig()
	log := logger.New(config.Environment)

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

	pixelbattle := service.NewBattleService(*rds, *br, log)

	router := server.InitRouter(pixelbattle, log)

	srv := server.New(config, router, log)

	if err := srv.Run(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
