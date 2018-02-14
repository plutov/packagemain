package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/blackjack/webcam"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/machinebox/sdk-go/facebox"
)

var (
	greetings = make(map[string]time.Time)
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
				if f.Confidence >= 0.5 {
					greeted := isGreeted(f.Name)
					log.Printf("face: %s, confidence: %.2f, greeted: %t", f.Name, f.Confidence, greeted)
					if !greeted {
						greetings[f.Name] = time.Now()
						speech := htgotts.Speech{Folder: "audio", Language: "en"}
						speech.Speak(fmt.Sprintf("Hi %s, how are you today?", f.Name))
						break
					}
				}
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

func isGreeted(name string) bool {
	g, ok := greetings[name]
	now := time.Now()
	return ok && now.Before(g.Time.Add(time.Hour*12))
}

// Record from mic to a file
// Stop if silence detected (<3% volume for 0.5s)
// timeLimitSecs is a maximum time
// rate 16000, bit 16
func record(fileName string, timeLimitSecs int) (err error) {
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "trim", "0", strconv.Itoa(timeLimitSecs), "silence", "1", "0.5", "3%", "1", "0.5", "3%")

	env := os.Environ()
	env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	return cmd.Run()
}
