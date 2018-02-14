### Building Face Detection Go program for Raspberry Pi. Part 2

In previous episode we've build a program in Go which captures image from Webcam and uses Facebox to recognize the face. But that's not smart as we just print the name in stdout. In this video we'll improve it a bit, we will make our program greet person using Text to Speech and recognize user's voice using Google Speech to Text voice recognition.

Let's start with Text to Speech.

We will use the code we've written last time.

Google Translate API has an option to get audio from text. It's easy to implement this API call, but also there is a package to do it faster.

```
go get github.com/plutov/htgo-tts
```

This package is good because it caches results, so when we need to have a speech of same text, it will just play a record. To play it we'll need to install mplayer on Raspberry Pi device:

```
sudo apt-get install mplayer
```

Let's change our main.go file to greet person when face is recognized, but let's do it only 1 time in 12 hours.

```
var (
	greetings = make(map[string]time.Time)
)

func isGreeted(name string) bool {
	g, ok := greetings[name]
	now := time.Now()
	return ok && now.Before(g.Add(time.Hour*12))
}
```

I don't use mutex here, because all access to `greetings` var are done synchronously.

Also let's check confidence. It should be higher than 50%:

```
log.Printf("found %d faces", len(faces))

if f.Confidence >= 0.5 {
	greeted := isGreeted(f.Name)
	log.Printf("face: %s, confidence: %.2f, greeted: %t", f.Name, f.Confidence, greeted)
	if !greeted {
		greetings[f.Name] = time.Now()

		speech := htgotts.Speech{Folder: "audio", Language: "en"}
		err = speech.Speak(fmt.Sprintf("Hi %s, how are you today?", f.Name))
		if err != nil {
			log.Printf("unable to run text-to-speech: %v", err)
			continue
		}

		break
	}
}
```

As you remember to build we need to set GOARCH and GOOS:

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
```

Now let's parse user's voice input. There is a nice linux tool `sox` to record audio. Let's install it first:

```
sudo apt-get install pulseaudio sox
```

`rec` command line tool will be also installed. We need to stop recording when user stop speaking, it can be done with `silence` option. Also we need to set a time limit, for example 5s.

As we will send this audio to Google Speech API, we need to set a specific format of audio. It should have 16000 rate, 1 channel.

Here is the command to do it in shell:

We set AUDIODEV env variable, it may be different on your device. I am using microphone built-in to webcam.

```
AUDIODEV=hw:1,0 rec -r 16000 -c 1 test1.wav trim 0 5 silence 0 1 0.5 3%
```

0 is to avoid losing all silence.
1 means that we enabled silence detection at the end
0.5 - silence time to stop
3% - volume threshold

We will execute this command from Go using exec package, let's write a function for it:

```
func record(fileName string, timeLimitSecs int) (err error) {
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "trim", "0", strconv.Itoa(timeLimitSecs), "silence", "0", "1", "0.5", "3%")

	env := os.Environ()
	env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	return cmd.Run()
}

...

file := fmt.Sprintf("%d.wav", time.Now().UnixNano())
err = record(file, 5)
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

I hope it was helpful and interesting, see you later!