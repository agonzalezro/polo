#!/bin/bash

all:
	# Build polo
	make clean
	make link
	make build

build:
	GOPATH=$(shell pwd) go get github.com/agonzalezro/polo

clean:
	# Remove the generated/downloaded stuff
	rm -Rf pages *.html src/agonzalezro bin pkg

link:
	# Do the funny symbolic links again
	mkdir -p src/github.com/agonzalezro
	rm src/github.com/agonzalezro/polo
	ln -s $(shell pwd) src/github.com/agonzalezro/polo

install:
	mkdir -p $(HOME)/bin
	rm $(HOME)/bin/polo
	cp bin/polo $(HOME)/bin
