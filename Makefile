.PHONY: default routegen all clean

GOBIN = $(shell pwd)/build/bin
GO ?= latest

default: routegen

all: routegen

routegen:
	go build -o=${GOBIN}/$@
	@echo "Done building."

clean:
	rm -f ${GOBIN}/routegen
