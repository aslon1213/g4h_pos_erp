package client

import (
	"encoding/json"
	"net/http"
)

func (c *Client) NewArrivals(branch_id string, arrival []string) (resp *http.Response, output map[string]interface{}, err error) {
	body, err := json.Marshal(arrival)
	if err != nil {
		return nil, output, err
	}
	resp, err = c.MakeRequest("POST", "/api/arrivals/"+branch_id, body, map[string]string{"Content-Type": "application/json"}, false)
	if err != nil {
		return nil, output, err
	}

	// decode a function
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return nil, output, err
	}

	return resp, output, err
}

func (c *Client) GetArrivals(branch_id string) (resp *http.Response, output []map[string]interface{}, err error) {
	resp, err = c.MakeRequest("GET", "/api/arrivals/"+branch_id, nil, map[string]string{"Content-Type": "application/json"}, false)
	if err != nil {
		return nil, output, err
	}
	return resp, output, err
}
