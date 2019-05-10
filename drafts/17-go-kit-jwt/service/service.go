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
