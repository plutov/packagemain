### Building Face Detection Go program for Raspberry Pi. Part 2

In previous episode we've build a program in Go which captures image from Webcam and uses Facebox to recognize the face. But that's not smart as we just print the name in stdout. In this video we'll improve it a bit, we will make our program greet person using Text to Speech and recognize user's voice using Google Speech API.

Let's start with Text to Speech.

We will use the code we've written last time.

Google Translate API has an option to get audio from text. It's easy to implement this API call, and I already prepared a package for easier use.

```
go get github.com/plutov/htgo-tts
```

This package caches results, so when we need to have a speech of same text, it will just play a record. To play the audio file package uses `omxplayer` which is already installed on Raspberry Pi device:

Let's change our main.go file to greet person when face is recognized, but let's do it only 1 time in 12 hours.

```
var (
	greetings = make(map[string]time.Time)
)

func isGreeted(name string) bool {
	g, ok := greetings[name]
	return ok && time.Now().Before(g.Add(time.Hour*12))
}
```

I don't use mutex here, because all access to `greetings` var are done synchronously.

Also let's check confidence. It should be higher than 50%:

```
speech := htgotts.Speech{Folder: "audio", Language: "en"}

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

		break
	}
}
```

I don't have portable speakers, so I connected my Raspi by HDMI to TV to be able to play the speech.

As you remember to build we need to set GOARCH and GOOS:

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
./capture
```

Now let's parse user's voice input. There are a lot of command line tools to record audio, let's go with `sox`.

```
sudo apt-get install pulseaudio sox
```

`rec` command line tool will be also installed.

As we will send this audio to Google Speech API, we need to set a specific format of audio. It should have 16000 rate, 1 channel.

Here is the command to do it in shell:

We set AUDIODEV env variable, it may be different on your device. I am using microphone built-in to webcam.

```
AUDIODEV=hw:1,0 rec -r 16000 -c 1 test1.wav trim 0 3
```

We set 3s limit.

We will execute this command from Go using exec package, let's write a function for it:

```
func record(fileName string, timeLimitSecs int) (err error) {
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "trim", "0", strconv.Itoa(timeLimitSecs))

	env := os.Environ()
	env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env

	return cmd.Run()
}

...

file := fmt.Sprintf("%d.wav", time.Now().UnixNano())
err = record(file, 3)
if err != nil {
	log.Printf("unable to record user voice input: %v", err)
	continue
}
```

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
./capture
```

Now let's send our audio to Google Speech API. First of all we need to create an application in Google Developer Console. Then create a service account and download a server JSON key I already downloaded it and sent to my Raspberyy Pi home folder.

And we need to get Go packages to work with Speech API:

```
go get cloud.google.com/go/speech/apiv1 google.golang.org/api/option google.golang.org/genproto/googleapis/cloud/speech/v1
```

Then we can use `speech.json` file as a token and get text transcript:

```
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
```

And let's make a simpliest parsing of yes|no keywords and reply something back:

```
text, err := speechToText(file)
if err != nil {
	log.Printf("unable to parse user voice: %v", err)
	continue
}

log.Printf("speech to text results: %s", text)

if strings.Contains(text, "yes") {
	speech.Speak(fmt.Sprintf("Have a nice day %s", f.Name))
}

if strings.Contains(text, "no") {
	speech.Speak("Oops, I am sorry")
}
```

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
./capture
```

I hope it was helpful and interesting, In the next episode we will add Natural Language Understanding feature to this program, so stay tuned.