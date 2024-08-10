# build stage
FROM golang:1.22 as builder
WORKDIR /root
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go

# final stage
FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/api .
CMD ["./api"]
