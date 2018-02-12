### Building Face Detection Go program for Raspberry Pi. Part 2

In previous episode we've build a program in Go which captures image from Webcam and uses Facebox to recognize the face. But that's not smart as we just print the name in stdout. In this video we'll improve it a bit, we will make our program greet person using Text to Speech and recognize user's voice using Google Speech to Text voice recognition.

Let's start with Text to Speech.

We will use the code we've written last time.

Google Translate API has an option to get audio from text. It's easy to implement this API call, but also there is a package to do it faster.

```
go get github.com/hegedustibor/htgo-tts
```

This package is good because it caches results, so when we need to have a speech of same text, it will just play a record. To play it we'll need to install mplayer on Raspberry Pi device:

```
apt-get install mplayer
```

Let's change our main.go file to greet person when face is recognized, but let's do it only 1 time in 12 hours.

```
var (
	greetings = make(map[string]time.Time)
)

func isGreeted(name string) bool {
	g, ok := greetings[name]
	now := time.Now()
	return ok && now.Before(g.Time.Add(time.Hour*12))
}
```

I don't use mutex here, because all access to `greetings` var are done synchronously.

Also let's check confidence. It should be higher than 50%:

```
if f.Confidence >= 0.5 {
	greeted := isGreeted(f.Name)
	log.Printf("face: %s, confidence: %.2f, greeted: %t", f.Name, f.Confidence, greeted)
	if !greeted {
		greetings[f.Name] = time.Now()
		speech := htgotts.Speech{Folder: "audio", Language: "en"}
		speech.Speak(fmt.Sprintf("Hi %s, how are you today?", f.Name))
	}
}
```

As you remember to build we need to set GOARCH and GOOS:

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
```

I hope it was helpful and interesting, see you later!