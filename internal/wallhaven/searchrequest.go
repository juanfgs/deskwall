package wallhaven

import (
	"errors"
	"net/url"
	"strconv"
)

var ErrNoPreviousPage = errors.New("no previous page")

type Color struct {
	R string
	G string
	B string
}

type SearchRequest struct {
	ApiKey      string
	AtLeast     []int32
	Categories  []string
	Colors      []Color
	Order       string
	Page        int
	Purity      string
	Query       string `json:"q"`
	Ratios      []int32
	Resolutions []int32
	Sorting     string
}

func NewSearchRequest(query string) SearchRequest {
	return SearchRequest{
		Query: query,
		Page:  1,
	}
}

func (r SearchRequest) NextPage() (SearchRequest, error) {
	r.Page++
	return r, nil
}

func (r SearchRequest) PrevPage() (SearchRequest, error) {
	if r.Page == 1 {
		return r, ErrNoPreviousPage
	}
	r.Page--
	return r, nil
}

func (r SearchRequest) Values() url.Values {
	v := url.Values{}
	v.Set("q", r.Query)
	v.Set("page", strconv.Itoa(r.Page))
	if r.Categories != nil {
		v.Set("categories", encodeToggle(r.Categories))
	}
	if r.Purity != "" {
		v.Set("purity", r.Purity)
	}
	if r.Sorting != "" {
		v.Set("sorting", r.Sorting)
	}
	if r.Order != "" {
		v.Set("order", r.Order)
	}
	return v
}

// encodeToggle joins a slice of "0"/"1" flags into a single string, e.g.
// ["1","0","1"] -> "101".
func encodeToggle(flags []string) string {
	out := ""
	for _, f := range flags {
		if f == "1" {
			out += "1"
		} else {
			out += "0"
		}
	}
	return out
}
