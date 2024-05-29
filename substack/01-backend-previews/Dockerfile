FROM golang:1.22 as builder

WORKDIR /
COPY main.go .
RUN go build -o goapp main.go

FROM alpine:latest
RUN apk add libc6-compat
COPY --from=builder /goapp .

CMD ["./goapp"]