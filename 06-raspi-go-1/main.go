package main

import (
	"bytes"
	"log"

	"github.com/blackjack/webcam"
	"github.com/machinebox/sdk-go/facebox"
)

func main() {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		log.Fatalf("unable to open webcam: %v", err)
	}
	defer cam.Close()

	err = cam.StartStreaming()
	if err != nil {
		log.Fatalf("unable to start streaming: %v", err)
	}

	fbox := facebox.New("http://192.168.1.216:8080")
	for {
		cam.WaitForFrame(500000)

		frame, err := cam.ReadFrame()
		if err != nil {
			log.Printf("unable to read frame: %v", err)
			continue
		}

		if len(frame) != 0 {
			faces, err := fbox.Check(bytes.NewReader(frame))
			if err != nil {
				log.Printf("unable to recognize face: %v", err)
				continue
			}

			for _, f := range faces {
				log.Printf("face: %s, confidence: %.2f", f.Name, f.Confidence)
			}
		}
	}
}
