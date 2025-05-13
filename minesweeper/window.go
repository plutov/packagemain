package main

import rl "github.com/gen2brain/raylib-go/raylib"

func centerWindow(width, height int) {
	rl.SetWindowSize(width, height)
	monitorWidth := rl.GetMonitorWidth(0)
	monitorHeight := rl.GetMonitorHeight(0)
	rl.SetWindowPosition((monitorWidth-width)/2, (monitorHeight-height)/2)
}
