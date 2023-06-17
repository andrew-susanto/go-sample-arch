package goenv

import (
	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"github.com/joho/godotenv"
)

// InitEnvironmentVariable initialize environment variable
func InitEnvironmentVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Error(err, nil, "godotenv.Load() got error - InitEnvironmentVariable")
	}
}
