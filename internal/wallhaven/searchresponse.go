package wallhaven

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// SearchResponse is the top-level payload returned by a wallhaven search.
type SearchResponse struct {
	Data []Wallpaper `json:"data"`
	Meta Meta        `json:"meta"`
}

// Wallpaper is a single wallpaper entry inside the data array.
type Wallpaper struct {
	ID         string  `json:"id"`
	URL        string  `json:"url"`
	ShortURL   string  `json:"short_url"`
	Views      int     `json:"views"`
	Favorites  int     `json:"favorites"`
	Source     string  `json:"source"`
	Purity     string  `json:"purity"`
	Category   string  `json:"category"`
	DimensionX int     `json:"dimension_x"`
	DimensionY int     `json:"dimension_y"`
	Resolution string  `json:"resolution"`
	Ratio      string  `json:"ratio"`
	FileSize   int64   `json:"file_size"`
	FileType   string  `json:"file_type"`
	CreatedAt  WHTime  `json:"created_at"`
	Colors     []string `json:"colors"`
	Path       string  `json:"path"`
	Thumbs     Thumbs  `json:"thumbs"`
}

// Thumbs groups the available thumbnail URLs for a wallpaper.
type Thumbs struct {
	Large    string `json:"large"`
	Original string `json:"original"`
	Small    string `json:"small"`
}

// Meta holds pagination info plus the polymorphic query and a nullable seed.
type Meta struct {
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	PerPage     int    `json:"per_page"`
	Total       int    `json:"total"`
	Query       Query  `json:"query"`
	Seed        *string `json:"seed"`
}

// TagQuery is the object form of meta.query, used for exact tag searches.
type TagQuery struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}

// Query represents meta.query, which the API encodes as one of:
//   - null
//   - a plain JSON string (e.g. "test")
//   - an object (e.g. {"id":1,"tag":"anime"})
// Exactly one of String or Tag is non-nil after a successful unmarshal.
type Query struct {
	String *string
	Tag    *TagQuery
}

// UnmarshalJSON decodes the polymorphic query field.
func (q *Query) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		q.String = &s
		return nil
	}
	var t TagQuery
	if err := json.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("wallhaven: query must be null, string, or tag object: %w", err)
	}
	q.Tag = &t
	return nil
}

// MarshalJSON re-encodes the query in its original polymorphic form.
func (q Query) MarshalJSON() ([]byte, error) {
	switch {
	case q.String != nil:
		return json.Marshal(*q.String)
	case q.Tag != nil:
		return json.Marshal(*q.Tag)
	default:
		return []byte("null"), nil
	}
}

// WHTime is a time.Time that (un)marshals using the wallhaven layout.
type WHTime struct {
	time.Time
}

const whTimeLayout = "2006-01-02 15:04:05"

// UnmarshalJSON parses created_at values like "2018-10-31 01:23:10".
func (t *WHTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	parsed, err := time.ParseInLocation(whTimeLayout, s, time.UTC)
	if err != nil {
		return fmt.Errorf("wallhaven: bad created_at %q: %w", s, err)
	}
	t.Time = parsed
	return nil
}

// MarshalJSON re-encodes the time in the wallhaven layout.
func (t WHTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return json.Marshal(t.Time.Format(whTimeLayout))
}
