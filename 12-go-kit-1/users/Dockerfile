FROM golang:1.20

WORKDIR /go/src/github.com/plutov/packagemain/12-go-kit-1/users

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

CMD [ "./server" ]
