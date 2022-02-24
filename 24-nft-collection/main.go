package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/pkg/errors"
)

type Layer struct {
	AssetsFolder string
	Position     image.Point
	NextLayer    *Layer
}

func main() {
	width := 1024
	height := 1024

	layers := &Layer{
		AssetsFolder: "./backgrounds",
		Position:     image.Point{0, 0},
		NextLayer: &Layer{
			AssetsFolder: "./gophers",
			Position:     image.Point{256, 256},
			NextLayer: &Layer{
				AssetsFolder: "./quotes",
				Position:     image.Point{668, 100},
			},
		},
	}

	// base image container with defined size
	baseImage := image.NewRGBA(image.Rect(0, 0, width, height))
	collection, err := addLayer([]*image.RGBA{baseImage}, layers)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for i, img := range collection {
		out, err := os.Create(fmt.Sprintf("./collection/%d.png", i+1))
		if err != nil {
			fmt.Printf("unable to create a file: %s", err.Error())
			os.Exit(1)
		}

		if err := png.Encode(out, img); err != nil {
			fmt.Printf("unable to encode an image: %s", err.Error())
			os.Exit(1)
		}
	}
}

func addLayer(prevImages []*image.RGBA, layer *Layer) ([]*image.RGBA, error) {
	if layer == nil {
		return prevImages, nil
	}

	// get all images from layer folder
	layerImages := []image.Image{}
	err := filepath.Walk(layer.AssetsFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "unable to open file")
		}

		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to decode image, path: %s", path))
		}

		layerImages = append(layerImages, img)

		return nil
	})

	if err != nil {
		return []*image.RGBA{}, err
	}

	newImages := []*image.RGBA{}
	for _, prevImage := range prevImages {
		for _, layerImage := range layerImages {
			// clone image into new variable dst
			dst := image.NewRGBA(prevImage.Bounds())
			draw.Draw(dst, prevImage.Bounds(), prevImage, image.Point{}, draw.Over)

			// add new layer
			draw.Draw(dst, layerImage.Bounds().Add(layer.Position), layerImage, image.Point{}, draw.Over)

			newImages = append(newImages, dst)
		}
	}

	return addLayer(newImages, layer.NextLayer)
}
