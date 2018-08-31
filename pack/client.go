package idex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client for requests
type Client struct {
	URL string
}

// NewClient instance of a client
func NewClient(url string) *Client {
	return &Client{URL: url}
}

func (c *Client) do(endpoint, payload string) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.URL, endpoint), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
