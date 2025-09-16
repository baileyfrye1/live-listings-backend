package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.userHandler.GetAgentById))
	defer server.Close()

	url := fmt.Sprintf("%s/agents/3", server.URL)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	expected := "{\"id\":\"3\",\"first_name\":\"Renee\",\"last_name\":\"Guzzo\",\"email\":\"guzzorenee@gmail.com\",\"created_at\":\"2025-08-15 14:20:51.962969+00\",\"updated_at\":\"2025-08-15 16:54:50.543905+00\"}"

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}

	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
