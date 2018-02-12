package main

import (
	"bytes"
	"log"
	"time"

	"github.com/blackjack/webcam"
	"github.com/machinebox/sdk-go/facebox"
)

func main() {
	cam, err := openWebcam("/dev/video0")
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

	fbox := facebox.New("http://192.168.1.216:8080")
	for {
		frame, err := cam.ReadFrame()
		if err != nil {
			log.Printf("unable to read frame: %v", err)
			continue
		}

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
	}
}

// On Raspberry Pi we can often see error open /dev/video0: device or resource busy
// This function will try to open camera again each second 10 times max
func openWebcam(path string) (*webcam.Webcam, error) {
	attemptsLeft := 10
	var (
		cam *webcam.Webcam
		err error
	)

	for attemptsLeft > 0 {
		cam, err = webcam.Open(path)
		if err == nil {
			return cam, nil
		}

		time.Sleep(time.Second)
		attemptsLeft--
	}

	return cam, err
}
