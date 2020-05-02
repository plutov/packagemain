FROM golang:1.13

ENV GO111MODULE=on

WORKDIR /go/src/github.com/plutov/packagemain/12-go-kit-1/notificator

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o server ./cmd/main.go

CMD [ "./server" ]
