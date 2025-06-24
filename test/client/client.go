package client

import (
	models "aslon1213/magazin_pos/pkg/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Client struct {
	Host     string
	Port     string
	Username string
	Passwrod string
	Token    string
}

func NewClient(host, port, username, password string) *Client {
	c := &Client{
		Host:     host,
		Port:     port,
		Username: username,
		Passwrod: password,
	}
	c.Token = c.Login()
	return c
}

func (c *Client) Login() string {
	return ""
}

func (c *Client) MakeRequest(method, path string, body []byte, headers map[string]string, auth_required bool) (*http.Response, error) {

	if auth_required {
		if c.Token == "" {
			c.Token = c.Login()
		}
		headers["Authorization"] = "Bearer " + c.Token
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, "http://"+c.Host+c.Port+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// if resp.StatusCode !=

	return resp, nil
}

func (c *Client) DecodeResponseMultiple(resp *http.Response, body []byte) (models.BranchFinanceOutput, error) {
	output := models.BranchFinanceOutput{}
	err := json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (c *Client) DecodeResponseSingle(resp *http.Response) (models.BranchFinanceOutputSingle, error) {
	output := models.BranchFinanceOutputSingle{}
	err := json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (c *Client) EncodeBody(body interface{}) ([]byte, error) {
	return json.Marshal(body)
}

// Get all branches
func (c *Client) GetAllBranches() (resp *http.Response, output models.BranchFinanceOutput, err error) {
	resp, err = c.MakeRequest("GET", "/finance/branches", nil, map[string]string{"Content-Type": "application/json"}, false)
	if err != nil {
		return nil, models.BranchFinanceOutput{}, err
	}
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, models.BranchFinanceOutput{}, err
	// }
	// log.Info().Str("resp", string(body)).Msg("Getting all branches")
	output, err = c.DecodeResponseMultiple(resp, nil)
	return resp, output, err
}

func (c *Client) GetBranchByID(branchID string) (resp *http.Response, output models.BranchFinanceOutputSingle, err error) {
	resp, err = c.MakeRequest("GET", fmt.Sprintf("/finance/id/%s", branchID), nil, map[string]string{"Content-Type": "application/json"}, false)
	if err != nil {
		return nil, models.BranchFinanceOutputSingle{}, err
	}

	output, err = c.DecodeResponseSingle(resp)
	return resp, output, err
}

func (c *Client) GetBranchByName(branchName string) (resp *http.Response, output models.BranchFinanceOutputSingle, err error) {
	resp, err = c.MakeRequest("GET", fmt.Sprintf("/finance/name/%s", branchName), nil, map[string]string{"Content-Type": "application/json"}, false)
	if err != nil {
		return nil, models.BranchFinanceOutputSingle{}, err
	}
	output, err = c.DecodeResponseSingle(resp)
	return resp, output, err
}

func (c *Client) CreateBranch(branch models.NewBranchFinanceInput) (resp *http.Response, output models.BranchFinanceOutputSingle, err error) {
	body, err := c.EncodeBody(branch)
	if err != nil {
		return nil, models.BranchFinanceOutputSingle{}, err
	}

	log.Info().Str("body", string(body)).Msg("Creating branch")

	resp, err = c.MakeRequest("POST", "/finance", body, map[string]string{"Content-Type": "application/json"}, false)
	if err != nil {
		return nil, models.BranchFinanceOutputSingle{}, err
	}
	output, err = c.DecodeResponseSingle(resp)
	return resp, output, err
}
