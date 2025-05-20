FROM golang:1.24 as builder

WORKDIR /build

COPY . .
RUN go mod download
RUN	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /release
COPY --from=builder /build/server .

CMD ["./server"]
