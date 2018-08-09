Hi Gophers, My name is Alex Pliutau.

Welcome to the "package main" new video.

## Microservices with go-kit. Part 2

In the previous video we prepared a local environment for our services using kit command line tool. In this video we'll continue to work with this code.

```
cp -R ../12-go-kit-1/* .
find . -type f -print0 | xargs -0 sed -i "" "s/12-go-kit-1/13-go-kit-2/g"
```

Let's implement our Notificator service first by writing the proto definition as it's supposed to be a gRPC service. We aleady have pre-generated file `notificator/pkg/grpc/pb/notificator.pb`:

```
syntax = "proto3";

package pb;

service Notificator {
	rpc SendEmail (SendEmailRequest) returns (SendEmailReply);
}

message SendEmailRequest {
	string email = 1;
	string content = 2;
}

message SendEmailReply {
	string id = 1;
}
```

Now we need to generate server and client stubs, we can use the `compile.sh` script already given us by kit tool.

```
cd notificator/pkg/grpc/pb
./compile.sh
```

If we check `notificator.pb.go` it is updated.

Now we need to implement the service itself, let's generate a uuid and return it. But first we have to edit a bit the service to match our Request / Response formats (new `id` return argument).

notificator/pkg/service/service.go:
```
import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// NotificatorService describes the service.
type NotificatorService interface {
	// Add your methods here
	SendEmail(ctx context.Context, email string, content string) (string, error)
}

type basicNotificatorService struct{}

func (b *basicNotificatorService) SendEmail(ctx context.Context, email string, content string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
```

notificator/pkg/service/middleware.go:
```
func (l loggingMiddleware) SendEmail(ctx context.Context, email string, content string) (string, error) {
	defer func() {
		l.logger.Log("method", "SendEmail", "email", email, "content", content)
	}()
	return l.next.SendEmail(ctx, email, content)
}
```

notificator/pkg/endpoint/endpoint.go
```
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
```

If we search for TODO `grep -R "TODO" notificator` we can see that we still need to implement Encoder and Decoder for gRPC request and response.

notificator/pkg/grpc/handler.go:
```
func decodeSendEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SendEmailRequest)
	return endpoint.SendEmailRequest{Email: req.Email, Content: req.Content}, nil
}

func encodeSendEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	reply := r.(*pb.SendEmailReply)
	return endpoint.SendEmailResponse{Id: reply.Id}, nil
}
```