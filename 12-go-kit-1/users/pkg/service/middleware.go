package service

import (
	"context"

	log "github.com/go-kit/kit/log"
)

type Middleware func(UsersService) UsersService

type loggingMiddleware struct {
	logger log.Logger
	next   UsersService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next UsersService) UsersService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Create(ctx context.Context, email string) (err error) {
	defer func() {
		l.logger.Log("method", "Create", "email", email, "err", err)
	}()
	return l.next.Create(ctx, email)
}
