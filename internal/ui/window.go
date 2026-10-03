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
// filter bar, a wallpaper grid, and an expandable image detail view.
func NewMainWindow(app fyne.App, client wallhaven.Client) *MainWindow {
	w := &MainWindow{Window: app.NewWindow("Wallplay")}
	w.SetMaster()
	w.Resize(fyne.NewSize(1200, 700))

	grid := NewWallpaperGrid(client)
	paginator := NewPaginator(grid)

	var filters *FilterBar
	runSearch := func() {
		req := filters.Request()
		if req.Query != "" {
			grid.SearchWithRequest(req)
		}
	}
	filters = NewFilterBar(runSearch)
	filters.OnSubmitted(func(_ string) { runSearch() })

	searchBtn := widget.NewButton("Search", runSearch)

	gridArea := container.NewBorder(nil, container.NewPadded(paginator), nil, nil, container.NewPadded(grid))

	expanded := NewExpandedView(client, nil)
	expanded.Hide()

	stack := container.NewStack(gridArea, expanded)

	expanded.onClose = func() {
		expanded.Hide()
	}

	grid.OnSelect = func(wp wallhaven.Wallpaper) {
		expanded.Show()
		expanded.Load(wp)
		expanded.GrowImage(400 * time.Millisecond)
	}

	filterRow := container.NewBorder(nil, nil, nil, searchBtn, filters)

	w.SetContent(container.NewBorder(container.NewPadded(filterRow), nil, nil, nil, stack))

	return w
}
