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
	return v
}
