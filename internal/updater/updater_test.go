package updater

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/adriangalilea/xdg-dirs/internal/logger"
)

// The determinism contract: identical state produces byte-identical output,
// sorted by variable name. This is what lets callers diff runs exactly.
func TestExportEnvDeterministicAndSorted(t *testing.T) {
	log := logger.NewLogger(false, filepath.Join(t.TempDir(), "test.log"))
	u := NewUpdater(log)

	userDirs := map[string]string{
		"XDG_DESKTOP_DIR":  "/home/x/Desktop",
		"XDG_DOWNLOAD_DIR": "/home/x/Downloads",
		"XDG_MUSIC_DIR":    "/home/x/Music",
		"NOT_XDG":          "must-not-appear",
	}

	first := u.ExportEnv(userDirs)
	for i := 0; i < 50; i++ {
		if got := u.ExportEnv(userDirs); got != first {
			t.Fatalf("run %d differs from first run:\n%s\n----\n%s", i, got, first)
		}
	}

	lines := strings.Split(first, "\n")
	if !sort.StringsAreSorted(lines) {
		t.Fatalf("export lines are not sorted:\n%s", first)
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "export XDG_") {
			t.Fatalf("unexpected line (non-XDG leaked or bad format): %q", line)
		}
	}
}

// User-defined values must win over defaults in the merged export.
func TestExportEnvUserOverridesDefaults(t *testing.T) {
	log := logger.NewLogger(false, filepath.Join(t.TempDir(), "test.log"))
	u := NewUpdater(log)

	userDirs := map[string]string{"XDG_DESKTOP_DIR": "/custom/desk"}
	out := u.ExportEnv(userDirs)
	if !strings.Contains(out, "export XDG_DESKTOP_DIR=\"/custom/desk\"") {
		t.Fatalf("user override lost in export:\n%s", out)
	}
}
