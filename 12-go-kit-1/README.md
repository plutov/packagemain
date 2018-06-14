## Microservices with go-kit. Part 1

Nowadays, Microservices is one of the most popular buzz-word in the field of software architecture. Applications designed using microservices consist of a set of several small services cooperating and communicating together. Separation between services is enforced by the service's external API. Each individual micro service can be scaled and deployed separately from the rest.

But where monolithic application have one source of logs, one source of metrics, one application to deploy, one API to rate limit, etc, microservice based application have multiple sources. Some of the common concerns of application design that are amplified in a microservices based application are:

 - Rate limiters
 - Serialization
 - Logging
 - Circuit breakers
 - Request tracing
 - Service discovery

And here go-kit comes to us with set of abstractions, packages and interfaces for the developer, so the implementations across your services becomes standard. 

With this video I will start in-depth tutorial on using go-kit. We'll create entire system combined of microservices, setup environment, review how it interacts.

### go-kit review

But first of all let's go to GitHub and review the go-kit project.

https://github.com/go-kit/kit

As you can see there are a lot of folders: sd, auth, circuit breaker, etc.