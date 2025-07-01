package client

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (c *Client) Login() (*http.Response, string, error) {

	body := map[string]string{
		"username": c.Username,
		"password": c.Passwrod,
	}

	// usmarshal body to json
	json_body, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal body")
		return nil, "", err
	}

	response, err := c.MakeRequest(
		"POST",
		"/auth/login",
		json_body,
		map[string]string{
			"Content-Type": "application/json",
		},
		true,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to login")
	}

	resp_body := map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&resp_body)
	if err != nil {
		return response, "", err
	}
	log.Info().Interface("resp_body", resp_body).Msg("resp_body")

	return response, resp_body["data"].(string), nil
}

func (c *Client) Register(username string, password string) (*http.Response, error) {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	// usmarshal body to json
	json_body, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal body")
		return nil, err
	}

	response, err := c.MakeRequest(
		"POST",
		"/auth/register",
		json_body,
		map[string]string{
			"Content-Type": "application/json",
		},
		false,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to register")
		return response, err
	}

	return response, nil
}
