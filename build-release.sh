#!/bin/bash

OSS=(darwin freebsd linux)
ARCHS=(386 amd64)

mkdir -p bin
rm -f bin/polo*

for os in "${OSS[@]}"; do
    for arch in "${ARCHS[@]}"; do
    	echo "Building for $os($arch)"
        GOOS=$os GOARCH=$arch make
        mv bin/polo bin/polo-$os-$arch
    done
done

# Link darwin amd64 to bin/polo
ln -s polo-darwin-amd64 bin/polo
