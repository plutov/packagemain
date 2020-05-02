Hi Gophers, My name is Alex Pliutau.

And today me and my buddy (Gopher toy) are going to make the second part of "Microservices with go-kit" series.

## Microservices with go-kit. Part 2

In the previous video we prepared a local environment for our services using kit command line tool. In this video we'll continue to work with this code.

I will copy all code to a new folder and replace the go imports.

```
cp -R ../12-go-kit-1/* .
find . -type f -print0 | xargs -0 sed -i "" "s/12-go-kit-1/13-go-kit-2/g"
```

Let's implement our Notificator service first by writing the proto definition as it's supposed to be a gRPC service. We aleady have pre-generated file `notificator/pkg/grpc/pb/notificator.pb`, let's make it really simple.

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

Now we need to generate server and client stubs, we can use the `compile.sh` script already given us by kit tool, it basically contains the `protoc` command.

```
cd notificator/pkg/grpc/pb
./compile.sh
```

If we check `notificator.pb.go` - it was updated.

Now we need to implement the service itself. Instead of sending a real email let's generate a uuid only and return it, pretending that it's sent. But first we have to edit a bit the service to match our Request / Response formats (new `id` return argument).

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
	id := uuid.NewV4()

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
	reply := r.(endpoint.SendEmailResponse)
	return &pb.SendEmailReply{Id: reply.Id}, nil
}
```

### Service discovery

The SendEmail will be invoked by User service, so User service needs to know the address of Notificator, the typical service discovery problem. Of course in our local environment we know how to connect to the service as we use Docker Compose, but it may be more difficult in real distributed environment.

Let's start with registering our Notificator service in the etcd. Basically etcd is a distributed reliable key-value store, widely used for service discovery. go-kit supports other technologies for service discovery: eureka, consul, zookeeper, etc.

Let's add it to our Docker Compose so it will be available for our servers. Copied from Internet:

docker-compose.yml:
```
etcd:
	image: 'quay.io/coreos/etcd:v3.1.7'
	restart: always
	ports:
		- '23791:2379'
		- '23801:2380'
	environment:
		ETCD_NAME: infra
		ETCD_INITIAL_ADVERTISE_PEER_URLS: 'http://etcd:2380'
		ETCD_INITIAL_CLUSTER: infra=http://etcd:2380
		ETCD_INITIAL_CLUSTER_STATE: new
		ETCD_INITIAL_CLUSTER_TOKEN: secrettoken
		ETCD_LISTEN_CLIENT_URLS: 'http://etcd:2379,http://localhost:2379'
		ETCD_LISTEN_PEER_URLS: 'http://etcd:2380'
		ETCD_ADVERTISE_CLIENT_URLS: 'http://etcd:2379'
```

Let's register Notificator in etcd, notificator/cmd/service/service.go:

```
registrar, err := registerService(logger)
if err != nil {
	logger.Log(err)
	return
}

defer registrar.Deregister()

func registerService(logger log.Logger) (*sdetcd.Registrar, error) {
	var (
		etcdServer = "http://etcd:2379"
		prefix     = "/services/notificator/"
		instance   = "notificator:8082"
		key        = prefix + instance
	)

	client, err := sdetcd.NewClient(context.Background(), []string{etcdServer}, sdetcd.ClientOptions{})
	if err != nil {
		return nil, err
	}

	registrar := sdetcd.NewRegistrar(client, sdetcd.Service{
		Key:   key,
		Value: instance,
	}, logger)

	registrar.Register()

	return registrar, nil
}
```

We should always remember to deregister service when our program is stopped or crashed. Now etcd knows about our service, in this example we have only 1 instance, but in real life it could be more of course.

Now let's test our Notificator service and check if it is able to register in etcd:

```
docker-compose up -d etcd
docker-compose up -d notificator
```

Now let's get back to our Users service and invoke the Notificator service, basically we're going to send a fictional notification to user after it's created.

As Notificator is a gRPC service, so we need to share a client stub file with our client, in our case Users service.

The protobuf client stub code is located in `notificator/pkg/grpc/pb/notificator.pb.go`, and we can just import this package to our cient.

users/pkg/service/service.go:
```
import (
	"github.com/plutov/packagemain/13-go-kit-2/notificator/pkg/grpc/pb"
	"google.golang.org/grpc"
)

type basicUsersService struct {
	notificatorServiceClient pb.NotificatorClient
}

func (b *basicUsersService) Create(ctx context.Context, email string) error {
	reply, err := b.notificatorServiceClient.SendEmail(context.Background(), &pb.SendEmailRequest{
		Email:   email,
		Content: "Hi! Thank you for registration...",
	})

	if reply != nil {
		log.Printf("Email ID: %s", reply.Id)
	}

	return err
}

// NewBasicUsersService returns a naive, stateless implementation of UsersService.
func NewBasicUsersService() UsersService {
	conn, err := grpc.Dial("notificator:8082", grpc.WithInsecure())
	if err != nil {
		log.Printf("unable to connect to notificator: %s", err.Error())
		return new(basicUsersService)
	}

	log.Printf("connected to notificator")

	return &basicUsersService{
		notificatorServiceClient: pb.NewNotificatorClient(conn),
	}
}
```

But as we registered Notificator in etcd we can replace hardcoded Notificator address by getting it from etcd.

```
var etcdServer = "http://etcd:2379"

client, err := sdetcd.NewClient(context.Background(), []string{etcdServer}, sdetcd.ClientOptions{})
if err != nil {
	log.Printf("unable to connect to etcd: %s", err.Error())
	return new(basicUsersService)
}

entries, err := client.GetEntries("/services/notificator/")
if err != nil || len(entries) == 0 {
	log.Printf("unable to get prefix entries: %s", err.Error())
	return new(basicUsersService)
}

conn, err := grpc.Dial(entries[0], grpc.WithInsecure())
```

We get the first entry as we have only one, but in real system it may be hundreads of entries, so we can apply some logic for instance selection, for example Round Robin.

Now let's start our Users service and test this out:

```
docker-compose up users
```

We're going to call the http endpoint to create a user:

```
curl -XPOST http://localhost:8802/create -d '{"email": "test"}'
```

### Conclusion

In this video we have implemented fictional Notificator gRPC service, registered it in etcd and invoked from another service Users.

In the next video we're going to review the service authorization through JWT (JSON Web Tokens).