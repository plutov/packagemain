FROM golang:alpine
ADD . /go
WORKDIR /go
RUN go build main.go
ENTRYPOINT ./main