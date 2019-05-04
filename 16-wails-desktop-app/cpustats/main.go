package main

import (
	"github.com/leaanthony/mewn"
	"github.com/plutov/packagemain/16-wails-desktop-app/cpustats/pkg/sys"
	"github.com/wailsapp/wails"
)

func main() {
	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	stats := &sys.Stats{}

	app := wails.CreateApp(&wails.AppConfig{
		Width:  512,
		Height: 512,
		Title:  "CPU Usage",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(stats)
	app.Run()
}
