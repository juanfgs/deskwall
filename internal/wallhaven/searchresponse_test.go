package wallhaven

import (
	"encoding/json"
	"testing"
	"time"
)

// searchResponseJSON is a well-formed subset of the API response sample.
const searchResponseJSON = `{
  "data": [
    {
      "id": "94x38z",
      "url": "https://wallhaven.cc/w/94x38z",
      "short_url": "http://whvn.cc/94x38z",
      "views": 6,
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
      "path": "https://w.wallhaven.cc/94/wallhaven-94x38z.jpg",
      "thumbs": {
        "large": "https://th.wallhaven.cc/lg/94/94x38z.jpg",
        "original": "https://th.wallhaven.cc/orig/94/94x38z.jpg",
        "small": "https://th.wallhaven.cc/small/94/94x38z.jpg"
      }
    },
    {
      "id": "ze1p56",
      "url": "https://wallhaven.cc/w/ze1p56",
      "short_url": "http://whvn.cc/ze1p56",
      "views": 11,
      "favorites": 0,
      "source": "",
      "purity": "sfw",
      "category": "anime",
      "dimension_x": 3779,
      "dimension_y": 2480,
      "resolution": "3779x2480",
      "ratio": "1.52",
      "file_size": 1011043,
      "file_type": "image/jpeg",
      "created_at": "2018-10-07 17:05:28",
      "colors": ["#424153", "#e7d8b1", "#cc3333", "#ffffff", "#cccccc"],
      "path": "https://w.wallhaven.cc/ze/wallhaven-ze1p56.jpg",
      "thumbs": {
        "large": "https://th.wallhaven.cc/lg/ze/ze1p56.jpg",
        "original": "https://th.wallhaven.cc/orig/ze/ze1p56.jpg",
        "small": "https://th.wallhaven.cc/small/ze/ze1p56.jpg"
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "last_page": 36,
    "per_page": 24,
    "total": 848,
    "query": "test",
    "seed": "abc123"
  }
}`

const metaTagQueryJSON = `{
  "data": [],
  "meta": {
    "current_page": 1,
    "last_page": 1,
    "per_page": 24,
    "total": 1,
    "query": {"id": 1, "tag": "anime"},
    "seed": null
  }
}`

const metaNullQueryJSON = `{
  "data": [],
  "meta": {
    "current_page": 1,
    "last_page": 1,
    "per_page": 24,
    "total": 0,
    "query": null,
    "seed": null
  }
}`

func TestUnmarshalSearchResponse(t *testing.T) {
	t.Parallel()

	var got SearchResponse
	if err := json.Unmarshal([]byte(searchResponseJSON), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(got.Data) != 2 {
		t.Fatalf("want 2 wallpapers, got %d", len(got.Data))
	}

	w0 := got.Data[0]
	if w0.ID != "94x38z" {
		t.Errorf("ID = %q, want 94x38z", w0.ID)
	}
	if w0.FileSize != 5070446 {
		t.Errorf("FileSize = %d, want 5070446", w0.FileSize)
	}
	wantTime := time.Date(2018, 10, 31, 1, 23, 10, 0, time.UTC)
	if !w0.CreatedAt.Time.Equal(wantTime) {
		t.Errorf("CreatedAt = %v, want %v", w0.CreatedAt.Time, wantTime)
	}
	if len(w0.Colors) != 5 || w0.Colors[0] != "#000000" {
		t.Errorf("Colors = %v", w0.Colors)
	}
	if w0.Thumbs.Large != "https://th.wallhaven.cc/lg/94/94x38z.jpg" {
		t.Errorf("Thumbs.Large = %q", w0.Thumbs.Large)
	}

	if got.Meta.Total != 848 {
		t.Errorf("Meta.Total = %d, want 848", got.Meta.Total)
	}
	if got.Meta.Query.String == nil || *got.Meta.Query.String != "test" {
		t.Errorf("Meta.Query.String = %v, want \"test\"", got.Meta.Query.String)
	}
	if got.Meta.Query.Tag != nil {
		t.Errorf("Meta.Query.Tag = %v, want nil", got.Meta.Query.Tag)
	}
	if got.Meta.Seed == nil || *got.Meta.Seed != "abc123" {
		t.Errorf("Meta.Seed = %v, want \"abc123\"", got.Meta.Seed)
	}
}

func TestUnmarshalTagQuery(t *testing.T) {
	t.Parallel()

	var got SearchResponse
	if err := json.Unmarshal([]byte(metaTagQueryJSON), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Meta.Query.String != nil {
		t.Errorf("Query.String = %v, want nil", got.Meta.Query.String)
	}
	if got.Meta.Query.Tag == nil {
		t.Fatal("Query.Tag = nil, want object")
	}
	if got.Meta.Query.Tag.ID != 1 || got.Meta.Query.Tag.Tag != "anime" {
		t.Errorf("Query.Tag = %+v, want {ID:1 Tag:anime}", got.Meta.Query.Tag)
	}
	if got.Meta.Seed != nil {
		t.Errorf("Seed = %v, want nil", got.Meta.Seed)
	}
}

func TestUnmarshalNullQuery(t *testing.T) {
	t.Parallel()

	var got SearchResponse
	if err := json.Unmarshal([]byte(metaNullQueryJSON), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Meta.Query.String != nil || got.Meta.Query.Tag != nil {
		t.Errorf("Query = %+v, want both nil", got.Meta.Query)
	}
	if got.Meta.Seed != nil {
		t.Errorf("Seed = %v, want nil", got.Meta.Seed)
	}
}

func TestQueryRoundTrip(t *testing.T) {
	t.Parallel()

	cases := []string{searchResponseJSON, metaTagQueryJSON, metaNullQueryJSON}
	for i, in := range cases {
		var got SearchResponse
		if err := json.Unmarshal([]byte(in), &got); err != nil {
			t.Fatalf("case %d unmarshal: %v", i, err)
		}
		out, err := json.Marshal(got)
		if err != nil {
			t.Fatalf("case %d marshal: %v", i, err)
		}
		var rt SearchResponse
		if err := json.Unmarshal(out, &rt); err != nil {
			t.Fatalf("case %d re-unmarshal: %v", i, err)
		}
		// Compare query shape after round-trip.
		if (got.Meta.Query.String == nil) != (rt.Meta.Query.String == nil) {
			t.Errorf("case %d: Query.String nil mismatch", i)
		}
		if (got.Meta.Query.Tag == nil) != (rt.Meta.Query.Tag == nil) {
			t.Errorf("case %d: Query.Tag nil mismatch", i)
		}
		if got.Meta.Query.Tag != nil && rt.Meta.Query.Tag != nil {
			if *got.Meta.Query.Tag != *rt.Meta.Query.Tag {
				t.Errorf("case %d: Query.Tag = %+v, want %+v", i, rt.Meta.Query.Tag, got.Meta.Query.Tag)
			}
		}
	}
}
