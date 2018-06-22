package service

import (
	"context"
)

// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, email string) (err error)
}

type basicUsersService struct{}

func (b *basicUsersService) Create(ctx context.Context, email string) (err error) {
	return err
}

// NewBasicUsersService returns a naive, stateless implementation of UsersService.
func NewBasicUsersService() UsersService {
	return &basicUsersService{}
}

// New returns a UsersService with all of the expected middleware wired in.
func New(middleware []Middleware) UsersService {
	var svc UsersService = NewBasicUsersService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
