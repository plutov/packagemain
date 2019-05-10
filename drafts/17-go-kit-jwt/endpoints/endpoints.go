package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/service"
)

// Endpoints .
type Endpoints struct {
	Auth endpoint.Endpoint
	Me   endpoint.Endpoint
}

// AuthRequest .
type AuthRequest struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

// AuthResponse .
type AuthResponse struct {
	JWT string `json:"jwt"`
	Err string `json:"error,omitempty"`
}

// MakeAuthEndpoint .
func MakeAuthEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthRequest)
		jwt, err := s.Auth(ctx, req.Email, req.Pass)
		if err != nil {
			return AuthResponse{Err: err.Error()}, nil
		}
		return AuthResponse{JWT: jwt}, nil
	}
}

// MeResponse .
type MeResponse struct {
	Email string `json:"email"`
	Err   string `json:"error,omitempty"`
}

// MakeMeEndpoint .
func MakeMeEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		email, err := s.Me(ctx)
		if err != nil {
			return MeResponse{Err: err.Error()}, nil
		}
		return MeResponse{Email: email}, nil
	}
}
