package platform

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// PicturesDir returns the platform-appropriate pictures directory:
//   - Linux: parses XDG user-dirs.dirs, falls back to ~/Pictures
//   - macOS: ~/Pictures
//   - Windows: %USERPROFILE%\Pictures
func PicturesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	switch runtime.GOOS {
	case "linux":
		return linuxPicturesDir(home)
	case "darwin":
		return filepath.Join(home, "Pictures")
	case "windows":
		return filepath.Join(home, "Pictures")
	default:
		return filepath.Join(home, "Pictures")
	}
}

// linuxPicturesDir parses ~/.config/user-dirs.dirs for XDG_PICTURES_DIR,
// falling back to ~/Pictures.
func linuxPicturesDir(home string) string {
	cfgPath := filepath.Join(home, ".config", "user-dirs.dirs")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return filepath.Join(home, "Pictures")
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "XDG_PICTURES_DIR=") {
			val := strings.Trim(strings.TrimPrefix(line, "XDG_PICTURES_DIR="), `"`)
			val = strings.ReplaceAll(val, "$HOME", home)
			return val
		}
	}
	return filepath.Join(home, "Pictures")
}
