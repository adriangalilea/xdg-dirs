package xdgdirs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"xdg-user-dirs-cross/internal/logger"
)

func TestXDGDirsBasicFunctionality(t *testing.T) {
	// Save the original ConfigHome and restore it after the test
	origConfigHome := xdg.ConfigHome
	defer func() { xdg.ConfigHome = origConfigHome }()

	// Create a temporary directory for our test files
	tmpDir, err := ioutil.TempDir("", "xdg-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	xdg.ConfigHome = tmpDir
	os.Setenv("HOME", tmpDir)

	log := logger.NewLogger(true)
	xdgDirs := NewXDGDirs(log)

	// Test ReadUserDirs and WriteUserDirs
	userDirs, err := xdgDirs.ReadUserDirs()
	if err != nil {
		t.Fatalf("ReadUserDirs() error = %v", err)
	}

	err = xdgDirs.WriteUserDirs(userDirs)
	if err != nil {
		t.Fatalf("WriteUserDirs() error = %v", err)
	}

	// Check if user-dirs.dirs was generated
	userDirsDirsPath := filepath.Join(tmpDir, "user-dirs.dirs")
	if _, err := os.Stat(userDirsDirsPath); os.IsNotExist(err) {
		t.Errorf("user-dirs.dirs file was not created")
	}

	// Read and parse the generated user-dirs.dirs file
	content, err := ioutil.ReadFile(userDirsDirsPath)
	if err != nil {
		t.Fatalf("Failed to read user-dirs.dirs: %v", err)
	}

	// Parse the content and check if directories were created
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "XDG_") && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				dirPath := strings.Trim(parts[1], "\"")
				dirPath = strings.Replace(dirPath, "$HOME", tmpDir, 1)
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					t.Errorf("Directory %s was not created", dirPath)
				}
			}
		}
	}
}

func TestXDGDirsWithCustomUserDirs(t *testing.T) {
	// Save the original ConfigHome and restore it after the test
	origConfigHome := xdg.ConfigHome
	defer func() { xdg.ConfigHome = origConfigHome }()

	// Create a temporary directory for our test files
	tmpDir, err := ioutil.TempDir("", "xdg-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	xdg.ConfigHome = tmpDir
	os.Setenv("HOME", tmpDir)

	// Create a custom user.dirs file
	customUserDirs := `XDG_DESKTOP_DIR="$HOME/CustomDesktop"
XDG_DOWNLOAD_DIR="$HOME/CustomDownloads"
`
	err = ioutil.WriteFile(filepath.Join(tmpDir, "user.dirs"), []byte(customUserDirs), 0644)
	if err != nil {
		t.Fatalf("Failed to create custom user.dirs: %v", err)
	}

	log := logger.NewLogger(true)
	xdgDirs := NewXDGDirs(log)

	// Test ReadUserDirs and WriteUserDirs
	userDirs, err := xdgDirs.ReadUserDirs()
	if err != nil {
		t.Fatalf("ReadUserDirs() error = %v", err)
	}

	err = xdgDirs.WriteUserDirs(userDirs)
	if err != nil {
		t.Fatalf("WriteUserDirs() error = %v", err)
	}

	// Check if user-dirs.dirs was generated
	userDirsDirsPath := filepath.Join(tmpDir, "user-dirs.dirs")
	content, err := ioutil.ReadFile(userDirsDirsPath)
	if err != nil {
		t.Fatalf("Failed to read user-dirs.dirs: %v", err)
	}

	// Check if custom directories are present in the generated file
	if !strings.Contains(string(content), `XDG_DESKTOP_DIR="$HOME/CustomDesktop"`) {
		t.Errorf("Custom Desktop directory not found in user-dirs.dirs")
	}
	if !strings.Contains(string(content), `XDG_DOWNLOAD_DIR="$HOME/CustomDownloads"`) {
		t.Errorf("Custom Downloads directory not found in user-dirs.dirs")
	}

	// Check if custom directories were created
	customDesktop := filepath.Join(tmpDir, "CustomDesktop")
	customDownloads := filepath.Join(tmpDir, "CustomDownloads")
	if _, err := os.Stat(customDesktop); os.IsNotExist(err) {
		t.Errorf("Custom Desktop directory was not created")
	}
	if _, err := os.Stat(customDownloads); os.IsNotExist(err) {
		t.Errorf("Custom Downloads directory was not created")
	}
}
