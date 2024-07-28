package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"xdg-user-dirs-cross/internal/logger"
	"xdg-user-dirs-cross/internal/updater"
)

var log *logger.Logger

func main() {
	// Parse command-line flags
	debug := flag.Bool("d", false, "Enable debug output")
	dryRun := flag.Bool("n", false, "Simulate changes without applying them")
	createDirs := flag.Bool("c", false, "Create directories if they don't exist")
	defaultLogFilePath := filepath.Join(os.TempDir(), "xdg-user-dirs-update.log")
	logFilePath := flag.String("l", defaultLogFilePath, "Specify the log file path")
	help := flag.Bool("help", false, "Show help message")
	flag.BoolVar(help, "h", false, "Show help message")
	flag.Parse()


	// Display help message if requested
	if *help {
		fmt.Println("xdg-user-dirs-cross: A cross-platform tool for managing XDG user directories")
		fmt.Println("\nUsage:")
		fmt.Println("xdg-user-dirs-cross [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -d, --debug        Enable debug output")
		fmt.Println("  -n, --dry-run      Simulate changes without applying them")
		fmt.Println("  -c, --create-dirs  Create directories if they don't exist")
		fmt.Printf("  -l, --log-file     Specify the log file path (default: %s)\n", *logFilePath)
		fmt.Println("  -h, --help         Show help message")
		fmt.Println("\nConfiguration:")
		fmt.Println("  This tool looks for ~/.config/xdg/user.dirs")
		fmt.Println("\nInfo:")
		fmt.Println("  This tool moves the ~/.config/user-dirs.dirs (deprecated from xdg-user-dirs and xdg-user-dirs-update) into ~/.config/xdg/user-dirs.dirs-backup")
		fmt.Println("  This tool generates the ~/.config/xdg/user.dirs file.")
		fmt.Println("\nFor more detailed information, please refer to the README.md file.")
		os.Exit(0)
	}

	log = logger.NewLogger(*debug, *logFilePath)

	// Pre-checks: Remove ~/.config/user-dirs.dirs if it exists and unset XDG environment variables
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory: %v", err)
	}
	userDirsFile := filepath.Join(homeDir, ".config", "user-dirs.dirs")
	backupFile := filepath.Join(homeDir, ".config", "xdg", "user-dirs.dirs-backup")
	if _, err := os.Stat(userDirsFile); err == nil {
		log.Debug("Backing up existing user-dirs.dirs file: %s to %s", userDirsFile, backupFile)
		if err := os.MkdirAll(filepath.Dir(backupFile), 0755); err != nil {
			log.Fatal("Failed to create directory: %v", err)
		}
		if err := os.Rename(userDirsFile, backupFile); err != nil {
			log.Fatal("Failed to backup existing user-dirs.dirs file: %v", err)
		}
	} else if !os.IsNotExist(err) {
		log.Fatal("Error checking for existing user-dirs.dirs file: %v", err)
	}

	xdgEnvVars := []string{
		"XDG_CACHE_HOME", "XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME",
		"XDG_RUNTIME_DIR", "XDG_DESKTOP_DIR", "XDG_DOWNLOAD_DIR", "XDG_DOCUMENTS_DIR",
		"XDG_MUSIC_DIR", "XDG_PICTURES_DIR", "XDG_VIDEOS_DIR", "XDG_TEMPLATES_DIR",
		"XDG_PUBLICSHARE_DIR",
	}
	log.Debug("Unsetting XDG environment variables")
	for _, envVar := range xdgEnvVars {
		log.Debug("Unsetting environment variable: %s", envVar)
		os.Unsetenv(envVar)
	}
	// Create updater instance
	updaterInstance := updater.NewUpdater(log)

	// Get user directories
	userDirs, err := updaterInstance.GetUserDirs()
	if err != nil {
		log.Fatal("Failed to get user directories: %v", err)
	}

	// Update user directorie
	if err := updaterInstance.Update(*createDirs, *dryRun); err != nil {
		log.Fatal("Failed to update user directories: %v", err)
	}

	// Get updated user directories
	userDirs, err = updaterInstance.GetUserDirs()
	if err != nil {
		log.Fatal("Failed to get user directories: %v", err)
	}

	// get the EXPORT env variables
	exports := updaterInstance.ExportEnv(userDirs)

	log.Debug("XDG environment variables to be exported:\n%s", exports)

	// Print the export commands for shell integration
	// Note: We don't set environment variables directly here because
	// this program is meant to be evaluated by the shell, not to modify
	// its own environment which would have no effect on the parent shell.
	log.Export(exports)

	log.Debug("Current environment variables:")
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "XDG_") {
			log.Debug(env)
		}
	}

}
