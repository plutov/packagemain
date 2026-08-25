package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plutov/gopkg/pkgsiteapi"
)

func TestDetailCmdIncludesPseudoVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/versions/example.com/module":
			if got := r.URL.Query().Get("pseudo"); got != "true" {
				t.Errorf("pseudo query parameter = %q, want true", got)
			}
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/symbols/example.com/module":
			_, _ = w.Write([]byte(`{"symbols":{"items":[]}}`))
		default:
			t.Errorf("unexpected request path: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := pkgsiteapi.NewClientWithResponses(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	message := detailCmd(client, &pkgsiteapi.SearchResultData{PackagePath: "example.com/module"})()
	if _, ok := message.(detailMsg); !ok {
		t.Fatalf("detailCmd() message = %T, want detailMsg", message)
	}
}
