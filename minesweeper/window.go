package main

import rl "github.com/gen2brain/raylib-go/raylib"

func centerWindow(width, height int32) {
	rl.SetWindowSize(int(width), int(height))
	monitorWidth := rl.GetMonitorWidth(0)
	monitorHeight := rl.GetMonitorHeight(0)
	rl.SetWindowPosition((monitorWidth-int(width))/2, (monitorHeight-int(height))/2)
}
