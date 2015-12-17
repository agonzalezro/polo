#!/bin/bash

OSS=(darwin freebsd linux)
ARCHS=(386 amd64)

mkdir -p bin
rm -f bin/polo*

cd cmd/polo

for os in "${OSS[@]}"; do
    for arch in "${ARCHS[@]}"; do
        GOOS=$os GOARCH=$arch go build
        mv polo ../../bin/polo-$os-$arch
    done
done

cd - > /dev/null
