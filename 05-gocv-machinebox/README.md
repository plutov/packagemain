### Face Detection in Go using OpenCV and MachineBox

Recently I found a very nice developer-friendly project MachineBox, which provides some machine learning tools inside Docker Container, including face detection, natural language understanding and few more. And it has SDK in Go, so in this video we will build a program which will detect my face. We will also use OpenCV to capture video from Web camera, it also has Go bindings.

Sounds interesting? Let's start!

As I said MachineBox can be installed very easy by running Docker container. First of all we need to register on MachineBox.io and get key, then we can set environment variable MB_KEY and run facebox on port 8080:

```
export MB_KEY=""
docker run -d -p 8080:8080 -e "MB_KEY=$MB_KEY" machinebox/facebox
```

Let's go now to localhost:8080 and verify that it's ready, now we can use Go SDK to communicate with facebox from our Go program.

Now to be able to capture video and recognize faces from the web camera we have to install OpenCV. I am using GoCV.io which has Go bindings for OpenCV. Installation may be complicated, but it worked on my Mac with OpenCV3:

```
brew install opencv3
```

Then we need to set someenvironment variables:

```
source $GOPATH/src/gocv.io/x/gocv/env.sh
```

On https://gocv.io/getting-started/osx/ you can find how to install it on Mac. To verify that GoCV works fine we can run:

```
go run $GOPATH/src/gocv.io/x/gocv/cmd/version/main.go
```

And let's get GoCV package:

```
go get gocv.io/x/gocv
```

Now let's create a main.go file and capture video from web camera:

```
package main

import (
	"log"

	"gocv.io/x/gocv"
)

func main() {
	// open webcam. 0 is the default device ID, change it if your device ID is different
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalf("error opening web cam: %v", err)
	}
	defer webcam.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// open display window
	window := gocv.NewWindow("packagemain")
	defer window.Close()

	for {
		if ok := webcam.Read(img); !ok || img.Empty() {
			log.Print("cannot read webcam")
			continue
		}

		// show the image in the window, and wait 100ms
		window.IMShow(img)
		window.WaitKey(100)
	}
}
```

With OpenCV we can also use classifier to recognize faces on the image. Let's do it, for this we'll use Haar Cascades classifier. I already downloaded this XML file to use, you can find it in GoCV repository. With GoCV it's also possible to draw rectangle for each face to highlight it.

```
var (
	blue          = color.RGBA{0, 0, 255, 0}
	faceAlgorithm = "haarcascade_frontalface_default.xml"
)

...

// load classifier to recognize faces
classifier := gocv.NewCascadeClassifier()
classifier.Load(faceAlgorithm)
defer classifier.Close()

for {
	...

	// detect faces
	rects := classifier.DetectMultiScale(img)
	for _, r := range rects {
		// Save each found face into the file
		imgFace := img.Region(r)
		imgFace.Close()

		// draw rectangle for the face
		size := gocv.GetTextSize("I don't know you", gocv.FontHersheyPlain, 3, 2)
		pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		gocv.PutText(img, "I don't know you", pt, gocv.FontHersheyPlain, 3, blue, 2)
		gocv.Rectangle(img, r, blue, 3)
	}

	...
}
```

To be able to recognize faces we need to train our model, in Facebox it's very easy, we can just upload few images and set the name. I will do it by saving images of my face into a folder.

```
imgName := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
gocv.IMWrite(imgName, imgFace)
```

Now let's go to localhost:8080. There are few options to train Facebox: we can use API, we can use web page or we can use Go SDK.

For this example the easiest way is to use form from web page.

Now after we trained it we can check if Facebox can find my face on the image and we can print my name.

```
fbox          = facebox.New("http://localhost:8080")
```

GoCV has no option to get io.reader of an image, so we will open our saved file.

```
buf, err := gocv.IMEncode("jpg", imgFace)
if err != nil {
	log.Printf("unable to encode matrix: %v", err)
	continue
}

faces, err := fbox.Check(bytes.NewReader(buf))
f.Close()
if err != nil {
	log.Printf("unable to recognize face: %v", err)
}

var caption = "I don't know you"
if len(faces) > 0 {
	caption = fmt.Sprintf("I know you %s", faces[0].Name)
}
```

It's not an advertisement, but I really think that MachineBox is a nice project to provide some Machine learning power to you within simple docker containers.