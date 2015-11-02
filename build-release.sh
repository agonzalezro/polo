#!/bin/bash

OSS=(darwin freebsd linux)
ARCHS=(386 amd64)

rm bin/polo*

for os in "${OSS[@]}"; do
    for arch in "${ARCHS[@]}"; do
        GOOS=$os GOARCH=$arch gb build
    done
done

mv bin/polo bin/polo-darwin-amd64
