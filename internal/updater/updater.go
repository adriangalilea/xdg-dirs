package updater

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"xdg-dirs/internal/logger"
	"xdg-dirs/internal/xdgdirs"
)

type Updater struct {
	logger  *logger.Logger
	xdgDirs *xdgdirs.XDGDirs
}

func NewUpdater(log *logger.Logger) *Updater {
	return &Updater{
		logger:  log,
		xdgDirs: xdgdirs.NewXDGDirs(log),
	}
}

func (u *Updater) Update(userDirs map[string]string, createDirs, dryRun bool) error {
	if dryRun {
		u.logger.Debug("Dry run mode: No changes will be applied")
		return nil
	}

	if createDirs {
		if err := u.ensureDirectories(userDirs, createDirs); err != nil {
			u.logger.Error("Failed to ensure directories: %v", err)
			return fmt.Errorf("failed to ensure directories: %w", err)
		}
	}

	generatedDirsPath := filepath.Join(os.Getenv("HOME"), ".config", "xdg", "generated.dirs")
	if err := u.xdgDirs.WriteUserDirs(userDirs); err != nil {
		u.logger.Error("Failed to write to %s: %v", generatedDirsPath, err)
		return fmt.Errorf("failed to write to %s: %w", generatedDirsPath, err)
	}

	u.logger.Debug("Wrote merged XDG directories to %s", generatedDirsPath)
	return nil
}

func (u *Updater) ensureDirectories(userDirs map[string]string, createDirs bool) error {
	if !createDirs {
		return nil
	}
	for key, dir := range userDirs {
		if dir == "" {
			continue
		}
		dir = filepath.Clean(os.ExpandEnv(dir)) // Expand environment variables like $HOME and clean the path

		// Check if the path is valid
		if !filepath.IsAbs(dir) {
			u.logger.Error("Invalid directory path for %s: %s", key, dir)
			return fmt.Errorf("invalid directory path for %s: %s", key, dir)
		}

		// Check if the path is a directory
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				u.logger.Error("Failed to create directory for %s: %v", key, err)
				return fmt.Errorf("failed to create directory for %s: %w", key, err)
			}
			u.logger.Debug("Created directory for %s: %s", key, dir)
		} else if err != nil {
			u.logger.Error("Failed to check directory for %s: %v", key, err)
			return fmt.Errorf("failed to check directory for %s: %w", key, err)
		} else if !info.IsDir() {
			u.logger.Error("Path exists but is not a directory for %s: %s", key, dir)
			return fmt.Errorf("path exists but is not a directory for %s: %s", key, dir)
		}
	}
	return nil
}

func (u *Updater) GetUserDirs() (map[string]string, error) {
	return u.xdgDirs.ReadUserDirs()
}

func (u *Updater) ExportEnv(userDirs map[string]string) string {
	var exports []string
	seen := make(map[string]bool)

	// Include user directories and base directories without duplication
	for key, value := range userDirs {
		if strings.HasPrefix(key, "XDG_") {
			exports = append(exports, fmt.Sprintf("export %s=\"%s\"", key, value))
			seen[key] = true
		}
	}

	for key, value := range u.xdgDirs.Dirs {
		if strings.HasPrefix(key, "XDG_") && !seen[key] {
			exports = append(exports, fmt.Sprintf("export %s=\"%s\"", key, value))
			seen[key] = true
		}
	}

	return strings.Join(exports, "\n")
}
