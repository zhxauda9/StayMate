package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080/api"

func TestRegister(t *testing.T) {
	registerPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	t.Logf("Trying to send request to register user, data: %v", registerPayload)
	respBody, err := sendRequest("POST", baseURL+"/register", registerPayload, nil)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	fmt.Printf("Register Response: %s\n", string(respBody))
}

func TestLogin(t *testing.T) {
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	respBody, err := sendRequest("POST", baseURL+"/login", loginPayload, nil)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	fmt.Printf("Login Response: %s\n", string(respBody))

	// Extract token from login response
	var loginResult map[string]string
	if err := json.Unmarshal(respBody, &loginResult); err != nil {
		t.Fatalf("Error parsing login response: %v", err)
	}

	token := loginResult["token"]
	if token == "" {
		t.Fatalf("Token is empty in login response")
	}
}

func TestValidateToken(t *testing.T) {
	token := loginAndGetToken(t)

	respBody, err := sendRequest("GET", baseURL+"/validate", nil, &token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	fmt.Printf("Validate Response: %s\n", string(respBody))
}

func TestGetProfile(t *testing.T) {
	token := loginAndGetToken(t)

	respBody, err := sendRequest("GET", baseURL+"/profile", nil, &token)
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}
	fmt.Printf("Profile Response: %s\n", string(respBody))
}

func TestLogout(t *testing.T) {
	token := loginAndGetToken(t)

	respBody, err := sendRequest("POST", baseURL+"/logout", nil, &token)
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}
	fmt.Printf("Logout Response: %s\n", string(respBody))
}

// Helper Functions

func loginAndGetToken(t *testing.T) string {
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	respBody, err := sendRequest("POST", baseURL+"/login", loginPayload, nil)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	var loginResult map[string]string
	if err := json.Unmarshal(respBody, &loginResult); err != nil {
		t.Fatalf("Error parsing login response: %v", err)
	}

	token := loginResult["token"]
	if token == "" {
		t.Fatalf("Token is empty in login response")
	}

	return token
}

func sendRequest(method, url string, payload interface{}, token *string) ([]byte, error) {
	var body []byte
	var err error

	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error response: %s", string(respBody))
	}

	return respBody, nil
}
