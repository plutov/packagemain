package http

import (
	"context"
	"encoding/json"

	"net/http"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/endpoints"
)

// NewHTTPHandler .
func NewHTTPHandler(endpoints endpoints.Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/auth", httptransport.NewServer(endpoints.Auth, DecodeAuthRequest, EncodeJSONResponse))
	m.Handle("/me", httptransport.NewServer(endpoints.Me, DecodeMeRequest, EncodeJSONResponse, httptransport.ServerBefore(gokitjwt.HTTPToContext())))
	return m
}

// DecodeAuthRequest .
func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// EncodeJSONResponse .
func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// DecodeMeRequest .
func DecodeMeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
