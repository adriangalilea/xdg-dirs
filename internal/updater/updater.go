package updater

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"xdg-user-dirs-cross/internal/logger"
	"xdg-user-dirs-cross/internal/xdgdirs"
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

func (u *Updater) Update(createDirs, dryRun bool) error {
	u.logger.Debug("Starting Update")

	// Read user directories
	userDirs, err := u.xdgDirs.ReadUserDirs()
	if err != nil {
		u.logger.Error("Failed to read user directories: %v", err)
		return fmt.Errorf("failed to read user directories: %w", err)
	}
	u.logger.Debug("Contents of userDirs: %v", userDirs)

	// Handle dry-run mode
	if dryRun {
		u.logger.Debug("Dry run mode: No changes will be applied")
	}

	// Ensure directories are created if needed
	if err := u.ensureDirectories(userDirs, createDirs, dryRun); err != nil {
		u.logger.Error("Failed to ensure directories: %v", err)
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Write updated directories to ~/.config/xdg/generated.dirs if not in dry-run mode
	if !dryRun {
		u.logger.Debug("Not in dry run mode, writing to ~/.config/xdg/generated.dirs")
		if err := u.xdgDirs.WriteUserDirs(u.xdgDirs.Dirs); err != nil {
			u.logger.Error("Failed to write to ~/.config/xdg/generated.dirs: %v", err)
			return fmt.Errorf("failed to write to ~/.config/xdg/generated.dirs: %w", err)
		}
		u.logger.Debug("Updated ~/.config/xdg/generated.dirs")
	}

	u.logger.Debug("XDG user directories update completed.")
	return nil
}

func (u *Updater) ensureDirectories(userDirs map[string]string, createDirs, dryRun bool) error {
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
			if dryRun {
				u.logger.Debug("Would create directory for %s: %s", key, dir)
			} else {
				err := os.MkdirAll(dir, 0700)
				if err != nil {
					u.logger.Error("Failed to create directory for %s: %v", key, err)
					return fmt.Errorf("failed to create directory for %s: %w", key, err)
				}
				u.logger.Debug("Created directory for %s: %s", key, dir)
			}
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
