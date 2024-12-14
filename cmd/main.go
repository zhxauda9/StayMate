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

	logger.Info().Msg("Starting application")

	mux, err := server.InitServer()
	if err != nil {
		logger.Fatal().Str("Error", err.Error()).Msg("Error inizializing server")
	}

	logger.Info().Str("Address", address).Msg("Starting server!")
	http.ListenAndServe(address, mux)
}
