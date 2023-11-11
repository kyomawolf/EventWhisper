package main

import (
	"os"
	"strings"

	"github.com/kyomawolf/EventWhisper/whisper-core/internal/api"
	"github.com/kyomawolf/EventWhisper/whisper-core/internal/configuration"
	log "github.com/sirupsen/logrus"
)

func main() {

	config, err := configuration.LoadConfig()
	if err != nil {
		log.Errorf("Error loading config: %v", err)
		os.Exit(1)
	}

	log.SetLevel(log.WarnLevel)

	if strings.ToLower(config.LogLevel) == "debug" {
		log.SetLevel(log.DebugLevel)
	}
	if strings.ToLower(config.LogLevel) == "info" {
		log.SetLevel(log.InfoLevel)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Errorf("Error creating server: %v", err)
		os.Exit(1)
	}

	err = server.ConfigureRouter()
	if err != nil {
		log.Errorf("Error configuring router: %v", err)
		os.Exit(1)
	}

	err = server.Start()
	if err != nil {
		log.Errorf("Error starting server: %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
