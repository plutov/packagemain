package demo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func FuzzProcessRequest(f *testing.F) {
	// Create sample inputs for the fuzzer
	testRequests := []Request{
		{Limit: -100, Offset: -100},
		{Limit: 0, Offset: 0},
		{Limit: 50, Offset: 50},
		{Limit: 2000, Offset: 2000},
	}

	// Add to the seed corpus
	for _, r := range testRequests {
		if data, err := json.Marshal(r); err == nil {
			f.Add(data)
		}
	}

	// Create a test server
	srv := httptest.NewServer(http.HandlerFunc(ProcessRequest))
	defer srv.Close()

	// Fuzz target with a single []byte argument
	f.Fuzz(func(t *testing.T, data []byte) {
		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			// Skip invalid JSON requests that may be generated during fuzz
			t.Skip("invalid json")
		}

		// Pass data to the server
		resp, err := http.DefaultClient.Post(srv.URL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("unable to call server: %v, data: %s", err, string(data))
		}

		defer resp.Body.Close()

		// Skip BadRequest errors
		if resp.StatusCode == http.StatusBadRequest {
			t.Skip("invalid json")
		}

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("non-200 status code %d", resp.StatusCode)
		}
	})
}
