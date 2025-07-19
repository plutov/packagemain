package main

import "github.com/gonutz/prototype/draw"

func getTextColor(neighbors int) draw.Color {
	switch neighbors {
	case 1:
		return draw.LightBlue
	case 2:
		return draw.Green
	case 3:
		return draw.LightRed
	default:
		return draw.White
	}
}
