# tools
GO = GO111MODULE=on GOPROXY=https://proxy.golang.org go

.PHONY: default build test install version

default: build

build: version.go
	${GO} build ./...

test: build
	${GO} test -v -race ./...

install: test
	${GO} install ./...

clean:
	${GO} clean ./...

# When releasing significant changes, make sure to update the semantic
# version number in `./VERSION`, merge changes, then run `make release_tag`.
version:
ifdef ver
		@echo 'Setting version to $(ver)'
		@./tag_version.sh $(ver)
else
		@echo 'ver not defined. call make ver=<version eg 1.2.3> version'
endif

release_tag:
	git tag -a v`cat ./VERSION`
	git push origin v`cat ./VERSION`
