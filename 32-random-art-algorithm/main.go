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
	// a number of leaves which produce inputs for this operation
	InputsCount() uint8
	// set children of the operation
	SetInputs([]Operation)
	// always returns 3 elements for RGB
	Eval(x float64, y float64) []float64
}

type OpVarX struct{}
type OpVarY struct{}
type OpConstant struct {
	constant float64
}
type OpColorMix struct {
	inputs []Operation
}

func (o *OpVarX) InputsCount() uint8 {
	return 0
}
func (o *OpVarX) SetInputs(inputs []Operation) {}
func (o *OpVarX) Eval(x float64, y float64) []float64 {
	return []float64{x, x, x}
}

func (o *OpVarY) InputsCount() uint8 {
	return 0
}
func (o *OpVarY) SetInputs(inputs []Operation) {}
func (o *OpVarY) Eval(x float64, y float64) []float64 {
	return []float64{y, y, y}
}

func (o *OpConstant) InputsCount() uint8 {
	return 0
}
func (o *OpConstant) SetInputs(inputs []Operation) {}
func (o *OpConstant) Eval(x float64, y float64) []float64 {
	return []float64{o.constant, o.constant, o.constant}
}

func (o *OpColorMix) InputsCount() uint8 {
	return 3
}
func (o *OpColorMix) SetInputs(inputs []Operation) {
	o.inputs = inputs
}
func (o *OpColorMix) Eval(x float64, y float64) []float64 {
	r := o.inputs[0].Eval(x, y)[0]
	g := o.inputs[1].Eval(x, y)[1]
	b := o.inputs[2].Eval(x, y)[2]

	return []float64{r, g, b}
}

func main() {
	var (
		phrase string
		depth  int
	)
	// same phrase will always result in the same image
	flag.StringVar(&phrase, "phrase", "", "phrase")
	flag.IntVar(&depth, "depth", 3, "depth of graph")
	flag.Parse()

	h := md5.New()
	h.Write([]byte(phrase))
	// this seed is not ideal, as it takes only the first 8 bytes. check out Mersenne-Twister algorithm
	seed := binary.BigEndian.Uint64(h.Sum(nil))
	r := rand.New(rand.NewSource(int64(seed)))

	// operations without inputs
	opsNoLeaves := []Operation{&OpVarX{}, &OpVarY{}, &OpConstant{r.Float64()}, &OpColorMix{}}
	// operations with inputs
	opsWithLeaves := []Operation{&OpColorMix{}}

	root := &OpColorMix{}
	root.SetInputs(generateNodeChildren(root.InputsCount(), depth))
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

// Given an expected children count and current depth generate the list of inputs
func generateNodeChildren(count uint8, depth int) []Operation {

}
