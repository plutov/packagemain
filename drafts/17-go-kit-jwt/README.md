# packagemain #17: Microservices with go-kit. Part 3. JWT

Install go-kit:
```bash
go get github.com/go-kit/kit/...
```

```go
// service/service.go
package service

import (
	"context"
	"errors"
)

// AuthService .
type AuthService struct{}

// Auth .
func (s AuthService) Auth(ctx context.Context, email string, pass string) (string, error) {
	panic(errors.New("not implemented"))
}

// Me .
func (s AuthService) Me(ctx context.Context) (string, error) {
	panic(errors.New("not implemented"))
}
```

```go
// endpoints/endpoints.go
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
```

```go
// http/http.go
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
```

```go
// cmd/authsvc/main.go
package main

import (
	"log"
	"os"

	stdhttp "net/http"

	jwt "github.com/dgrijalva/jwt-go"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/endpoints"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/http"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/service"
)

func main() {
	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY"))
	svc := service.AuthService{
		SigningKey: signingKey,
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) { return signingKey, nil }
	handler := http.NewHTTPHandler(endpoints.Endpoints{
		Auth: endpoints.MakeAuthEndpoint(svc),
		Me:   gokitjwt.NewParser(keyFunc, jwt.SigningMethodHS256, gokitjwt.StandardClaimsFactory)(endpoints.MakeMeEndpoint(svc)),
	})

	log.Fatal(stdhttp.ListenAndServe(":8888", handler))
}
```

```go
// service/service.go
package service

import (
	"context"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
)

// AuthService .
type AuthService struct {
	SigningKey []byte
}

var (
	validEmail = "alex@packagemain.com"
	validPass  = "12345"
)

// Auth .
func (s AuthService) Auth(ctx context.Context, email string, pass string) (string, error) {
	if email == validEmail && pass == validPass {
		return s.generateJWT(email)
	}

	return "", errors.New("invalid credentials")
}

// Me .
func (s AuthService) Me(ctx context.Context) (string, error) {
	if claims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*jwt.StandardClaims); ok {
		return claims.Id, nil
	}

	return "", errors.New("user not found")
}

func (s AuthService) generateJWT(email string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.SigningKey)
}
```

```bash
JWT_SIGNING_KEY=verysecretkey go run cmd/authsvc/main.go
```

```bash
curl -XPOST -H "Content-Type: application/json" -d '{"email": "alex@packagemain.com", "pass": "12345"}' http://localhost:8888/auth
> {
	"jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTc1NzcxMTMsImp0aSI6ImFsZXhAcGFja2FnZW1haW4uY29tIiwiaWF0IjoxNTU3NDkwNzEzfQ.D-p2K0UPQ5cE3BGdk8BgXUm-S4puNqdI6a4ZxecWASc",
	"error":""
}
```

```bash
curl -XGET http://localhost:8888/me
> token up for parsing was not passed through the context
```

```bash
curl -XGET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTc1NzcxMTMsImp0aSI6ImFsZXhAcGFja2FnZW1haW4uY29tIiwiaWF0IjoxNTU3NDkwNzEzfQ.D-p2K0UPQ5cE3BGdk8BgXUm-S4puNqdI6a4ZxecWASc" http://localhost:8888/me
> token up for parsing was not passed through the context
```