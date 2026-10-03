package wallhaven

import (
	"encoding/json"
	"testing"
	"time"
)

const wallpaperDetailJSON = `{
  "data": {
    "id": "94x38z",
    "url": "https://wallhaven.cc/w/94x38z",
    "short_url": "http://whvn.cc/94x38z",
    "uploader": {
      "username": "test-user",
      "group": "User",
      "avatar": {
        "200px": "https://wallhaven.cc/images/user/avatar/200/11_3339efb2a813.png",
        "128px": "https://wallhaven.cc/images/user/avatar/128/11_3339efb2a813.png",
        "32px": "https://wallhaven.cc/images/user/avatar/32/11_3339efb2a813.png",
        "20px": "https://wallhaven.cc/images/user/avatar/20/11_3339efb2a813.png"
      }
    },
    "views": 12,
    "favorites": 0,
    "source": "",
    "purity": "sfw",
    "category": "anime",
    "dimension_x": 6742,
    "dimension_y": 3534,
    "resolution": "6742x3534",
    "ratio": "1.91",
    "file_size": 5070446,
    "file_type": "image/jpeg",
    "created_at": "2018-10-31 01:23:10",
    "colors": ["#000000", "#abbcda", "#424153", "#66cccc", "#333399"],
    "path": "https://w.wallhaven.cc/full/94/wallhaven-94x38z.jpg",
    "thumbs": {
      "large": "https://th.wallhaven.cc/lg/94/94x38z.jpg",
      "original": "https://th.wallhaven.cc/orig/94/94x38z.jpg",
      "small": "https://th.wallhaven.cc/small/94/94x38z.jpg"
    },
    "tags": [
      {
        "id": 1,
        "name": "anime",
        "alias": "Chinese cartoons",
        "category_id": 1,
        "category": "Anime & Manga",
        "purity": "sfw",
        "created_at": "2015-01-16 02:06:45"
      }
    ]
  }
}`

func TestUnmarshalWallpaperDetail(t *testing.T) {
	t.Parallel()

	var got WallpaperDetailResponse
	if err := json.Unmarshal([]byte(wallpaperDetailJSON), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	d := got.Data
	if d.ID != "94x38z" {
		t.Errorf("ID = %q, want 94x38z", d.ID)
	}
	// Embedded Wallpaper fields
	if d.Views != 12 {
		t.Errorf("Views = %d, want 12", d.Views)
	}
	if d.FileSize != 5070446 {
		t.Errorf("FileSize = %d, want 5070446", d.FileSize)
	}
	wantTime := time.Date(2018, 10, 31, 1, 23, 10, 0, time.UTC)
	if !d.CreatedAt.Time.Equal(wantTime) {
		t.Errorf("CreatedAt = %v, want %v", d.CreatedAt.Time, wantTime)
	}
	if d.Path != "https://w.wallhaven.cc/full/94/wallhaven-94x38z.jpg" {
		t.Errorf("Path = %q", d.Path)
	}
	if d.Thumbs.Small != "https://th.wallhaven.cc/small/94/94x38z.jpg" {
		t.Errorf("Thumbs.Small = %q", d.Thumbs.Small)
	}

	// Uploader
	if d.Uploader.Username != "test-user" {
		t.Errorf("Uploader.Username = %q", d.Uploader.Username)
	}
	if d.Uploader.Group != "User" {
		t.Errorf("Uploader.Group = %q", d.Uploader.Group)
	}
	if d.Uploader.Avatar.Px200 != "https://wallhaven.cc/images/user/avatar/200/11_3339efb2a813.png" {
		t.Errorf("Uploader.Avatar.Px200 = %q", d.Uploader.Avatar.Px200)
	}
	if d.Uploader.Avatar.Px20 != "https://wallhaven.cc/images/user/avatar/20/11_3339efb2a813.png" {
		t.Errorf("Uploader.Avatar.Px20 = %q", d.Uploader.Avatar.Px20)
	}

	// Tags
	if len(d.Tags) != 1 {
		t.Fatalf("Tags len = %d, want 1", len(d.Tags))
	}
	tag := d.Tags[0]
	if tag.ID != 1 || tag.Name != "anime" {
		t.Errorf("Tag = %+v, want {ID:1 Name:anime}", tag)
	}
	if tag.Alias != "Chinese cartoons" {
		t.Errorf("Tag.Alias = %q", tag.Alias)
	}
	if tag.CategoryID != 1 {
		t.Errorf("Tag.CategoryID = %d, want 1", tag.CategoryID)
	}
	if tag.Category != "Anime & Manga" {
		t.Errorf("Tag.Category = %q", tag.Category)
	}
	if tag.Purity != "sfw" {
		t.Errorf("Tag.Purity = %q", tag.Purity)
	}
	wantTagTime := time.Date(2015, 1, 16, 2, 6, 45, 0, time.UTC)
	if !tag.CreatedAt.Time.Equal(wantTagTime) {
		t.Errorf("Tag.CreatedAt = %v, want %v", tag.CreatedAt.Time, wantTagTime)
	}
}
