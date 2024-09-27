FROM golang:1.23-alpine AS builder

RUN apk add build-base

WORKDIR /
COPY main.go .
RUN go build -o app main.go

FROM alpine:latest
COPY --from=builder /app .

CMD ["./app"]
