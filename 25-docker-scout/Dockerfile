FROM golang:1.19

WORKDIR /
COPY main.go .
RUN go build -o goapp main.go

CMD ["./goapp"]