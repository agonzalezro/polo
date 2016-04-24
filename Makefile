all: polo

deps:
	mkdir -p bin

polo: deps
	cd cmd/polo&&go generate&&go build&&mv polo ../../bin
