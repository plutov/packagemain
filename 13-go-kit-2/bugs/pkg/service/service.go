package service

import "context"

// BugsService describes the service.
type BugsService interface {
	// Add your methods here
	Create(ctx context.Context, bug string) error
}

type basicBugsService struct{}

func (b *basicBugsService) Create(ctx context.Context, bug string) (e0 error) {
	// TODO implement the business logic of Create
	return e0
}

// NewBasicBugsService returns a naive, stateless implementation of BugsService.
func NewBasicBugsService() BugsService {
	return &basicBugsService{}
}

// New returns a BugsService with all of the expected middleware wired in.
func New(middleware []Middleware) BugsService {
	var svc BugsService = NewBasicBugsService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
