package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	_ "image/jpeg"
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
			Position:     image.Point{0, 0},
			NextLayer: &Layer{
				AssetsFolder: "./quotes",
				Position:     image.Point{0, 0},
			},
		},
	}

	baseImage := image.NewRGBA(image.Rect(0, 0, width, height))
	generatedImages, err := addLayer([]image.Image{baseImage}, layers)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(len(generatedImages))
}

func addLayer(prevImages []image.Image, layer *Layer) ([]image.Image, error) {
	if layer == nil {
		return prevImages, nil
	}

	layerImages := []image.Image{}

	// get all images from layer folder
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
		return []image.Image{}, err
	}

	newImages := []image.Image{}
	for range prevImages {
		for _, layerImage := range layerImages {
			newImages = append(newImages, layerImage)
		}
	}

	return addLayer(newImages, layer.NextLayer)
}
