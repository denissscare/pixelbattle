package main

import (
	"pixelbattle/internal/config"
	"pixelbattle/internal/server"
	"pixelbattle/pkg/logger"
)

func main() {
	config := config.LoadConfig()
	log := logger.New(config.Enviroment)

	router := server.InitRouter()
	srv := server.New(config, router, log)

	if err := srv.Run(); err != nil {
		panic("failed start")
	}
}
