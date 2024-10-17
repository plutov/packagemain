FROM golang:1.21

WORKDIR /go/src/github.com/plutov/packagemain/12-go-kit-1/notificator

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

CMD [ "./server" ]
