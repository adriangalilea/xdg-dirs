package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	
	"xdg-dirs/internal/conf"
	"xdg-dirs/internal/logger"
	"xdg-dirs/internal/setup"
	"xdg-dirs/internal/updater"
)

var log *logger.Logger

func main() {
	// Parse command-line flags
	debug := flag.Bool("d", false, "Enable debug output")
	dryRun := flag.Bool("n", false, "Simulate changes without applying them")
	createDirs := flag.Bool("c", false, "Create directories if they don't exist")
	logFilePath := flag.String("l", conf.DefaultLogFilePath, "Specify the log file path")
	help := flag.Bool("help", false, "Show help message")
	flag.BoolVar(help, "h", false, "Show help message")
	flag.Parse()

	// Display help message if requested
	if *help {
		fmt.Println(conf.HelpMessage)
		os.Exit(0)
	}

	log = logger.NewLogger(*debug, *logFilePath)

	// Perform initial setup
	if err := setup.Prepare(log); err != nil {
		log.Fatal("Failed to perform initial setup: %v", err)
	}

	// Create updater instance
	updaterInstance := updater.NewUpdater(log)

	// Get user directories
	userDirs, err := updaterInstance.GetUserDirs()
	if err != nil {
		log.Fatal("Failed to get user directories: %v", err)
	}

	// Update user directories
	if err := updaterInstance.Update(userDirs, *createDirs, *dryRun); err != nil {
		log.Fatal("Failed to update user directories: %v", err)
	}

	// Get the EXPORT env variables
	exports := updaterInstance.ExportEnv(userDirs)

	log.Debug("XDG environment variables to be exported:\n%s", exports)

	// Print the export commands for shell integration
	log.Export(exports)

	log.Debug("Current environment variables:")
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "XDG_") {
			log.Debug(env)
		}
	}
}
