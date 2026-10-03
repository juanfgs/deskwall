// Package ui contains the Fyne windows and widgets used by the application.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"wallplay/internal/wallhaven"
)

// MainWindow is the application's primary window.
type MainWindow struct {
	fyne.Window
}

// NewMainWindow creates and configures the main application window with a
// search bar and a wallpaper grid backed by the given wallhaven client.
func NewMainWindow(app fyne.App, client wallhaven.Client) *MainWindow {
	w := &MainWindow{Window: app.NewWindow("Wallplay")}
	w.SetMaster()
	w.Resize(fyne.NewSize(900, 650))

	grid := NewWallpaperGrid(client)
	paginator := NewPaginator(grid)

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
	paddedGrid := container.NewPadded(grid)
	paddedPaginator := container.NewPadded(paginator)

	w.SetContent(container.NewBorder(
		container.NewPadded(bar), paddedPaginator, nil, nil, paddedGrid,
	))

	return w
}
