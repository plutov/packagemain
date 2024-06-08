FROM golang:1.22 as builder

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy main.go

FROM alpine:latest
COPY --from=builder /proxy .

EXPOSE 3307

CMD ["./proxy"]