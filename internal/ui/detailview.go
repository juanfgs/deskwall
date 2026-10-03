package ui

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"deskwall/internal/wallhaven"
)

// ExpandedView shows a large wallpaper image with an info sidebar.
// It overlays the grid and grows elastically when shown.
type ExpandedView struct {
	widget.BaseWidget

	client  wallhaven.Client
	image   *canvas.Image
	spinner *widget.ProgressBarInfinite
	info    *fyne.Container
	scroll  *container.Scroll
	close    *widget.Button
	download *widget.Button
	apply    *widget.Button
	current  wallhaven.Wallpaper
	onClose  func()
	bg       *canvas.Rectangle
}

// NewExpandedView constructs a hidden expanded view bound to the given client.
func NewExpandedView(client wallhaven.Client, onClose func()) *ExpandedView {
	v := &ExpandedView{
		client:  client,
		image:   canvas.NewImageFromResource(nil),
		spinner: widget.NewProgressBarInfinite(),
		onClose: onClose,
		bg:      canvas.NewRectangle(color.Transparent),
	}
	v.image.FillMode = canvas.ImageFillContain
	v.spinner.Hide()
	v.bg.FillColor = theme.Color(theme.ColorNameBackground)

	v.close = widget.NewButton("Close", func() {
		if v.onClose != nil {
			v.onClose()
		}
	})

	v.download = widget.NewButton("Download", nil)
	v.download.Importance = widget.HighImportance
	v.download.OnTapped = func() {
		go v.downloadWallpaper(v.current)
	}

	v.apply = widget.NewButton("Apply", nil)
	v.apply.Importance = widget.SuccessImportance
	v.apply.OnTapped = func() {
		go v.applyWallpaper(v.current)
	}

	v.info = container.NewVBox()

	v.scroll = container.NewScroll(v.image)

	v.ExtendBaseWidget(v)
	v.Hide()
	return v
}

// Load fetches the full image and detail for the given wallpaper.
func (v *ExpandedView) Load(wp wallhaven.Wallpaper) {
	v.current = wp
	v.spinner.Show()
	v.image.Resource = nil
	v.image.Refresh()
	v.info.Objects = nil
	v.info.Refresh()

	go v.loadImage(wp.ID, wp.Path)
	go v.loadDetail(wp.ID)
}

// downloadWallpaper downloads the full-resolution image to the XDG
// Pictures directory (the default backgrounds location for GNOME and KDE).
func (v *ExpandedView) downloadWallpaper(wp wallhaven.Wallpaper) {
	fyne.Do(func() {
		v.download.SetText("Downloading…")
		v.download.Disable()
	})

	destDir := picturesDir()
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		fyne.Do(func() {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			v.download.SetText("Download")
			v.download.Enable()
		})
		return
	}

	ext := filepath.Ext(wp.Path)
	if ext == "" {
		ext = ".jpg"
	}
	filename := "wallhaven-" + wp.ID + ext
	destPath := filepath.Join(destDir, filename)

	res, err := http.Get(wp.Path)
	if err != nil || res.StatusCode != http.StatusOK {
		if res != nil {
			res.Body.Close()
		}
		fyne.Do(func() {
			dialog.ShowError(fmt.Errorf("download failed: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
			v.download.SetText("Download")
			v.download.Enable()
		})
		return
	}

	out, err := os.Create(destPath)
	if err != nil {
		res.Body.Close()
		fyne.Do(func() {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			v.download.SetText("Download")
			v.download.Enable()
		})
		return
	}

	_, err = io.Copy(out, res.Body)
	res.Body.Close()
	out.Close()

	fyne.Do(func() {
		if err != nil {
			dialog.ShowError(fmt.Errorf("save failed: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
		} else {
			dialog.ShowInformation("Download complete", "Saved to:\n"+destPath, fyne.CurrentApp().Driver().AllWindows()[0])
		}
		v.download.SetText("Download")
		v.download.Enable()
	})
}

// applyWallpaper downloads the full image to the XDG Pictures directory
// and applies it as the GNOME wallpaper via gsettings.
func (v *ExpandedView) applyWallpaper(wp wallhaven.Wallpaper) {
	fyne.Do(func() {
		v.apply.SetText("Applying…")
		v.apply.Disable()
	})

	destDir := picturesDir()
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		fyne.Do(func() {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			v.apply.SetText("Apply")
			v.apply.Enable()
		})
		return
	}

	ext := filepath.Ext(wp.Path)
	if ext == "" {
		ext = ".jpg"
	}
	filename := "wallhaven-" + wp.ID + ext
	destPath := filepath.Join(destDir, filename)

	res, err := http.Get(wp.Path)
	if err != nil || res.StatusCode != http.StatusOK {
		if res != nil {
			res.Body.Close()
		}
		fyne.Do(func() {
			dialog.ShowError(fmt.Errorf("download failed: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
			v.apply.SetText("Apply")
			v.apply.Enable()
		})
		return
	}

	out, err := os.Create(destPath)
	if err != nil {
		res.Body.Close()
		fyne.Do(func() {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			v.apply.SetText("Apply")
			v.apply.Enable()
		})
		return
	}

	_, err = io.Copy(out, res.Body)
	res.Body.Close()
	out.Close()

	if err != nil {
		fyne.Do(func() {
			dialog.ShowError(fmt.Errorf("save failed: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
			v.apply.SetText("Apply")
			v.apply.Enable()
		})
		return
	}

	// Apply via gsettings (GNOME).
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", "file://"+destPath)
	if err := cmd.Run(); err != nil {
		fyne.Do(func() {
			dialog.ShowError(fmt.Errorf("could not apply wallpaper (gsettings): %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
		})
	} else {
		cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+destPath)
		cmd.Run()
		fyne.Do(func() {
			dialog.ShowInformation("Wallpaper applied", "Set to:\n"+destPath, fyne.CurrentApp().Driver().AllWindows()[0])
		})
	}

	fyne.Do(func() {
		v.apply.SetText("Apply")
		v.apply.Enable()
	})
}

func (v *ExpandedView) loadImage(id, fullURL string) {
	res, err := http.Get(fullURL)
	if err != nil || res.StatusCode != http.StatusOK {
		if res != nil {
			res.Body.Close()
		}
		return
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil || !decodeable(body) {
		return
	}
	fyne.Do(func() {
		v.image.Resource = fyne.NewStaticResource(id, body)
		v.image.Refresh()
		v.spinner.Hide()
	})
}

func (v *ExpandedView) loadDetail(id string) {
	detail, err := v.client.Wallpaper(id)
	if err != nil {
		return
	}
	fyne.Do(func() {
		v.populateInfo(detail.Data)
	})
}

func (v *ExpandedView) populateInfo(det wallhaven.WallpaperDetail) {
	v.info.Objects = nil
	v.info.Objects = append(v.info.Objects,
		infoRow("ID", det.ID),
		infoRow("Resolution", det.Resolution),
		infoRow("Ratio", det.Ratio),
		infoRow("Category", det.Category),
		infoRow("Purity", det.Purity),
		infoRow("File type", det.FileType),
		infoRow("File size", fmt.Sprintf("%d bytes", det.FileSize)),
		infoRow("Views", fmt.Sprintf("%d", det.Views)),
		infoRow("Favorites", fmt.Sprintf("%d", det.Favorites)),
		infoRow("Uploaded", det.CreatedAt.Time.Format("2006-01-02")),
	)
	if det.Source != "" {
		v.info.Objects = append(v.info.Objects, infoRow("Source", det.Source))
	}
	if det.Uploader.Username != "" {
		v.info.Objects = append(v.info.Objects, infoRow("Uploader", det.Uploader.Username))
	}
	if len(det.Tags) > 0 {
		var tags string
		for i, t := range det.Tags {
			if i > 0 {
				tags += ", "
			}
			tags += t.Name
		}
		v.info.Objects = append(v.info.Objects, infoRow("Tags", tags))
	}
	v.info.Refresh()
}

func infoRow(label, value string) *fyne.Container {
	return container.NewVBox(
		widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value),
	)
}

// picturesDir returns the XDG Pictures directory, falling back to ~/Pictures.
// GNOME and KDE both look for user wallpapers here by default.
func picturesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	// Parse ~/.config/user-dirs.dirs for XDG_PICTURES_DIR.
	cfgPath := filepath.Join(home, ".config", "user-dirs.dirs")
	if data, err := os.ReadFile(cfgPath); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "XDG_PICTURES_DIR=") {
				val := strings.Trim(strings.TrimPrefix(line, "XDG_PICTURES_DIR="), `"`)
				val = strings.ReplaceAll(val, "$HOME", home)
				return val
			}
		}
	}
	return filepath.Join(home, "Pictures")
}

// GrowImage animates the main image growing from a small size to full
// with an elastic curve and a fade-in.
func (v *ExpandedView) GrowImage(duration time.Duration) {
	v.image.Translucency = 1.0
	target := v.scroll.Size()
	if target.Width <= 1 || target.Height <= 1 {
		target = fyne.NewSize(800, 600)
	}

	anim := fyne.NewAnimation(duration, func(val float32) {
		e := elasticCurve(val)
		v.image.Translucency = 1.0 - float64(e)
		w := float32(1) + e*target.Width
		h := float32(1) + e*target.Height
		v.image.SetMinSize(fyne.NewSize(w, h))
		v.image.Refresh()
	})
	anim.Curve = fyne.AnimationEaseOut
	anim.Start()
}

// CreateRenderer lays out: close button top-right, large image center, info sidebar right.
// A background rectangle is stacked behind everything so the grid doesn't
// show through the transparent image areas.
func (v *ExpandedView) CreateRenderer() fyne.WidgetRenderer {
	v.bg.FillColor = theme.Color(theme.ColorNameBackground)
	v.bg.Refresh()

	infoScroll := container.NewScroll(v.info)
	infoScroll.SetMinSize(fyne.NewSize(250, 0))
	buttons := container.NewHBox(v.apply, v.download, v.close)
	topBar := container.NewBorder(nil, nil, nil, buttons, widget.NewLabel(""))
	content := container.NewBorder(topBar, nil, nil, infoScroll, container.NewPadded(v.scroll))
	stack := container.NewStack(v.bg, content)
	return widget.NewSimpleRenderer(stack)
}
