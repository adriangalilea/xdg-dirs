package xdgdirs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"runtime"

	"github.com/adriangalilea/xdg-dirs/internal/logger"
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

// The defaults ARE the XDG spec literals, deliberately NOT platform-adapted.
// This tool exists to give every platform the classic paths (lowercase, no
// spaces, predictable - see README); the adrg/xdg library's macOS opinion
// ("Application Support" for everything) was exactly what it replaces, and
// shipping that opinion as the fallback made a fresh machine diverge from
// the fleet the moment no user.dirs existed. Native mappings are an OPT-IN
// via user.dirs, never a default. XDG_RUNTIME_DIR is deliberately absent:
// the spec says the SYSTEM provides it (lifetime + permission semantics no
// user tool can fake); set it in user.dirs if you must.
func getDefaultXDGDirs() map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get user home directory: %v", err))
	}
	videos := filepath.Join(home, "Videos")
	if runtime.GOOS == "darwin" {
		videos = filepath.Join(home, "Movies")
	}
	return map[string]string{
		"XDG_CACHE_HOME":      filepath.Join(home, ".cache"),
		"XDG_CONFIG_HOME":     filepath.Join(home, ".config"),
		"XDG_DATA_HOME":       filepath.Join(home, ".local", "share"),
		"XDG_STATE_HOME":      filepath.Join(home, ".local", "state"),
		"XDG_DESKTOP_DIR":     filepath.Join(home, "Desktop"),
		"XDG_DOWNLOAD_DIR":    filepath.Join(home, "Downloads"),
		"XDG_DOCUMENTS_DIR":   filepath.Join(home, "Documents"),
		"XDG_MUSIC_DIR":       filepath.Join(home, "Music"),
		"XDG_PICTURES_DIR":    filepath.Join(home, "Pictures"),
		"XDG_VIDEOS_DIR":      videos,
		"XDG_TEMPLATES_DIR":   filepath.Join(home, "Templates"),
		"XDG_PUBLICSHARE_DIR": filepath.Join(home, "Public"),
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
	userDirsPath := filepath.Join(configHome, "xdg", "user.dirs")

	if _, err := os.Stat(userDirsPath); err == nil {
		content, err := os.ReadFile(userDirsPath)
		if err != nil {
			x.logger.Error("Failed to read user.dirs file: %v", err)
			return nil, fmt.Errorf("failed to read user.dirs file: %w", err)
		}

		x.logger.Debug("Contents of %s:\n%s", userDirsPath, string(content))

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "XDG_") && strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					// Strip inline comments
					if idx := strings.Index(value, "#"); idx != -1 {
						value = strings.TrimSpace(value[:idx])
					}
					value = strings.Trim(value, "\"")
					value = os.ExpandEnv(value) // Expand environment variables
					userDirs[key] = value
				}
			}
		}
	} else {
		x.logger.Debug("user.dirs file not found at %s", userDirsPath)
	}

	// Merge user-defined directories with defaults, preferring user-defined values
	defaults := getDefaultXDGDirs()
	for key, defaultValue := range defaults {
		if value, exists := userDirs[key]; !exists || value == "" {
			userDirs[key] = defaultValue
		}
	}

	// Log all merged user directories
	var logEntries []string
	for key, value := range userDirs {
		logEntries = append(logEntries, fmt.Sprintf("%s=%s", key, value))
	}
	x.logger.Debug("Merged user directories:\n%s", strings.Join(logEntries, "\n"))
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

	_, err = file.WriteString("# This file is written by xdg-dirs. Do not edit: it is regenerated on\n# every run. To override a directory, edit user.dirs in the same folder.\n# Entries are sorted by name so identical state diffs byte-identically.\n#\n")
	if err != nil {
		x.logger.Error("Failed to write to generated.dirs file: %v", err)
		return fmt.Errorf("failed to write to generated.dirs file: %w", err)
	}

	keys := make([]string, 0, len(userDirs))
	for key := range userDirs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		_, err = file.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, userDirs[key]))
		if err != nil {
			x.logger.Error("Failed to write to generated.dirs file: %v", err)
			return fmt.Errorf("failed to write to generated.dirs file: %w", err)
		}
	}

	x.logger.Debug("Generated generated.dirs")
	return nil
}
