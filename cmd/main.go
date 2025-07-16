package main

import (
	"log/slog"
	"os"

	"github.com/ashtonx86/mocker/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	isProd := os.Getenv("PRODUCTION")
	if isProd != "true" {
		slog.Info("Starting application in development mode...")

		if err := godotenv.Load(".env.local"); err != nil {
			panic(err)
		}
	}

	server := server.NewWebServer()
	if server == nil {
		return
	}

	server.App.Listen(":3000")
}
