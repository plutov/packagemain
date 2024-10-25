package demo

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Response struct {
	Results    []int `json:"items"`
	PagesCount int   `json:"pagesCount"`
}

func ProcessRequest(w http.ResponseWriter, r *http.Request) {
	var req Request

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply offset and limit to some static data
	if req.Limit <= 0 {
		req.Limit = 1
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	start := req.Offset
	end := req.Offset + req.Limit
	all := make([]int, 1000)
	if req.Offset > len(all) {
		start = len(all) - 1
	}
	if end > len(all) {
		end = len(all)
	}
	res := Response{
		PagesCount: len(all) / req.Limit,
		Results:    all[start:end],
	}

	// Send JSON response
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
