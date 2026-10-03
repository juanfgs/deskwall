package wallhaven

import (
	//	"encoding/json"
	"errors"

	"testing"
)

func TestNextPage(t *testing.T) {
	t.Parallel()
	var query string = "foo"
	var firstPage int = 1
	r := NewSearchRequest(query)
	if r.Query != query {
		t.Errorf("r.Query = %s. want %s", r.Query, query)
	}

	if r.Page != firstPage {
		t.Errorf("r.Query = %d. want %d", r.Page, firstPage)
	}

	nextPage, _ := r.NextPage()
	if nextPage.Page != (firstPage + 1) {
		t.Errorf("r.Query = %d. want %d", nextPage.Page, firstPage+1)
	}

}

func TestPrevPage(t *testing.T) {
	t.Parallel()
	var query string = "foo"
	var secondPage int = 2
	r := NewSearchRequest(query)
	r.Page = secondPage
	if r.Query != query {
		t.Errorf("r.Query = %s. want %s", r.Query, query)
	}

	prevPage, _ := r.PrevPage()
	if prevPage.Page != (secondPage - 1) {
		t.Errorf("r.Query = %d. want %d", prevPage.Page, secondPage-1)
	}
	_, err := prevPage.PrevPage()
	if !errors.Is(err, ErrNoPreviousPage) {
		t.Fatalf("expected ErrNoPreviousPage, got %v", err)
	}

}

func TestValues(t *testing.T) {
	t.Parallel()
	r := NewSearchRequest("foo")

	got := r.Values()
	if got.Get("q") != "foo" {
		t.Errorf("expected query to be %s got %v", "foo", got.Get("q"))
	}
	if got.Get("page") != "1" {
		t.Errorf("expected query to be %s got %s", "1", got.Get("page"))
	}
}
