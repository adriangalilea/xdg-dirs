package xdgdirs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"io/ioutil"
	"strings"
	"sync"
	"xdg-dirs/internal/logger"
)

type XDGDirs struct {
	logger *logger.Logger
	mu     sync.Mutex
	Dirs   map[string]string
}

func init() {
}

func NewXDGDirs(log *logger.Logger) *XDGDirs {
	return &XDGDirs{
		logger: log,
		Dirs:   getDefaultXDGDirs(),
	}
}

func getDefaultXDGDirs() map[string]string {
	return map[string]string{
		"XDG_CACHE_HOME":     xdg.CacheHome,
		"XDG_CONFIG_HOME":    xdg.ConfigHome,
		"XDG_DATA_HOME":      xdg.DataHome,
		"XDG_STATE_HOME":     xdg.StateHome,
		"XDG_RUNTIME_DIR":    xdg.RuntimeDir,
		"XDG_DESKTOP_DIR":    xdg.UserDirs.Desktop,
		"XDG_DOWNLOAD_DIR":   xdg.UserDirs.Download,
		"XDG_DOCUMENTS_DIR":  xdg.UserDirs.Documents,
		"XDG_MUSIC_DIR":      xdg.UserDirs.Music,
		"XDG_PICTURES_DIR":   xdg.UserDirs.Pictures,
		"XDG_VIDEOS_DIR":     xdg.UserDirs.Videos,
		"XDG_TEMPLATES_DIR":  xdg.UserDirs.Templates,
		"XDG_PUBLICSHARE_DIR": xdg.UserDirs.PublicShare,
	}
}

func (x *XDGDirs) ReadUserDirs() (map[string]string, error) {
	userDirs := getDefaultXDGDirs()
	x.mu.Lock()
	defer x.mu.Unlock()

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			x.logger.Error("Failed to get user home directory: %v", err)
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	userDirsPath := filepath.Join(configHome, "xdg", "generated.dirs")
	x.logger.Debug("Attempting to read generated.dirs from: %s", userDirsPath)

	if _, err := os.Stat(userDirsPath); err == nil {
		content, err := ioutil.ReadFile(userDirsPath)
		if err != nil {
			x.logger.Error("Failed to read generated.dirs file: %v", err)
			return nil, fmt.Errorf("failed to read generated.dirs file: %w", err)
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "XDG_") && strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
					userDirs[key] = value
				}
			}
		}
	}

	x.logger.Debug("Contents of generated.dirs: %v", userDirs)

	// Merge user-defined directories with defaults, preferring user-defined values
	defaults := getDefaultXDGDirs()
	for key, defaultValue := range defaults {
		if value, exists := userDirs[key]; !exists || value == "" {
			userDirs[key] = defaultValue
		}
	}
	x.logger.Debug("Merged user directories: %v", userDirs)
	return userDirs, nil
}

func (x *XDGDirs) WriteUserDirs(userDirs map[string]string) error {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			x.logger.Error("Failed to get user home directory: %v", err)
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		configHome = filepath.Join(homeDir, ".config")
	}
	xdgConfigDir := filepath.Join(configHome, "xdg")
	if err := os.MkdirAll(xdgConfigDir, 0755); err != nil {
		x.logger.Error("Failed to create XDG config directory: %v", err)
		return fmt.Errorf("failed to create XDG config directory: %w", err)
	}
	userDirsFile := filepath.Clean(filepath.Join(configHome, "xdg", "generated.dirs"))

	file, err := os.OpenFile(userDirsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		x.logger.Error("Failed to create generated.dirs file: %v", err)
		return fmt.Errorf("failed to create generated.dirs file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("# This file is written by xdg-user-dirs-update\n# If you want to change or add directories, just edit the line you're\n# interested in. All local changes will be retained on the next run.\n# Format is XDG_xxx_DIR=\"$HOME/yyy\", where yyy is a shell-escaped\n# homedir-relative path, or XDG_xxx_DIR=\"/yyy\", where /yyy is an\n# absolute path. No other format is supported.\n#\n")
	if err != nil {
		x.logger.Error("Failed to write to generated.dirs file: %v", err)
		return fmt.Errorf("failed to write to generated.dirs file: %w", err)
	}

	for key, value := range userDirs {
		_, err = file.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, value))
		if err != nil {
			x.logger.Error("Failed to write to generated.dirs file: %v", err)
			return fmt.Errorf("failed to write to generated.dirs file: %w", err)
		}
	}

	x.logger.Debug("Generated generated.dirs")
	return nil
}
