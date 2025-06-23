package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AuthResponse represents the structure of the authentication response
type AuthResponse struct {
	Token string `json:"token"`
}

// LoginAndGetToken is a utility function to authenticate and retrieve a token
func LoginAndGetToken(t *testing.T, username, password string) string {
	// Define the login payload
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}

	// Marshal the payload to JSON
	body, err := json.Marshal(loginPayload)
	assert.NoError(t, err)

	// Create a new HTTP request for login
	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Use http.DefaultClient to send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check if the response status is OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response body to get the token
	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NoError(t, err)

	// Return the token
	return authResp.Token
}
