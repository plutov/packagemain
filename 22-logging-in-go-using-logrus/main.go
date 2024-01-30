package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}

func main() {
	e := echo.New()
	e.Use(loggingMiddleware)

	e.GET("/isInt", func(c echo.Context) error {
		a := c.QueryParam("a")

		logCtx := log.WithField("a", a)
		logCtx.Debug("parsing the string")

		_, err := strconv.Atoi(a)
		if err != nil {
			logCtx.Error("unable to parse a string")
			return c.String(http.StatusBadRequest, "")
		}
		logCtx.Debug("string is parsed")

		return c.String(http.StatusOK, "ok")
	})

	e.Start(":5050")
}

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
