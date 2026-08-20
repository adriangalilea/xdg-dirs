package xdgdirs

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/adriangalilea/xdg-dirs/internal/logger"
)

// Test 1: Core feature - user config actually overrides defaults
func TestUserConfigOverridesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Unsetenv("HOME")

	// User wants cache in .local/cache instead of default
	xdgDir := filepath.Join(tmpDir, "xdg")
	os.MkdirAll(xdgDir, 0755)
	os.WriteFile(filepath.Join(xdgDir, "user.dirs"), []byte(`XDG_CACHE_HOME="$HOME/.local/cache"`), 0644)

	log := logger.NewLogger(false, "")
	x := NewXDGDirs(log)
	dirs, _ := x.ReadUserDirs()

	expected := filepath.Join(tmpDir, ".local/cache")
	if dirs["XDG_CACHE_HOME"] != expected {
		t.Errorf("User config not respected: got %s, want %s", dirs["XDG_CACHE_HOME"], expected)
	}
}

// Test 2: Critical - environment variables expand correctly
func TestEnvironmentVariableExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	os.Setenv("HOME", "/home/testuser")
	defer os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Unsetenv("HOME")

	xdgDir := filepath.Join(tmpDir, "xdg")
	os.MkdirAll(xdgDir, 0755)
	os.WriteFile(filepath.Join(xdgDir, "user.dirs"), []byte(`XDG_DESKTOP_DIR="$HOME/Desktop"`), 0644)

	log := logger.NewLogger(false, "")
	x := NewXDGDirs(log)
	dirs, _ := x.ReadUserDirs()

	if dirs["XDG_DESKTOP_DIR"] != "/home/testuser/Desktop" {
		t.Errorf("$HOME not expanded: got %s", dirs["XDG_DESKTOP_DIR"])
	}
}

// Test 3: Robustness - handles real-world messy configs
func TestHandlesMalformedConfig(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Unsetenv("HOME")

	xdgDir := filepath.Join(tmpDir, "xdg")
	os.MkdirAll(xdgDir, 0755)

	messyConfig := `# User's comments
XDG_DESKTOP_DIR="$HOME/Desktop"  # inline comment
MALFORMED LINE WITHOUT EQUALS
XDG_DOWNLOAD_DIR="$HOME/Downloads"
`
	os.WriteFile(filepath.Join(xdgDir, "user.dirs"), []byte(messyConfig), 0644)

	log := logger.NewLogger(false, "")
	x := NewXDGDirs(log)

	// Should not crash
	dirs, err := x.ReadUserDirs()
	if err != nil {
		t.Fatalf("Should handle malformed config gracefully: %v", err)
	}

	// Should parse valid lines
	if !strings.HasSuffix(dirs["XDG_DOWNLOAD_DIR"], "Downloads") {
		t.Error("Failed to parse valid lines from messy config")
	}

	// Bug: Currently fails due to inline comment not being stripped
	if strings.Contains(dirs["XDG_DESKTOP_DIR"], "#") {
		t.Error("Inline comments not stripped from values")
	}
}

// Test 4: Platform-specific behavior
func TestPlatformSpecificDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Unsetenv("HOME")

	log := logger.NewLogger(false, "")
	x := NewXDGDirs(log)
	dirs, _ := x.ReadUserDirs()

	// Just verify platform differences exist
	if runtime.GOOS == "darwin" {
		if !strings.Contains(dirs["XDG_VIDEOS_DIR"], "Movies") {
			t.Error("macOS should use Movies folder")
		}
	}
}
