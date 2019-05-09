package main

import (
	"log"

	stdhttp "net/http"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/endpoints"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/http"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/service"
)

var signingKey = []byte("verysecretkey")

func main() {
	svc := service.Service{
		SigningKey: signingKey,
	}

	keyFunc := func(token *stdjwt.Token) (interface{}, error) { return signingKey, nil }
	handler := http.NewHTTPHandler(endpoints.Endpoints{
		Auth: endpoints.MakeAuthEndpoint(svc),
		Me:   jwt.NewParser(keyFunc, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory)(endpoints.MakeMeEndpoint(svc)),
	})

	log.Fatal(stdhttp.ListenAndServe(":8888", handler))
}
