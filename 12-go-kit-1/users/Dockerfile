FROM golang

RUN mkdir -p /go/src/github.com/plutov/packagemain/12-go-kit-1

ADD . /go/src/github.com/plutov/packagemain/12-go-kit-1

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run github.com/plutov/packagemain/12-go-kit-1/users/cmd  -watch github.com/plutov/packagemain/12-go-kit-1/users
