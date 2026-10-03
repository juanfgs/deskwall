// Package ui contains the Fyne windows and widgets used by the application.
package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"wallplay/internal/wallhaven"
)

// elasticCurve is an ease-out curve with a slight overshoot for a
// spring-like "grow" feel.
func elasticCurve(v float32) float32 {
	t := v - 1
	return t*t*((1.5*t)+1.5) + 1
}

// MainWindow is the application's primary window.
type MainWindow struct {
	fyne.Window
}

// NewMainWindow creates and configures the main application window with a
// search bar, a wallpaper grid, and an expandable image detail view.
func NewMainWindow(app fyne.App, client wallhaven.Client) *MainWindow {
	w := &MainWindow{Window: app.NewWindow("Wallplay")}
	w.SetMaster()
	w.Resize(fyne.NewSize(1200, 700))

	grid := NewWallpaperGrid(client)
	paginator := NewPaginator(grid)

	gridArea := container.NewBorder(nil, container.NewPadded(paginator), nil, nil, container.NewPadded(grid))

	expanded := NewExpandedView(client, nil)
	expanded.Hide()

	// Overlay the expanded view on top of the grid.
	stack := container.NewStack(gridArea, expanded)

	// Wire up the close behavior: hide expanded view, show grid.
	expanded.onClose = func() {
		expanded.Hide()
	}

	// When a wallpaper is selected, show the expanded view with a grow animation.
	grid.OnSelect = func(wp wallhaven.Wallpaper) {
		expanded.Show()
		expanded.Load(wp)
		expanded.GrowImage(400 * time.Millisecond)
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Search wallpapers…")

	search := func() {
		if q := entry.Text; q != "" {
			grid.Search(q)
		}
	}
	entry.OnSubmitted = func(_ string) { search() }
	button := widget.NewButton("Search", search)

	bar := container.NewBorder(nil, nil, nil, button, entry)

	w.SetContent(container.NewBorder(container.NewPadded(bar), nil, nil, nil, stack))

	return w
}
