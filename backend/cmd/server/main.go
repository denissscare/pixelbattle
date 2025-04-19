package main

import (
	"pixelbattle/internal/config"
	"pixelbattle/internal/server"
)

func main() {
	config := config.LoadConfig()

	router := server.InitRouter()
	srv := server.New(config, router)

	if err := srv.Run(); err != nil {
		panic("failed start")
	}
}
