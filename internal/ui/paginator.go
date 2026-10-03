package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Paginator is a pagination control bound to a WallpaperGrid.
// It shows the current page / last page and provides Prev/Next buttons
// that call into the grid. The control refreshes itself via the grid's
// OnPageChange callback whenever the page state changes.
type Paginator struct {
	widget.BaseWidget

	grid   *WallpaperGrid
	prev   *widget.Button
	next   *widget.Button
	label  *widget.Label
}

// NewPaginator creates a Paginator bound to the given grid and registers
// itself for page-change notifications.
func NewPaginator(grid *WallpaperGrid) *Paginator {
	p := &Paginator{grid: grid}
	p.prev = widget.NewButton("Prev", p.grid.PrevPage)
	p.next = widget.NewButton("Next", p.grid.NextPage)
	p.label = widget.NewLabel("")
	p.label.Alignment = fyne.TextAlignCenter
	p.ExtendBaseWidget(p)

	grid.SetOnPageChange(p.refresh)

	return p
}

// CreateRenderer lays out the paginator as [Prev | page x of y | Next].
func (p *Paginator) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewBorder(nil, nil, p.prev, p.next, p.label)
	return widget.NewSimpleRenderer(content)
}

// refresh updates the label text and button enabled states. Called from
// the main thread via the grid's OnPageChange callback.
func (p *Paginator) refresh() {
	current, last, loading := p.grid.PageInfo()
	if last == 0 {
		p.label.SetText("")
	} else {
		p.label.SetText(fmt.Sprintf("Page %d of %d", current, last))
	}
	p.prev.Disable()
	p.next.Disable()
	if loading {
		return
	}
	if current > 1 {
		p.prev.Enable()
	}
	if last > 0 && current < last {
		p.next.Enable()
	}
}
