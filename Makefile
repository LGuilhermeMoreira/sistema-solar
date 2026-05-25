BINARY ?= sistema-solar
GOCACHE ?= $(CURDIR)/.cache/go-build

.PHONY: all build clean

all: build

build:
	GOCACHE="$(GOCACHE)" go build -buildvcs=false -o $(BINARY) .

clean:
	rm -f $(BINARY)
