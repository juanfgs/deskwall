package wallhaven

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
)

const BaseUrl = "https://wallhaven.cc/api/v1"

type Client struct {
	ApiKey string
	http.Client
}

func NewClient() Client {
	return Client{}
}

func NewClientWithApiKey(key string) Client {
	return Client{
		ApiKey: key,
	}
}

func buildUrl(endpoint string, query string) string {
	return fmt.Sprintf("%s/%s?%s", BaseUrl, endpoint, query)
}

// buildUrlPath builds a URL with a path segment, e.g. buildUrlPath("w", "94x38z")
// produces "https://wallhaven.cc/api/v1/w/94x38z".
func buildUrlPath(endpoint, id string) string {
	return fmt.Sprintf("%s/%s/%s", BaseUrl, endpoint, id)
}

func (c Client) Search(q string) (SearchResponse, error) {
	return c.SearchWithRequest(NewSearchRequest(q))
}

// SearchWithRequest performs a search using the given SearchRequest, allowing
// callers to control the page number and (future) filter parameters.
func (c Client) SearchWithRequest(r SearchRequest) (SearchResponse, error) {
	var response SearchResponse
	queryParams := r.Values()
	if c.ApiKey != "" {
		queryParams.Set("apikey", c.ApiKey)
	}

	req, err := http.NewRequest("GET", buildUrl("search", queryParams.Encode()), nil)
	if err != nil {
		log.Printf(" Error: %v ", err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Printf(" Error: %v ", resp)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return SearchResponse{}, err
	}
	return response, nil
}

// Wallpaper fetches the detail of a single wallpaper by its id, e.g. "94x38z".
func (c Client) Wallpaper(id string) (WallpaperDetailResponse, error) {
	var response WallpaperDetailResponse
	req, err := http.NewRequest("GET", buildUrlPath("w", id), nil)
	if err != nil {
		log.Printf(" Error: %v ", err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Printf(" Error: %v ", resp)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return WallpaperDetailResponse{}, err
	}
	return response, nil
}
