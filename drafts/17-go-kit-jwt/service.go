package service

import "context"

type Service interface {
	Auth(ctx context.Context, c Credentials) (string, error)
	Me(ctx context.Context) (User, error)
}

type Credentials struct {
	Email string `json:"email,omitempty"`
	Pass  string `json:"pass,omitempty"`
}

type User struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}
