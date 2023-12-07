package main

import (
	"log/slog"
	"os"

	"github.com/EventWhisper/EventWhisper/whisper-core/internal/api"
	"github.com/EventWhisper/EventWhisper/whisper-core/internal/configuration"
)

func main() {

	logOpts := slog.HandlerOptions{
		Level: configuration.GetSlogLevel(os.Getenv("LOG_LEVEL")),
	}

	textHandler := slog.NewTextHandler(os.Stdout, &logOpts)
	logger := slog.New(textHandler)
	slog.SetDefault(logger)

	slog.Info("Starting EventWhisper")

	config, err := configuration.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	server, err := api.NewServer(config)
	if err != nil {
		slog.Error("Error creating server", "error", err)
		os.Exit(1)
	}

	err = server.Start()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
