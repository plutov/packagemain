### Building Face Detection Go program for Raspberry Pi. Part 1

I always wanted to try running Go programs on Raspberry Pi, as it allows you to run something interesting on a small device. I got the Raspberry Pi 3 and I will do few videos about writing and running Go programs there. I think it's better to cover it in a few videos, because we will work with different things here: writing and building program itself, capture images from web camera, using Google Speech API there and maybe more.

Sounds cool, let's start!

In this video we will write a Go program which can capture image from web camera and send it to MachineBox to detect the face and tell us the name.

My Raspberry Pi is already connected to wi-fi and has SSH interface. Also I have webcam connected to my Raspberry Pi.

Let's write a small program to capture image from web camera. In previous video I did it with OpenCV, I found a package blackjack/webcam, which uses V4L2 Linux framework, and it's already installed on Raspberry Pi. It makes life easier, because OpenCV installation on RasPi may be complicated.

But this package works only in Linux, so when we do go get we need to set GOARCH and GOOS environment variables.

```
sudo GOARCH=arm GOOS=linux go get github.com/blackjack/webcam
```

```
package main

import (
	"log"

	"github.com/blackjack/webcam"
)

func main() {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		log.Fatalf("unable to open webcam: %v", err)
	}
	defer cam.Close()

	for code, formatName := range cam.GetSupportedFormats() {
		if formatName == "Motion-JPEG" {
			cam.SetImageFormat(code, 1280, 720)
		}
	}

	err = cam.StartStreaming()
	if err != nil {
		log.Fatalf("unable to start streaming: %v", err)
	}

	for {
		cam.WaitForFrame(500000)

		frame, err := cam.ReadFrame()
		if err != nil {
			log.Printf("unable to read frame: %v", err)
			continue
		}

		if len(frame) != 0 {
			log.Println("ok")
		}
	}
}
```

And when we build we also need to set GOARCH and GOOS:

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
```

In previous video I was showing how to run Facebox container to detect faces via Docker, and in this video I was planning to run Facebox on Raspberry Pi, but currently it's not possible to run it on ARM architecture. But that's fine, as anyway it's better to run it on server with faster CPU. So our Raspberry Pi device will work as client, and I will skip the Facebox installation step, I already have it running.

I just want to check that it's accessible from Raspberry Pi:
```
curl http://192.168.1.216:8080/info
```

Cool, now we can send frame to Facebox for recognition. But first of all let's get facebox Go package:
```
go get github.com/machinebox/sdk-go/facebox
```

```
fbox := facebox.New("http://192.168.1.216:8080")
```

I found that mjpeg stream frames do not have Huffman Table information (DHT), so we need to add DHT segment manually, so it will be a standard JPEG file. In file `dht.go` I have a small function to change a frame slice of bytes.

```
if len(frame) != 0 {
	frame = addMotionDht(frame)

	faces, err := fbox.Check(bytes.NewBuffer(frame))
	if err != nil {
		log.Printf("unable to recognize face: %v", err)
		continue
	}

	for _, f := range faces {
		log.Printf("face: %s, confidence: %.2f", f.Name, f.Confidence)
	}
}
```

Cool, now we can capture faces, and in the next video we'll add Google Text to Speech integration to our Raspberry Pi client, so it can say something if face is recognized or not.

I hope it was helpful and interesting, see you later!