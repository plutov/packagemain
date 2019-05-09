package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/auth/jwt"

	stdjwt "github.com/dgrijalva/jwt-go"
)

type Credentials struct {
	Email string `json:"email,omitempty"`
	Pass  string `json:"pass,omitempty"`
}
type User struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}
type Service struct {
	SigningKey []byte
}

type Claims struct {
	Me User `json:"me"`
	stdjwt.StandardClaims
}

var (
	validEmail = "alex@packagemain.com"
	validPass  = "12345"
	validName  = "Alex"
)

func generateJWT(signingKey []byte, user User) (string, error) {
	claims := stdjwt.StandardClaims{
		Id:        user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  stdjwt.TimeFunc().Unix(),
	}
	token := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func (s Service) Auth(ctx context.Context, c Credentials) (string, error) {
	if c.Email == validEmail && c.Pass == validPass {
		return generateJWT(s.SigningKey, User{
			Email: validEmail,
			Name:  validName,
		})
	}

	return "", errors.New("invalid credentials")
}

func (s Service) Me(ctx context.Context) (User, error) {
	if claims, ok := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.StandardClaims); ok {
		return User{
			Email: claims.Id,
		}, nil
	}

	return User{}, errors.New("user not found")
}
