
export GOPATH := $(GOPATH):$(PWD)

SRC=src/epl/*.go

.PHONY: all deps test

all: test

deps:

test:
	go test epl -test.v

