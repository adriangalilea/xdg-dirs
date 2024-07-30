package conf

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	DefaultLogFilePath string
	HelpMessage        string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}

	DefaultLogFilePath = filepath.Join(homeDir, ".local", "state", "xdg-user-dirs-cross", "xdg-user-dirs-cross.log")

	HelpMessage = fmt.Sprintf(`xdg-user-dirs-cross: A cross-platform tool for managing XDG user directories

Usage:
xdg-user-dirs-cross [options]

Options:
  -d, --debug        Enable debug output
  -n, --dry-run      Simulate changes without applying them
  -c, --create-dirs  Create directories if they don't exist
  -l, --log-file     Specify the log file path (default: %s)
  -h, --help         Show help message

Configuration:
  This tool looks for ~/.config/xdg/user.dirs

Info:
  This tool moves the ~/.config/user-dirs.dirs (deprecated from xdg-user-dirs and xdg-user-dirs-update) into ~/.config/xdg/user-dirs.dirs-backup
  This tool generates the ~/.config/xdg/user.dirs file.

For more detailed information, please refer to the README.md file.`, DefaultLogFilePath)
}
