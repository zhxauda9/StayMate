package config

import (
	"log"

	"github.com/joho/godotenv"
)

var AvailiableMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
