## Microservices with go-kit. Part 1

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

And if you are a Go developer, here go-kit comes to us with set of abstractions, packages and interfaces for the developer, so the implementations across your services becomes standard. 

With this video I want to start an in-depth tutorial on using go-kit. We'll create a system built on microservices, setup environment, review how services interact with each other.

### go-kit review

We should understand that go-kit is not a framework, it's a toolkit for building microservices in Go, including packages and interfaces. It is similar to Java Spring Boot but smaller in scope.

Let's go to GitHub project page and review the go-kit project.

https://github.com/go-kit/kit

As you can see there are a lot of folders: sd, auth, circuit breaker, etc. which we can import into our project and implement. There is a `kitgen` command line tool to generate a service from template which is not ready to be used yet, but there are other packages which can help you.

### go-kit CLI

There is a separate package to create a service from template:

```
go get github.com/kujtimiihoxha/kit
```

Let's create our first service:

```
kit new service hello
```