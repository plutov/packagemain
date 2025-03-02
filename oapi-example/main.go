package main

import (
	"embed"
	"net/http"
	"oapiexample/pkg/api"

	"github.com/labstack/echo/v4"
)

//go:embed pkg/api/index.html
//go:embed openapi.yaml
var swaggerUI embed.FS

func main() {
	server := api.NewServer()

	e := echo.New()

	api.RegisterHandlers(e, api.NewStrictHandler(
		server,
		// add middlewares here if needed
		[]api.StrictMiddlewareFunc{},
	))

	// serve swagger docs
	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerUI)))))

	e.Start("127.0.0.1:8080")
}
