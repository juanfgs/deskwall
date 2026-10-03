package ui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"wallplay/internal/wallhaven"
)

// thumbCellSize is the target size of each grid cell. NewGridWrap reflows
// the number of columns to fit the available width.
var thumbCellSize = fyne.NewSize(200, 140)

// glowInset is the margin between the gradient glow and the thumbnail image.
const glowInset = 4

// glowStart/glowEnd define the diagonal gradient rendered around the
// selected thumbnail.
var (
	glowStart = color.NRGBA{R: 0x35, G: 0x84, B: 0xe4, A: 0xff}
	glowEnd   = color.NRGBA{R: 0x91, G: 0x41, B: 0xe8, A: 0xff}
)

// WallpaperGrid displays wallhaven search results as a scrollable grid of
// thumbnails loaded asynchronously.
type WallpaperGrid struct {
	widget.BaseWidget

	client wallhaven.Client
	status *widget.Label
	grid   *fyne.Container
	scroll *container.Scroll

	mu           sync.Mutex
	selectedCell *thumbCell
	req          wallhaven.SearchRequest
	lastPage     int
	currentPage  int
	loading      bool
	onPageChange func()
	OnSelect     func(wallhaven.Wallpaper)
}

// NewWallpaperGrid constructs a grid ready to receive Search calls.
func NewWallpaperGrid(client wallhaven.Client) *WallpaperGrid {
	g := &WallpaperGrid{
		client: client,
		status: widget.NewLabel(""),
		grid:   container.NewGridWrap(thumbCellSize),
	}
	g.scroll = container.NewScroll(g.grid)
	g.ExtendBaseWidget(g)
	return g
}

// Search starts a new search for the given query at page 1.
func (g *WallpaperGrid) Search(query string) {
	g.loadPage(wallhaven.NewSearchRequest(query), true)
}

// SearchWithRequest starts a new search using the given request (with filters)
// at page 1.
func (g *WallpaperGrid) SearchWithRequest(req wallhaven.SearchRequest) {
	g.loadPage(req, true)
}

// NextPage loads the next page if one exists and we're not already loading.
func (g *WallpaperGrid) NextPage() {
	g.mu.Lock()
	if g.loading || g.currentPage >= g.lastPage {
		g.mu.Unlock()
		return
	}
	req, err := g.req.NextPage()
	if err != nil {
		g.mu.Unlock()
		return
	}
	g.mu.Unlock()
	g.loadPage(req, false)
}

// PrevPage loads the previous page if one exists and we're not already loading.
func (g *WallpaperGrid) PrevPage() {
	g.mu.Lock()
	if g.loading || g.currentPage <= 1 {
		g.mu.Unlock()
		return
	}
	req, err := g.req.PrevPage()
	if err != nil {
		g.mu.Unlock()
		return
	}
	g.mu.Unlock()
	g.loadPage(req, false)
}

// PageInfo returns the current page, last page, and whether a request is
// in flight. Safe to call from the UI thread.
func (g *WallpaperGrid) PageInfo() (current, last int, loading bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.currentPage, g.lastPage, g.loading
}

// SetOnPageChange sets a callback invoked whenever the page state changes
// (after a search completes or a page boundary error is hit).
func (g *WallpaperGrid) SetOnPageChange(fn func()) {
	g.mu.Lock()
	g.onPageChange = fn
	g.mu.Unlock()
}

// loadPage executes a search request and populates the grid with results.
// When reset is true the selection and request are reset to a fresh search.
func (g *WallpaperGrid) loadPage(req wallhaven.SearchRequest, reset bool) {
	g.mu.Lock()
	if g.loading {
		g.mu.Unlock()
		return
	}
	g.loading = true
	if reset {
		g.req = req
		g.selectedCell = nil
	}
	g.status.SetText("Searching…")
	g.status.Show()
	g.grid.Objects = nil
	g.grid.Refresh()
	g.scroll.Offset = fyne.NewPos(0, 0)
	g.mu.Unlock()
	g.notifyPageChange()

	go func() {
		resp, err := g.client.SearchWithRequest(req)
		g.mu.Lock()
		g.loading = false
		if err != nil {
			g.mu.Unlock()
			fyne.Do(func() {
				g.status.SetText(fmt.Sprintf("Error: %v", err))
				g.status.Show()
			})
			g.notifyPageChange()
			return
		}
		g.lastPage = resp.Meta.LastPage
		g.currentPage = resp.Meta.CurrentPage
		g.mu.Unlock()

		if len(resp.Data) == 0 {
			fyne.Do(func() {
				g.status.SetText("No results")
				g.status.Show()
			})
			g.notifyPageChange()
			return
		}

		cells := make([]fyne.CanvasObject, len(resp.Data))
		for i, wp := range resp.Data {
			cell := newThumbCell(g, wp)
			cells[i] = cell
			g.loadThumbnail(cell.img, wp)
		}

		fyne.Do(func() {
			g.grid.Objects = cells
			g.grid.Refresh()
			g.status.Hide()
		})
		g.notifyPageChange()
	}()
}

func (g *WallpaperGrid) notifyPageChange() {
	g.mu.Lock()
	fn := g.onPageChange
	g.mu.Unlock()
	if fn != nil {
		fyne.Do(fn)
	}
}

// selectCell deselects the current selection and marks c as selected,
// then fires the OnSelect callback (if set) with the cell's wallpaper.
func (g *WallpaperGrid) selectCell(c *thumbCell) {
	if g.selectedCell == c {
		return
	}
	if g.selectedCell != nil {
		g.selectedCell.setSelected(false)
	}
	g.selectedCell = c
	c.setSelected(true)
	if g.OnSelect != nil {
		g.OnSelect(c.wallpaper)
	}
}

// loadThumbnail fetches the small thumbnail for a wallpaper in a background
// goroutine and updates the image on the main thread when ready.
func (g *WallpaperGrid) loadThumbnail(img *canvas.Image, wp wallhaven.Wallpaper) {
	go func() {
		res, err := http.Get(wp.Thumbs.Small)
		if err != nil {
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}
		if !decodeable(body) {
			return
		}
		fyne.Do(func() {
			img.Resource = fyne.NewStaticResource(wp.ID, body)
			img.Refresh()
		})
	}()
}

// decodeable reports whether body is a decodable image (lets us skip setting
// a resource that would render as a broken placeholder).
func decodeable(body []byte) bool {
	_, _, err := image.DecodeConfig(bytes.NewReader(body))
	return err == nil
}

// CreateRenderer satisfies widget.BaseWidget. The status label overlays the
// top of the scroll area so messages remain visible above the grid.
func (g *WallpaperGrid) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewBorder(g.status, nil, nil, nil, g.scroll)
	return widget.NewSimpleRenderer(content)
}

// thumbCell is a single selectable grid cell containing a thumbnail image and
// a gradient glow shown when selected.
type thumbCell struct {
	widget.BaseWidget

	img       *canvas.Image
	selected  bool
	grid      *WallpaperGrid
	wallpaper wallhaven.Wallpaper
}

func newThumbCell(grid *WallpaperGrid, wp wallhaven.Wallpaper) *thumbCell {
	c := &thumbCell{
		img:       canvas.NewImageFromResource(nil),
		grid:      grid,
		wallpaper: wp,
	}
	c.img.FillMode = canvas.ImageFillCover
	c.ExtendBaseWidget(c)
	return c
}

// Tapped selects this cell.
func (c *thumbCell) Tapped(*fyne.PointEvent) {
	c.grid.selectCell(c)
}

func (c *thumbCell) setSelected(b bool) {
	c.selected = b
	c.Refresh()
}

// CreateRenderer builds a renderer that stacks the gradient glow behind the
// thumbnail image, inset by glowInset to reveal the glow as a border.
func (c *thumbCell) CreateRenderer() fyne.WidgetRenderer {
	glow := canvas.NewLinearGradient(glowStart, glowEnd, 45)
	return &thumbCellRenderer{cell: c, glow: glow}
}

// thumbCellRenderer lays out the glow and image and toggles glow visibility.
type thumbCellRenderer struct {
	cell *thumbCell
	glow *canvas.LinearGradient
}

func (r *thumbCellRenderer) Layout(size fyne.Size) {
	r.glow.Resize(size)
	r.cell.img.Resize(fyne.NewSize(size.Width-2*glowInset, size.Height-2*glowInset))
	r.cell.img.Move(fyne.NewPos(glowInset, glowInset))
}

func (r *thumbCellRenderer) MinSize() fyne.Size {
	return thumbCellSize
}

func (r *thumbCellRenderer) Refresh() {
	if r.cell.selected {
		r.glow.Show()
	} else {
		r.glow.Hide()
	}
	r.glow.Refresh()
	r.cell.img.Refresh()
}

func (r *thumbCellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.glow, r.cell.img}
}

func (r *thumbCellRenderer) Destroy() {}
