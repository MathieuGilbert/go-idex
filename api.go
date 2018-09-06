package idex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// API for requests
type API struct {
	URL string
}

// Post returns the result of a POST to the endpoint with the payload
func (a *API) Post(endpoint, payload string) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", a.URL, endpoint), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
