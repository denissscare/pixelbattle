package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"pixelbattle/internal/config"
	"pixelbattle/internal/pixcelbattle/broker"
	"pixelbattle/internal/pixcelbattle/metrics"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/internal/pixcelbattle/storage/redis"
	"pixelbattle/internal/server"
	"pixelbattle/pkg/logger"
)

func main() {
	ctx := context.Background()
	config := config.LoadConfig()
	log := logger.New(config.Environment)
	metrics := metrics.NewPrometheusMetrics()

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

	pixelbattle := service.NewBattleService(*rds, *br, log, metrics)

	router := server.InitRouter(pixelbattle, log, metrics)

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
