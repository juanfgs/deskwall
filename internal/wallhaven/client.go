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
