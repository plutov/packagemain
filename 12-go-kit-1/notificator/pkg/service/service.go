package service

import "context"

// NotificatorService describes the service.
type NotificatorService interface {
	// Add your methods here
	SendEmail(ctx context.Context, email string, content string) error
}

type basicNotificatorService struct{}

func (b *basicNotificatorService) SendEmail(ctx context.Context, email string, content string) (e0 error) {
	// TODO implement the business logic of SendEmail
	return e0
}

// NewBasicNotificatorService returns a naive, stateless implementation of NotificatorService.
func NewBasicNotificatorService() NotificatorService {
	return &basicNotificatorService{}
}

// New returns a NotificatorService with all of the expected middleware wired in.
func New(middleware []Middleware) NotificatorService {
	var svc NotificatorService = NewBasicNotificatorService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
