// Package app wires the Fyne application together and owns its lifecycle.
package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"deskwall/internal/ui"
	"deskwall/internal/wallhaven"
)

// AppID is the reverse-DNS identifier used by Fyne for storage and metadata.
const AppID = "dev.deskwall.Deskwall"

// App wraps a fyne.App and manages its windows and lifecycle.
type App struct {
	fyne fyne.App
}

// New constructs a new application instance with the given ID.
func New() *App {
	return &App{fyne: app.NewWithID(AppID)}
}

// Run creates the main window and blocks until the application exits.
// It returns the process exit code (0 on normal exit).
func (a *App) Run() int {
	c := wallhaven.NewClient()
	w := ui.NewMainWindow(a.fyne, c)
	w.ShowAndRun()
	return 0
}
