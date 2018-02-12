### Building Face Detection Go program for Raspberry Pi. Part 2

In previous episode we've build a program in Go which captures image from Webcam and uses Facebox to recognize the face. But that's not smart as we just print the name in stdout. In this video we'll improve it a bit, we will make our program greet person using Text to Speech and recognize user's voice using Google Speech to Text voice recognition.

We will use the code we've written last time.

As you remember to build we also need to set GOARCH and GOOS:

```
GOARCH=arm GOOS=linux go build -o capture
rsync capture pi@192.168.1.49:~/
```

I hope it was helpful and interesting, see you later!