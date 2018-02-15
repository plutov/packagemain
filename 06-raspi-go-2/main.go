package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/blackjack/webcam"
	"github.com/machinebox/sdk-go/facebox"
	htgotts "github.com/plutov/htgo-tts"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
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
	speech := htgotts.Speech{Folder: "audio", Language: "en"}

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

			log.Printf("found %d faces", len(faces))

			for _, f := range faces {
				greeted := isGreeted(f.Name)

				log.Printf("face: %s, confidence: %.2f, greeted: %t", f.Name, f.Confidence, greeted)
				if !greeted && f.Confidence >= 0.5 {
					greetings[f.Name] = time.Now()

					err = speech.Speak(fmt.Sprintf("Hello, is that you %s?", f.Name))
					if err != nil {
						log.Printf("unable to run text-to-speech: %v", err)
						continue
					}

					file := fmt.Sprintf("record/%d.wav", time.Now().UnixNano())
					log.Printf("recording voice input into %s", file)
					err = record(file, 3)
					if err != nil {
						log.Printf("unable to record user voice input: %v", err)
						continue
					}

					log.Printf("recording completed: %s", file)

					text, err := speechToText(file)
					if err != nil {
						log.Printf("unable to run speech-to-text: %v", err)
						continue
					}

					log.Printf("speech to text results: %s", text)

					if strings.Contains(text, "yes") {
						speech.Speak(fmt.Sprintf("Have a nice day %s", f.Name))
					}

					if strings.Contains(text, "no") {
						speech.Speak("Oops, I am sorry")
					}

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
	return ok && time.Now().Before(g.Add(time.Hour*12))
}

func record(fileName string, timeLimitSecs int) (err error) {
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "trim", "0", strconv.Itoa(timeLimitSecs))

	env := os.Environ()
	env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env

	return cmd.Run()
}

func speechToText(filename string) (string, error) {
	ctx := context.Background()

	client, err := speech.NewClient(ctx, option.WithServiceAccountFile("speech.json"))
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return "", err
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			return alt.Transcript, err
		}
	}

	return "", err
}
