package main

import (
	"crypto/md5"
	"encoding/binary"
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

type Operation interface {
	InputsCount() uint8
	// always returns 3 elements for RGB
	Eval(x float64, y float64) []float64
}

type OpVarX struct{}

func (o *OpVarX) InputsCount() uint8 {
	return 0
}
func (o *OpVarX) Eval(x float64, y float64) []float64 {
	return []float64{x, x, x}
}

type OpVarY struct{}

func (o *OpVarY) InputsCount() uint8 {
	return 0
}
func (o *OpVarY) Eval(x float64, y float64) []float64 {
	return []float64{y, y, y}
}

// operations without inputs
var leaves = []Operation{&OpVarX{}, &OpVarY{}}

// operations with inputs
var ops = []Operation{}

func main() {
	var phrase string
	// same phrase will always result in the same image
	flag.StringVar(&phrase, "phrase", "", "phrase")
	flag.Parse()

	h := md5.New()
	h.Write([]byte(phrase))
	// this seed is not ideal, as it takes only the first 8 bytes. check out Mersenne-Twister algorithm
	seed := binary.BigEndian.Uint64(h.Sum(nil))
	r := rand.New(rand.NewSource(int64(seed)))

	width := 100
	height := 100

	upLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, bottomRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Initial values from 0 to 1
			fx := float64(x) / float64(width)
			fy := float64(y) / float64(height)

			i := r.Intn(len(leaves) - 1)
			rgba := leaves[i].Eval(fx, fy)
			r := uint8(rgba[0] * 0xff)
			g := uint8(rgba[1] * 0xff)
			b := uint8(rgba[2] * 0xff)
			img.Set(x, y, color.RGBA{r, g, b, 0xff})
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
