package main

import rl "github.com/gen2brain/raylib-go/raylib"

func getTextColor(neighbors int) rl.Color {
	switch neighbors {
	case 1:
		return rl.Blue
	case 2:
		return rl.Green
	case 3:
		return rl.Red
	default:
		return rl.Black
	}
}
