FROM golang:1.11

RUN mkdir -p /go/src/github.com/plutov/go-snake-telnet

WORKDIR /go/src/github.com/plutov/go-snake-telnet

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o snake-telnet main.go

ENTRYPOINT ["./snake-telnet"]