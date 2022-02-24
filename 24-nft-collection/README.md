## packagemain #24: Generate an NFT Collection in Go

Nowadays almost any conversation will inevitably end up with discussing NFTs or non-fungible tokens. However, I'm not going to talk about it much today, as I am not an expert. I will rather focus on a technical part of generating a collection of unique images by blending the layers in Go. That should be fun. And probably later if I gain more knowledge in this space we can "mint" them as well.

### How the images will be generated?

To generate the images for a collection, we'll be using mutiple types of layers and blend them programmatically. For this example I decided to use the following 3 layers:
- 10 background images
- 10 gophers
- 10 quotes that gophers say

This will allow us to create 1000 unique images for our collection. Each emage will contain only 1 background, 1 gopher and 1 quote.

I am not an artist, so I had to find some assets to do this video. I want to expect the authors' licenses, so I carefully selected few assets that we can legally use here.

- Backgrounds: https://unsplash.com/@cocoloris (license: https://unsplash.com/license)
- Gophers: https://github.com/MariaLetta/free-gophers-pack (license: https://github.com/MariaLetta/free-gophers-pack/blob/master/LICENSE)
- Quotes: I generated myself 10 simple quotes

I already uploaded all images that we'll use into the separate folders.

### What the program will do

Before writing the program, let's define what it'll do. We will build an executable program in which we'll define where to find layers' images. The prohram then will blend all unique combinations programmatically, and store them into a separate folder.

This program will rely a lot on [image](https://pkg.go.dev/image) package.

### Define the Layer type

Now, as everything is ready, let's define some basic types and functions. For simplicity, I'll code the whole program in a single `main.go` file, but some functions/types could be placed in corresponding packages if you'd like to.

`Layer` type:

```go
type Layer struct {
	AssetsFolder string
	Position     image.Point
	NextLayer    *Layer
}
```

Each layer needs some configuration, like a path to the folder where all layer's images are stored, an offset telling where to put an item on the final image, and a pointer to the next layer. Why pointer though? I think this program will benefit from the recursive approach, so linked list data structure would work well here.

Now let's define our layers by creating the instances of our new type:

```go
quoteLayer := Layer{
    AssetsFolder: "./quotes",
    Position:     image.Point{668, 100},
}
gopherLayer := Layer{
    AssetsFolder: "./gophers",
    Position:     image.Point{256, 256},
    NextLayer:    &quoteLayer,
}
backgroundLayer := Layer{
    AssetsFolder: "./backgrounds",
    Position:     image.Point{0, 0},
    NextLayer:    &gopherLayer,
}
```

Our final image will be 1024x1024 as the size of my background images, and the `Position` you see above was approximately colculated by, so the Gopher goes into the middle of the final image and quote goes to the top-right corner.

### Write a recursive function to iterate over layers and images

Now we can create a recursive function that will go layer by layer iterating over all images in each folder and creating all combinations.

The function will receive previously generated images and for each previous image it will add a new layer. It will return all new images.

```go
func addLayer(prevImages []*image.RGBA, layer *Layer) ([]*image.RGBA, error) {
    if layer == nil {
		return prevImages, nil
	}

    // get all images from layer folder
	layerImages := []image.Image{}
    // TODO: traverse the layer folder

    // for each previous image create all combinations of current layer
	newImages := []*image.RGBA{}
	for _, prevImage := range prevImages {
		for _, layerImage := range layerImages {
			// clone image into new variable dst
			dst := image.NewRGBA(prevImage.Bounds())

            // TODO: blend layers

			newImages = append(newImages, dst)
		}
	}

	return addLayer(newImages, layer.NextLayer)
}
```