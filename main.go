package main

import (
	"log"
	"os"

	"npm-download-spammer/pkg/config"
	"npm-download-spammer/pkg/spammer"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration loading error: %v", err)
	}

	// If configuration was not set through environment variables or file,
	// request it through CLI
	if len(cfg.GetPackageNames()) == 0 {
		cliConfig, err := config.GetConfigFromCLI()
		if err != nil {
			log.Fatalf("CLI configuration error: %v", err)
		}
		cfg = cliConfig
	}

	// Run main logic
	if err := spammer.Run(cfg); err != nil {
		log.Fatalf("Execution error: %v", err)
		os.Exit(1)
	}
}
