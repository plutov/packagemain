### Snake over Telnet in Go [![Build Status](https://travis-ci.org/plutov/go-snake-telnet.svg?branch=master)](https://travis-ci.org/plutov/go-snake-telnet)


### Run it with go

```bash
go get github.com/plutov/go-snake-telnet
go-snake-telnet
```

## Run with Docker

```bash
docker pull pltvs/go-snake-telnet .
docker run -d -p 8080:8080 pltvs/go-snake-telnet
```

## Play!

```bash
telnet localhost 8080
```

### Tests

```
go test ./... -bench=.
```
