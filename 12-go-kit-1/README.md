Hi Gophers, My name is Alex Pliutau.

Welcome to the "package main" new video.

## Microservices with go-kit. Part 1

Before we talk about go-kit, I wanna share some good news with you.

If you're from Vietnam and watching this video, we're organizing the first GopherCon in Vietnam, Ho Chi Minh city. It will be hold in November. Sponsors and speakers are welcomed, you may find all information on gophercon.vn.

Nowadays, Microservices is one of the most popular buzz-word in the field of software architecture.

There are different definitions of the word "microservice", I like to say that Microservice is what  single programmer can design, implement, deploy, and maintain.

In a monolithic application, components invoke one another via language‑level method or function calls. In contrast, a microservices‑based application is a distributed system running on multiple machines. Each service instance is typically a process. So services must interact using an inter‑process communication.

Simplest possible solution for communication between services is to use JSON over HTTP, however there are much more options: gRPC, pub/sub, etc.

Sounds cool, but there are challenges which come with microservices

 - Serialization
 - Logging
 - Circuit breakers
 - Request tracing
 - Service discovery

And if you are a Go developer, here go-kit comes to us with set of abstractions, packages and interfaces for the developer, so the implementations across your services become standard. 

With this video I want to start an in-depth tutorial on using go-kit tool. We'll create a system built on microservices, setup environment, review how services interact with each other.

We will create a fictional bug tracker system with help of few microservices:

 - Users
 - Bugs
 - Notificator

Some of them will be accessible with JSON over HTTP, and internal communication will be done with gRPC.

### go-kit review

We should understand that go-kit is not a framework, it's a toolkit for building microservices in Go, including packages and interfaces. It is similar to Java Spring Boot but smaller in scope.

Let's init our project.

There is a `kitgen` command line tool to generate a service from template which is not ready to be used yet, but there are other packages which can help you.

### go-kit CLI

2020 Update: "github.com/kujtimiihoxha/kit" is not supported anymore.

There is a separate package to create a service from template:

```
go get github.com/go-kit/kit
go get github.com/kujtimiihoxha/kit
```

Let's create our services:

```
kit new service users
kit new service bugs
kit new service notificator
```

This will generate the initial folder structure and the service interface. The interface is empty by default, let's define the functions in our interface. We need a function for User creation, let's start with this.

users:
```
package service

import "context"

// UsersService describes the service.
type UsersService interface {
	Create(ctx context.Context, email string) error
}
```

bugs:
```
package service

import "context"

// BugsService describes the service.
type BugsService interface {
	// Add your methods here
	Create(ctx context.Context, bug string) error
}
```

notifcator:
```
package service

import "context"

// NotificatorService describes the service.
type NotificatorService interface {
	// Add your methods here
	SendEmail(ctx context.Context, email string, content string) error
}

```

Then we need to run a command to generate a service, it will create the service boilerplate, service middleware and endpoint code. It also creates `cmd/` package to run our service.

```
kit generate service users --dmw
kit generate service bugs --dmw
```

--dmw creates default endpoint middleware, logging middleware.

This command has added go-kit packages to our code already: endpoint and http transport. What we need to do now is to implement our business logic in 1 place only.

We will continue with business logic in the next video.

Notificator should not have REST API as it's an internal service, so we generate service with gRPC transport. gRPC stands for Google RPC framework, if you never used it before, check https://grpc.io.

For this we need to [install protoc and protobuf](https://developers.google.com/protocol-buffers/docs/gotutorial) first.

```
kit generate service notificator -t grpc --dmw
```

This also created .pb file, but we will fill it in the next video.

go-kit CLI can also create a boilerplate docker-compose setup, let's try it.

```
kit generate docker
```

So it created Dockerfile, docker-compose.yml with ports mapping. Let's run our environment and trigger our `/create` endpoint.

```
docker-compose up
```

Dockerfiles are using `watcher` go package, which is updating and restarting binary files if Go code has been changed, which is very convenient on local environment.

Now our services are running on the ports 8800, 8801, 8802. Let's call the endpoint of Users:

```
curl -XPOST http://localhost:8800/create -d '{"email": "test"}'
```

### Conclusion

We haven't implemented services yet, but we prepared a good local environment in few mins, which can be deployed later to your infrastructure as we containerized it with Docker.

In the next video we'll connect these services between each other and implement some logic and we'll try some go-kit packages to manage microservices.