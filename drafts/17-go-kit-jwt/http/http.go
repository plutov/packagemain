package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/endpoints"
)

func NewHTTPHandler(endpoints endpoints.Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/auth", httptransport.NewServer(endpoints.Auth, DecodeAuthRequest, EncodeAuthResponse))
	m.Handle("/me", httptransport.NewServer(endpoints.Me, DecodeMeRequest, EncodeMeResponse, httptransport.ServerBefore(jwt.HTTPToContext())))
	return m
}
func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func EncodeAuthResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func DecodeMeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
func EncodeMeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
