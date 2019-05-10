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
