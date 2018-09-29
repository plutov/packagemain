FROM golang

RUN mkdir -p /go/src/github.com/plutov/packagemain/14-go-kit-3

ADD . /go/src/github.com/plutov/packagemain/14-go-kit-3

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run github.com/plutov/packagemain/14-go-kit-3/users/cmd  -watch github.com/plutov/packagemain/14-go-kit-3/users
