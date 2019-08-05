## Rate Limiting HTTP Requests in Go based on IP address

If you are running HTTP server and want to rate limit requests to the endpoints, you can use well-maintained tools such as [github.com/didip/tollbooth](https://github.com/didip/tollbooth). But if you're building somethinig very simple, it's not that hard to implement it on your own.

There is already experimental Go package `x/time/rate`, which we can use.

In this tutorial we'll create a simple middleware for rate limiting based on the user's IP address.

### Pure HTTP Server

Let's start with building a simple HTTP server, that has very simple endpoint. It could be a heavy endpoint, that's why we want to add a rate limit there.

### Rate Limit middleware

The limiter permits you to consume an average of r tokens per second, with a maximum of b tokens in any single 'burst'. So in the code above our limiter allows 2 tokens to be consumed per second, with a maximum burst size of 5.

In the `limitMiddleware` function we call the global limiter's Allow() method each time the middleware receives a HTTP request. If there are no tokens left in the bucket Allow() will return false and we send the user a 429 Too Many Requests response. Otherwise, calling Allow() will consume exactly one token from the bucket and we pass on control to the next handler in the chain.

### Build & Run

```
go get golang.org/x/time/rate
go build .
./17-ratelimit
```

### Test

There is one very nice tool I like to use for HTTP load testing, called [vegeta](https://github.com/tsenart/vegeta), which is also written in Go.

```
brew install vegeta
```

We need to create a simple config file saying what requests do we want to produce.

#### vegeta.conf

```
GET http://localhost:8888/
```

And then run attack for 10 seconds with 100 requests per time unit.

```
vegeta attack -duration=10s -rate=100 -targets=vegeta.conf | vegeta report
```