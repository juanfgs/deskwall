// Package platform provides cross-platform abstractions for OS-specific
// operations like setting the desktop wallpaper and resolving the
// user's Pictures directory.
package platform

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// WallpaperSetter is the strategy interface for applying a wallpaper file
// to the desktop background. Each platform provides its own implementation.
type WallpaperSetter interface {
	// SetWallpaper applies the image at the given absolute path as the
	// desktop wallpaper. Returns an error if the platform is unsupported
	// or the operation fails.
	SetWallpaper(path string) error
}

// NewWallpaperSetter returns the platform-appropriate WallpaperSetter.
func NewWallpaperSetter() WallpaperSetter {
	switch runtime.GOOS {
	case "linux":
		return &gnomeSetter{}
	case "darwin":
		return &macosSetter{}
	case "windows":
		return &windowsSetter{}
	default:
		return &unsupportedSetter{os: runtime.GOOS}
	}
}

// --- Linux / GNOME ---

type gnomeSetter struct{}

func (g *gnomeSetter) SetWallpaper(path string) error {
	uri := "file://" + path
	for _, key := range []string{"picture-uri-dark", "picture-uri"} {
		cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", key, uri)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("gsettings: %w", err)
		}
	}
	return nil
}

// --- macOS ---

type macosSetter struct{}

func (m *macosSetter) SetWallpaper(path string) error {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	script := fmt.Sprintf(`tell application "System Events" to set picture of every desktop to POSIX file "%s"`, abspath)
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("osascript: %w", err)
	}
	return nil
}

// --- Windows ---

type windowsSetter struct{}

func (w *windowsSetter) SetWallpaper(path string) error {
	// Windows uses SPI_SETDESKWALLPAPER (0x0014) via SystemParametersInfoW.
	// We invoke it through a tiny inline Go syscall via the syscall package
	// since cgo on Windows requires a C compiler (mingw). Using powershell
	// with a P/Invoke is cleaner and avoids the cgo dependency.
	abspath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	ps := fmt.Sprintf(`Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class WP {
  [DllImport("user32.dll", CharSet=CharSet.Unicode, SetLastError=true)]
  public static extern bool SystemParametersInfoW(uint uAction, uint uParam, string lpvParam, uint fuWin);
}
"@
[WP]::SystemParametersInfoW(0x0014, 0, "%s", 0x02)`, abspath)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", ps)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("SystemParametersInfoW: %w", err)
	}
	return nil
}

// --- Unsupported fallback ---

type unsupportedSetter struct {
	os string
}

func (u *unsupportedSetter) SetWallpaper(_ string) error {
	return fmt.Errorf("setting wallpaper is not supported on %s", u.os)
}
