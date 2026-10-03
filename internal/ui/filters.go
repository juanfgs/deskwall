package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"wallplay/internal/wallhaven"
)

// categoryNames and purityNames map positions to their API labels.
var (
	categoryNames = []string{"General", "Anime", "People"}
	purityNames   = []string{"SFW", "Sketchy", "NSFW"}
	sortingOpts   = []string{"relevance", "date_added", "random", "views", "favorites", "toplist"}
	orderOpts     = []string{"desc", "asc"}
)

// FilterBar holds the filter controls and produces a wallhaven.SearchRequest
// reflecting the current selection.
type FilterBar struct {
	widget.BaseWidget

	query    *widget.Entry
	catBtns  []*widget.Button
	purBtns  []*widget.Button
	sorting  *widget.Select
	order    *widget.Select

	onChange func()
	ready    bool
}

// NewFilterBar creates a filter bar. onChange is fired whenever any filter
// changes.
func NewFilterBar(onChange func()) *FilterBar {
	f := &FilterBar{onChange: onChange}

	f.query = widget.NewEntry()
	f.query.SetPlaceHolder("Search wallpapers…")

	for _, name := range categoryNames {
		btn := widget.NewButton(name, f.fire)
		btn.Importance = widget.HighImportance
		f.catBtns = append(f.catBtns, btn)
	}

	for _, name := range purityNames {
		btn := widget.NewButton(name, f.fire)
		btn.Importance = widget.HighImportance
		f.purBtns = append(f.purBtns, btn)
	}

	f.sorting = widget.NewSelect(sortingOpts, func(string) { f.fire() })
	f.sorting.SetSelected("relevance")

	f.order = widget.NewSelect(orderOpts, func(string) { f.fire() })
	f.order.SetSelected("desc")

	f.ExtendBaseWidget(f)
	f.ready = true
	return f
}

// fire toggles a category button and notifies the listener.
// It is suppressed until the filter bar has finished constructing.
func (f *FilterBar) fire() {
	if !f.ready || f.onChange == nil {
		return
	}
	f.onChange()
}

// Query returns the search entry text.
func (f *FilterBar) Query() string {
	return f.query.Text
}

// OnSubmitted sets the callback for when the user presses Enter in the entry.
func (f *FilterBar) OnSubmitted(fn func(string)) {
	f.query.OnSubmitted = fn
}

// Request builds a wallhaven.SearchRequest from the current filter state.
func (f *FilterBar) Request() wallhaven.SearchRequest {
	req := wallhaven.NewSearchRequest(f.query.Text)

	cats := make([]string, len(f.catBtns))
	for i, b := range f.catBtns {
		if b.Importance == widget.HighImportance {
			cats[i] = "1"
		}
	}
	req.Categories = cats

	purity := ""
	for _, b := range f.purBtns {
		if b.Importance == widget.HighImportance {
			purity += "1"
		} else {
			purity += "0"
		}
	}
	if purity != "000" {
		req.Purity = purity
	}

	if f.sorting.Selected != "" {
		req.Sorting = f.sorting.Selected
	}
	if f.order.Selected != "" {
		req.Order = f.order.Selected
	}

	return req
}

// toggleButton flips a button between active (High) and inactive (Low).
func toggleButton(b *widget.Button) {
	if b.Importance == widget.HighImportance {
		b.Importance = widget.LowImportance
	} else {
		b.Importance = widget.HighImportance
	}
	b.Refresh()
}

// CreateRenderer lays out the filter bar as two rows:
//   Row 1: [query entry]
//   Row 2: [General Anime People | SFW Sketchy NSFW | Sorting | Order]
func (f *FilterBar) CreateRenderer() fyne.WidgetRenderer {
	for _, b := range f.catBtns {
		btn := b
		orig := btn.OnTapped
		btn.OnTapped = func() {
			toggleButton(btn)
			orig()
		}
	}
	for _, b := range f.purBtns {
		btn := b
		orig := btn.OnTapped
		btn.OnTapped = func() {
			toggleButton(btn)
			orig()
		}
	}

	catRow := container.NewHBox()
	for _, b := range f.catBtns {
		catRow.Add(b)
	}

	purRow := container.NewHBox()
	for _, b := range f.purBtns {
		purRow.Add(b)
	}

	buttonsRow := container.NewHBox(
		catRow,
		widget.NewSeparator(),
		purRow,
		widget.NewSeparator(),
		f.sorting,
		f.order,
	)

	content := container.NewVBox(
		container.NewPadded(f.query),
		container.NewPadded(buttonsRow),
	)
	return widget.NewSimpleRenderer(content)
}
