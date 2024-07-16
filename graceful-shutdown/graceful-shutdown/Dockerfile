FROM golang:1.22 as builder

WORKDIR /root

COPY . .
RUN go mod download
RUN	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root
COPY --from=builder /root/server .

CMD ["./server"]