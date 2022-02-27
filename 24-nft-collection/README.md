## packagemain #24: Generate an NFT Collection in Go

Nowadays almost any conversation will inevitably end up with discussing NFTs or non-fungible tokens. However, I'm not going to talk about it much today, as I am not an expert. I will rather focus on a technical part of generating a collection of unique images by blending the layers in Go. That should be fun. And probably later if I gain more knowledge in this space we can "mint" them as well.

### How the images will be generated?

To generate the images for a collection, we'll be using multiple types of layers and blend them programmatically. For this example I decided to use the following 3 layers:
- 10 background images
- 10 gophers
- 10 quotes that gophers say

This will allow us to create 1000 unique images for our collection. Each image will contain only 1 background, 1 gopher and 1 quote.

I am not an artist, so I had to find some assets to do this video. I want to respect the authors' licenses, so I carefully selected few assets that we can legally use here.

- Backgrounds: https://unsplash.com/@cocoloris
- Gophers: https://github.com/MariaLetta/free-gophers-pack
- Quotes: I generated myself 10 simple quotes

I already uploaded all the images that we'll use into the separate folders.

### What the program will do

Before writing the program, let's define what it'll do. We will build an executable program in which we'll define where to find layers' images. The program then will blend all unique combinations programmatically, and store them into a separate folder.

As you can already imagine, this program will rely a lot on [image](https://pkg.go.dev/image) package.

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

Our final image will be 1024x1024 as the size of my background images, and the `Position` you see above was approximately colculated by me, so the Gopher goes into the middle of the final image and quote goes to the top-right corner.

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

    // add next layer to all current images
	return addLayer(newImages, layer.NextLayer)
}
```

Now we can call the `addLayer` function in our `main` function. To do so we'll create a base image container.

```go
func main() {
    // base image container with defined size
    baseImage := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
    collection, err := addLayer([]*image.RGBA{baseImage}, &backgroundLayer)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // just for the debugging, we should get 1000 final images
    fmt.Println(len(collection))
}
```

### Getting all images from the layout folder

Let's go back to `addLayer` function and implement missing pieces.

```go
// get all images from layer folder
layerImages := []image.Image{}
err := filepath.Walk(layer.AssetsFolder, func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }

    // skip directories
    if info.IsDir() {
        return nil
    }

    file, fileErr := os.Open(path)
    if fileErr != nil {
        return errors.Wrap(fileErr, "unable to open file")
    }

    defer file.Close()

    img, _, decodeErr := image.Decode(file)
    if decodeErr != nil {
        return errors.Wrap(decodeErr, fmt.Sprintf("unable to decode image, path: %s", path))
    }

    layerImages = append(layerImages, img)

    return nil
})
```

### Blending layers

To add image on top of the other image we'll use `image/draw` package. First of all we have to clone the previous image so we don't mutate it, and then add new layer on top.

```go
// clone image into new variable dst
dst := image.NewRGBA(prevImage.Bounds())
draw.Draw(dst, prevImage.Bounds(), prevImage, image.Point{}, draw.Over)

// add new layer
draw.Draw(dst, layerImage.Bounds().Add(layer.Position), layerImage, image.Point{}, draw.Over)

newImages = append(newImages, dst)
```

### Saving final images into a folder

Now our images are ready, but they stay only in memory, let's put them on disc:

```go
for i, img := range collection {
    // use random unique name for images
    out, err := os.Create(fmt.Sprintf("./collection/%s.png", uuid.NewString()))
    if err != nil {
        fmt.Printf("unable to create a file: %s", err.Error())
        os.Exit(1)
    }

    if err := png.Encode(out, img); err != nil {
        fmt.Printf("unable to encode an image: %s", err.Error())
        os.Exit(1)
    }

    out.Close()
}
```

### Results

Looks like everything is ready to run our program.

```
go run main.go
```

It will take few seconds to execute it, since png.Encode can take some time for 1000 files.

Now we can open our `./collection` folder to see what we've got. Some images look pretty funny IMO :)

![output.png](https://raw.githubusercontent.com/plutov/packagemain/master/24-nft-collection/output.png)

### What's next?

All right, we learned a bit how to generate images from layers in Go. What can we do with that? Like I said, if the audience is insterested about NFT topic in Go and I become more familiar with it as well, we can use this collection as our sample for "minting".
