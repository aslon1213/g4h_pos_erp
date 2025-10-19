package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Tests ----

func TestArrivals(t *testing.T) {
	ChangeLogging()

	// configs := getConfig(t)

	client := getClient(t)

	branches := []string{"xonobod", "polevoy"}

	test_arrivals_branch_xonobod := []string{
		"Arrival 1",
		"Arrival 2",
		"Arrival 3",
	}

	test_arrivals_branch_polevoy := []string{
		"Arrival 4",
		"Arrival 5",
		"Arrival 6",
	}

	for _, branch := range branches {
		if branch == "xonobod" {
			resp, output, err := client.NewArrivals(branch, test_arrivals_branch_xonobod)
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
			assert.Nil(t, err, "Expected no error, but got one")
			assert.Equal(t, "Proposal created successfully", output["message"], "Expected message to be 'Proposal created successfully', but got %s", output["message"])
		}
		if branch == "polevoy" {
			resp, output, err := client.NewArrivals(branch, test_arrivals_branch_polevoy)
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
			assert.Nil(t, err, "Expected no error, but got one")
			assert.Equal(t, "Proposal created successfully", output["message"], "Expected message to be 'Proposal created successfully', but got %s", output["message"])
		}
	}

	// get from backend and ensure that the arrivals are created
	for _, branch := range branches {
		resp, output, err := client.GetArrivals(branch)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, but got %d", resp.StatusCode)
		assert.Nil(t, err, "Expected no error, but got one")
		assert.Equal(t, len(test_arrivals_branch_xonobod), len(output), "Expected %d arrivals, but got %d", len(test_arrivals_branch_xonobod), len(output))
	}

}
