package main

import (
	"net/http"

	"github.com/zhxauda9/StayMate/internal/config"
	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/server"
)

const address = "127.0.0.1:8080"

func main() {
	config.LoadEnvVariables()
	logger := myLogger.NewZeroLogger()

	// В файле main.go
	logger.Info().Msg("Starting the application")

	mux, err := server.InitServer()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error initializing server")
	}

	logger.Info().Msg("Server initialized successfully, starting HTTP server...")
	if err := http.ListenAndServe(address, mux); err != nil {
		logger.Fatal().Err(err).Msg("Error starting HTTP server")
	}
}
