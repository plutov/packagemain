package endpoints

import "context"

import "github.com/go-kit/kit/endpoint"

import "github.com/plutov/packagemain/drafts/17-go-kit-jwt/service"

type AuthRequest struct {
	C service.Credentials `json:"credentials"`
}

type AuthResponse struct {
	JWT string `json:"jwt"`
	Err string `json:"err"`
}

func MakeAuthEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthRequest)
		string1, err := s.Auth(ctx, req.C)
		return AuthResponse{JWT: string1, Err: getErrorMsg(err)}, nil
	}
}

type MeRequest struct {
}
type MeResponse struct {
	Me  service.User `json:"me"`
	Err string       `json:"err"`
}

func MakeMeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		U, err := s.Me(ctx)
		return MeResponse{Me: U, Err: getErrorMsg(err)}, nil
	}
}

type Endpoints struct {
	Auth endpoint.Endpoint
	Me   endpoint.Endpoint
}

func getErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
