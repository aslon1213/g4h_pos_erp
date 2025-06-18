package test

import (
	"aslon1213/magazin_pos/pkg/app"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func StartTest(m *testing.M) {
	app := app.NewApp()
	app.Run()
	os.Exit(m.Run())
}

func TestCreateSuppliers(t *testing.T) {

	app := app.NewApp()
	go app.Run()
	// Create test HTTP client
	client := &http.Client{}

	// Test GetSuppliers
	resp, err := client.Get("/suppliers")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test GetSupplierByID
	resp, err = client.Get("/suppliers/123")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 404 { // Expect 404 since supplier doesn't exist
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}

	// Test CreateSupplier
	resp, err = client.Post("/suppliers", "application/json", strings.NewReader(`{
		"name": "Test Supplier",
		"email": "test@example.com",
		"phone": "1234567890",
		"address": "Test Address",
		"inn": "123456789",
		"notes": "Test Notes",
		"branch": "Test Branch"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}

	// Parse response to get created supplier ID
	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatal(err)
	}
	supplierData := createResp["data"].(map[string]interface{})
	supplierID := supplierData["id"].(string)

	// Test UpdateSupplier
	req, err := http.NewRequest("PUT", "/suppliers/"+supplierID, strings.NewReader(`{
		"name": "Updated Test Supplier"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test NewTransaction
	resp, err = client.Post("/suppliers/"+supplierID+"/transactions", "application/json", strings.NewReader(`{
		"amount": 1000,
		"type": "credit",
		"description": "Test Transaction"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Test DeleteSupplier
	req, err = http.NewRequest("DELETE", "/suppliers/"+supplierID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

}
