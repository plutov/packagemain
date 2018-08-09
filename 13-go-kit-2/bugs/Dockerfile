FROM golang

RUN mkdir -p /go/src/github.com/plutov/packagemain/13-go-kit-2

ADD . /go/src/github.com/plutov/packagemain/13-go-kit-2

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run github.com/plutov/packagemain/13-go-kit-2/bugs/cmd  -watch github.com/plutov/packagemain/13-go-kit-2/bugs
