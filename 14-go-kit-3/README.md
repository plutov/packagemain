# Microservices with go-kit. Part 3

In the previous video we have implemented fictional Notificator gRPC service, registered it in etcd and invoked from another service Users.

In this video we'll continue to work with the code we already have.

I will copy all code to a new folder and replace the go imports.

```bash
cp -R ../13-go-kit-2/* .
find . -type f -print0 | xargs -0 sed -i "" "s/13-go-kit-2/14-go-kit-3/g"
```

When one service synchronously invokes another there is always the possibility that the other service is unavailable or is exhibiting such high latency it is essentially unusable. Previous resources such as threads might be consumed in the caller while waiting for the other service to respond. This might lead to resource exhaustion, which would make the calling service unable to handle other requests. The failure of one service can potentially cascade to other services throughout the application.

So how to prevent a network or service failure from cascading to other services?

A service client should invoke a remote service via a proxy that functions in a similar fashion to an electrical circuit breaker. When the number of consecutive failures crosses a threshold, the circuit breaker trips, and for the duration of a timeout period all attempts to invoke the remote service will fail immediately. After the timeout expires the circuit breaker allows a limited number of test requests to pass through. If those requests succeed the circuit breaker resumes normal operation. Otherwise, if there is a failure the timeout period begins again.

And go-kit has a package that implements the circuit breaker pattern. It provides several implementations in this package, but if you're looking for guidance, Gobreaker is probably the best place to start. It has a simple and intuitive API, and is well-tested.

Let's implement it in our fictional system. As you remember Users service calls Notificator, and we have to wrap calls to Notificator with circuit breaker pattern.