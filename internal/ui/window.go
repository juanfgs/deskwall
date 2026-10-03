// Package ui contains the Fyne windows and widgets used by the application.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MainWindow is the application's primary window.
type MainWindow struct {
	fyne.Window
}

// NewMainWindow creates and configures the main application window.
func NewMainWindow(app fyne.App) *MainWindow {
	w := &MainWindow{Window: app.NewWindow("Wallplay")}
	w.SetMaster()
	w.Resize(fyne.NewSize(800, 600))

	label := widget.NewLabel("Welcome to Wallplay")
	w.SetContent(container.NewCenter(label))

	return w
}
