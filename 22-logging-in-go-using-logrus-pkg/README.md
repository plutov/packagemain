## packagemain #22: Logging in Go using logrus

Logging is a very essential part of large software, it's hard to overstate the importance of logging, be it performance metrics logging, error logging, or debug logging for troubleshooting later.

Go standard library has a `log` package, which can print messages, can panic, but feels very limited when working on enterprise-level software when you need better control over formatting, structure and distribution. That's why a lot of third-party logging packages were born, for example logrus, oklog, zerolog, etc. They have a lot of similarities so we're not going to review all of them, but can take `logrus` as an example, which I'm currently using a lot. Usually these packages are fully compatible with built-in `log` package, so replacing the log package in your project should not be a problem.

Let's talk a bit about what can we log. Actually, you can log so many things, so sometimes it's easier to ask what not to log :)

Some things to log:
- Errors/Warnings, whenever something serious happened that developers need to know about
- Debug messages with all contextual data, to be used later while troubleshooting
- Performance metrics, request latencies, memory consumptions

What not to log:
- In general you shouldn't log PII details, such as email addresses, credit card numbers, etc. But it may depend on your project.

Let's now write a simple API service, in which we will call some functions and see what can we log.

```go
func main() {
	e := echo.New()

	e.GET("/isInt", func(c echo.Context) error {
		a := c.QueryParam("a")

		_, err := strconv.Atoi(a)
		if err != nil {
			return c.String(http.StatusBadRequest, "")
		}

		return c.String(http.StatusOK, "ok")
	})

	e.Start(":5050")
}
```

It's very important to understand that logs can be sent not only to `stdout`, but also to files, `stderr` or somewhere else. I'll be using `stdout` only, assuming that later some logs aggregation software will pick up logs from there (for example logstash, fluentd, etc.).

Using the logrus we can set where we want to send the logs output.

```go
log.SetOutput(os.Stdout)
```

Now let's log all the requests, together with status code, latency, etc. We will do it through the Echo middleware.

```go
func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		res := next(c)

		log.WithFields(log.Fields{
			"method":     c.Request().Method,
			"path":       c.Path(),
			"status":     c.Response().Status,
			"latency_ns": time.Since(start).Nanoseconds(),
		}).Info("request details")

		return res
	}
}
```

```
INFO[0001] request details latency_ns=29370 method=GET path=/isInt status=200
INFO[0002] request details latency_ns=10697 method=GET path=/isInt status=400
```

Sometimes your log aggregation tools require logs to be in JSON format which is easier for parsing, for example when using Splunk. I feel nowadays it's generally a default practice to log in JSON.

```go
log.SetFormatter(&log.JSONFormatter{})
```

```
{"latency_ns":25315,"level":"info","method":"GET","msg":"request details","path":"/isInt","status":200,"time":"2021-04-08T12:59:48+02:00"}
{"latency_ns":7581,"level":"info","method":"GET","msg":"request details","path":"/isInt","status":400,"time":"2021-04-08T12:59:59+02:00"}
```

So far we logged only Info errors, now let's add some debugging info.

```go
log.WithField("a", a).Debug("parsing the string")
// ...
log.WithField("a", a).Debug("string is parsed")
```

Unfortunately when I run the code now, I don't see these logs in the stdout. That's because the default log level is set to Info, which means the logs of level Info or higher severity are only printed. To change this, we can configure our program to accept LOG_LEVEL environment variable for example, so we can have `LOG_LEVEL=debug` in dev environment, or `LOG_LEVEL=info` in production. I will also move logrus set up code into `init()` function.

```go
func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}
```

```
LOG_LEVEL=debug go run main.go
{"a":"10","level":"debug","msg":"parsing the string","time":"2021-04-08T13:12:02+02:00"}
{"a":"10","level":"debug","msg":"string is parsed","time":"2021-04-08T13:12:02+02:00"}
```

Sometimes multiple log messages use the same context, in our debug case it is `log.WithField("a", a)`, which can be stored separately and passed between functions even.

```go
logCtx := log.WithField("a", a)
logCtx.Debug("parsing the string")
// ...
```

As we discussed previously, logs can be sent to different outputs. There is one more nice feature of logrus, called Hooks, where you can send specific logs to log collectors or aggregators, such as Sentry.io, Stackdriver, etc.

```go
hook := sentryhook.New(nil)
hook.SetAsync(logrus.ErrorLevel)
logrus.AddHook(hook)
```

Which will send all errors to Sentry.io.

There are multiple other helper functions in `logrus`, but you can check them on your own.

Ok, what do we do now with all these gigabytes of logs we collected? Ideally, you use some tools behind it, which can collect and segment your logs, store them with some defined retention policy.

I use logs as well to build some metrics charts, for example using log-based metrics from Google Cloud Logging, where you can define how to parse your log messages for specific fields and display them. For example in our example we logged `latency_ns` which can be parsed by log server and displayed as a chart by applying some aggregation function to it, so we can monitor the avg/max/... latency for your services.

That's it for today, logging world is quite diverse, there are many packages, tools and ideas. Let me know how you log in Go in the comments below.