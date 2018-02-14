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
	"github.com/machinebox/sdk-go/facebox"
	htgotts "github.com/plutov/htgo-tts"
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

	//fbox := facebox.New("http://192.168.1.216:8080")
	for {
		frame, err := cam.ReadFrame()
		if err != nil {
			log.Printf("unable to read frame: %v", err)
			continue
		}

		if len(frame) != 0 {
			frame = addMotionDht(frame)

			faces := []facebox.Face{facebox.Face{
				Name:       "Alex",
				Confidence: 0.5,
			}}
			//faces, err := fbox.Check(bytes.NewBuffer(frame))
			if err != nil {
				log.Printf("unable to recognize face: %v", err)
				continue
			}

			log.Printf("found %d faces", len(faces))

			for _, f := range faces {
				greeted := isGreeted(f.Name)

				log.Printf("face: %s, confidence: %.2f, greeted: %t", f.Name, f.Confidence, greeted)
				if !greeted && f.Confidence >= 0.5 {
					greetings[f.Name] = time.Now()

					msg := fmt.Sprintf("Hello, is that you %s?", f.Name)
					log.Printf("greeting user: %s", msg)

					speech := htgotts.Speech{Folder: "audio", Language: "en"}
					err = speech.Speak(msg)
					if err != nil {
						log.Printf("unable to run text-to-speech: %v", err)
						continue
					}

					file := fmt.Sprintf("record/%d.wav", time.Now().UnixNano())
					log.Printf("recording voice input into %s", file)
					err = record(file, 5)
					if err != nil {
						log.Printf("unable to record user voice input: %v", err)
						continue
					}

					log.Printf("recording completed: %s", file)

					break
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
	return ok && now.Before(g.Add(time.Hour*12))
}

// Record from mic to a file
// Stop if silence detected (<3% volume for 0.5s)
// timeLimitSecs is a maximum time
// rate 16000, bit 16
func record(fileName string, timeLimitSecs int) (err error) {
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "trim", "0", strconv.Itoa(timeLimitSecs))

	env := os.Environ()
	env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	return cmd.Run()
}
