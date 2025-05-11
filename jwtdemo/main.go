package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

func init() {
	publicKeyData, err := os.ReadFile("./public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyData, err := os.ReadFile("./private_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/login", login)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaims)
		},
		SigningKey:    publicKey,
		SigningMethod: jwt.SigningMethodRS256.Name,
	}

	g := e.Group("/api")
	g.Use(echojwt.WithConfig(config))
	g.GET("/greet", greet)

	e.Start("127.0.0.1:4242")
}

type jwtClaims struct {
	jwt.RegisteredClaims
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "package" || password != "main" {
		return echo.ErrUnauthorized
	}

	claims := &jwtClaims{
		jwt.RegisteredClaims{
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	t, err := token.SignedString(privateKey)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func greet(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtClaims)
	sub := claims.Subject
	return c.String(http.StatusOK, fmt.Sprintf("hi %s!", sub))
}
