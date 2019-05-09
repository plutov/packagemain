# packagemain #17: Microservices with go-kit. Part 3. JWT

Install go-kit with kitgen:
```bash
go get github.com/go-kit/kit/...
go get github.com/go-kit/kit/cmd/kitgen
```

Prepare service interface:
```go
// service.go
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
```

Generate service:
```bash
kitgen ./service.go
```

Create cmd/authsvc/main.go file:
```go
// cmd/authsvc/main.go
package main

import (
	"log"

	stdhttp "net/http"

	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/endpoints"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/http"
	"github.com/plutov/packagemain/drafts/17-go-kit-jwt/service"
)

func main() {
	svc := service.Service{}

	handler := http.NewHTTPHandler(endpoints.Endpoints{
		PostProfile: endpoints.MakePostProfileEndpoint(svc),
	})

	log.Fatal(stdhttp.ListenAndServe(":8888", handler))
}
```