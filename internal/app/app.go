// Package app wires the Fyne application together and owns its lifecycle.
package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
	"wallplay/internal/ui"
	"wallplay/internal/wallhaven"
)

// AppID is the reverse-DNS identifier used by Fyne for storage and metadata.
const AppID = "dev.wallplay.Wallplay"

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
	log.Println(c.Search("metal gear"))
	w := ui.NewMainWindow(a.fyne)
	w.ShowAndRun()
	return 0
}
