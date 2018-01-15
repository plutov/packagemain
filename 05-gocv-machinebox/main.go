package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/machinebox/sdk-go/facebox"
	"gocv.io/x/gocv"
)

var (
	blue          = color.RGBA{0, 0, 255, 0}
	faceAlgorithm = "haarcascade_frontalface_default.xml"
	fbox          = facebox.New("http://localhost:8080")
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

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	classifier.Load(faceAlgorithm)
	defer classifier.Close()

	for {
		if ok := webcam.Read(img); !ok {
			log.Print("cannot read webcam")
			continue
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		for _, r := range rects {
			// Save each found face into the file
			imgFace := img.Region(r)
			imgName := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
			gocv.IMWrite(imgName, imgFace)
			imgFace.Close()

			f, err := os.Open(imgName)
			if err != nil {
				log.Printf("unable to open saved img: %v", err)
				continue
			}

			faces, err := fbox.Check(f)
			f.Close()
			if err != nil {
				log.Printf("unable to recognize face: %v", err)
			}

			var caption = "I don't know you"
			if len(faces) > 0 {
				caption = fmt.Sprintf("I know you %s", faces[0].Name)
			}

			// draw rectangle for the face
			size := gocv.GetTextSize(caption, gocv.FontHersheyPlain, 3, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(img, caption, pt, gocv.FontHersheyPlain, 3, blue, 2)
			gocv.Rectangle(img, r, blue, 3)
		}

		// show the image in the window, and wait 100ms
		window.IMShow(img)
		window.WaitKey(100)
	}
}
