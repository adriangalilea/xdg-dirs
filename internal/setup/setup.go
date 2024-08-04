package setup

import (
	"os"
	"path/filepath"
	"strings"

	"xdg-dirs/internal/logger"
)

// Prepare performs initial setup tasks like unsetting environment variables
// and backing up the user-dirs.dirs file.
func Prepare(log *logger.Logger) error {
	if err := unsetXDGEnvVars(log); err != nil {
		return err
	}
	return backupUserDirsFile(log)
}

func unsetXDGEnvVars(log *logger.Logger) error {
	xdgEnvVars := []string{
		"XDG_CACHE_HOME", "XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME",
		"XDG_RUNTIME_DIR", "XDG_DESKTOP_DIR", "XDG_DOWNLOAD_DIR", "XDG_DOCUMENTS_DIR",
		"XDG_MUSIC_DIR", "XDG_PICTURES_DIR", "XDG_VIDEOS_DIR", "XDG_TEMPLATES_DIR",
		"XDG_PUBLICSHARE_DIR",
	}
	var unsetVars []string
	for _, envVar := range xdgEnvVars {
		if os.Getenv(envVar) != "" {
			unsetVars = append(unsetVars, envVar)
			os.Unsetenv(envVar)
		}
	}
	log.Debug("Unsetting XDG environment variables:\n%s", strings.Join(unsetVars, "\n"))
	return nil
}

func backupUserDirsFile(log *logger.Logger) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Debug("Failed to get user home directory: %v", err)
		return err
	}
	userDirsFile := filepath.Join(homeDir, ".config", "user-dirs.dirs")
	
	if _, err := os.Stat(userDirsFile); os.IsNotExist(err) {
		log.Debug("user-dirs.dirs didn't exist.")
		return nil
	} else if err != nil {
		log.Debug("Error checking user-dirs.dirs file: %v", err)
		return err
	}

	// If the file exists, proceed with backup
	backupFile := filepath.Join(homeDir, ".config", "xdg", "user-dirs.dirs-backup")
	if err := os.MkdirAll(filepath.Dir(backupFile), 0755); err != nil {
		log.Debug("Failed to create backup directory: %v", err)
		return err
	}
	if err := os.Rename(userDirsFile, backupFile); err != nil {
		log.Debug("Failed to rename user-dirs.dirs to backup file: %v", err)
		return err
	}
	log.Debug("user-dirs.dirs was backed up on %s and deleted.", backupFile)
	return nil
}
