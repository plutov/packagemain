package endpoint

import (
	"context"

	endpoint "github.com/go-kit/kit/endpoint"
	service "github.com/plutov/packagemain/13-go-kit-2/notificator/pkg/service"
)

// SendEmailRequest collects the request parameters for the SendEmail method.
type SendEmailRequest struct {
	Email   string `json:"email"`
	Content string `json:"content"`
}

// SendEmailResponse collects the response parameters for the SendEmail method.
type SendEmailResponse struct {
	Id string
	E0 error `json:"e0"`
}

// MakeSendEmailEndpoint returns an endpoint that invokes SendEmail on the service.
func MakeSendEmailEndpoint(s service.NotificatorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendEmailRequest)
		id, e0 := s.SendEmail(ctx, req.Email, req.Content)
		return SendEmailResponse{Id: id, E0: e0}, nil
	}
}

// Failed implements Failer.
func (r SendEmailResponse) Failed() error {
	return r.E0
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// SendEmail implements Service. Primarily useful in a client.
func (e Endpoints) SendEmail(ctx context.Context, email string, content string) (e0 error) {
	request := SendEmailRequest{
		Content: content,
		Email:   email,
	}
	response, err := e.SendEmailEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SendEmailResponse).E0
}
