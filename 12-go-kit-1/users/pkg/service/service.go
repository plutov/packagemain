package service

import (
	"context"
	"os"

	log "github.com/go-kit/kit/log"
)

// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, email string) (err error)
}

type basicUsersService struct {
	logger log.Logger
}

func (b *basicUsersService) Create(ctx context.Context, email string) (err error) {
	b.logger.Log("created user with email", email)
	return err
}

// NewBasicUsersService returns a naive, stateless implementation of UsersService.
func NewBasicUsersService() UsersService {
	return &basicUsersService{
		logger: log.NewJSONLogger(os.Stderr),
	}
}

// New returns a UsersService with all of the expected middleware wired in.
func New(middleware []Middleware) UsersService {
	var svc UsersService = NewBasicUsersService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
