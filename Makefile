BINARY = $(notdir $(PWD))
BIND = 127.0.0.1:3000

VERSION := $(shell git describe --tags --dirty --always 2> /dev/null || echo "dev")
SOURCES = $(wildcard *.go **/*.go)

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags "-X main.version=`git describe --long --tags --dirty --always`" -o "$@"

deps:
	go get ./...

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY) --bind "$(BIND)" -vv

debug: $(BINARY)
	./$(BINARY) --pprof :6060 -vv

test:
	go test ./...
	golint ./...
