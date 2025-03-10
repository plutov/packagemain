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

func main() {
	var output string
	var phrase string
	// higher depth requires more time but produces better results
	var depth int
	// same phrase will always result in the same image
	flag.StringVar(&phrase, "phrase", "", "phrase")
	flag.IntVar(&depth, "depth", 5, "depth of graph")
	flag.StringVar(&output, "out", "image.png", "output file")
	flag.Parse()

	// seeded pseudo-random number generator
	prng := getPRNG(phrase)

	// graph root
	root := &OpColorMix{}

	// populate the graph
	root.SetInputs(generateGraphNodes(root.InputsCount(), depth, prng))

	width := 600
	height := 600

	upLeft := image.Point{0, 0}
	bottomRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, bottomRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Initial x, y values from 0 to 1
			xInit := float64(x) / float64(width)
			yInit := float64(y) / float64(height)

			// apply the graph to each pixel
			rgb := root.Eval(xInit, yInit)
			r := uint8(rgb[0]*0xff) % 0xff
			g := uint8(rgb[1]*0xff) % 0xff
			b := uint8(rgb[2]*0xff) % 0xff
			img.Set(x, y, color.RGBA{r, g, b, 0xff})
		}
	}

	// Encode as PNG.
	f, _ := os.Create(output)
	png.Encode(f, img)
}

// Pseudo Random Number Generator with md5 phrase as a seed
func getPRNG(phrase string) *rand.Rand {
	h := md5.New()
	h.Write([]byte(phrase))
	// this seed is not ideal, as it takes only the first 8 bytes. check out Mersenne-Twister algorithm
	seed := binary.BigEndian.Uint64(h.Sum(nil))
	return rand.New(rand.NewSource(int64(seed)))
}

// Given an expected children count and current depth generate the list of children nodes
// Recursive
func generateGraphNodes(count uint8, depth int, prng *rand.Rand) []Operation {
	if count == 0 {
		return []Operation{}
	}

	ops := []Operation{}
	for i := 0; i < int(count); i++ {
		op := pickOperation(prng, depth)
		op.SetInputs(generateGraphNodes(op.InputsCount(), depth-1, prng))
		ops = append(ops, op)
	}

	return ops
}
