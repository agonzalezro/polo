docker build -t gb-cross .

TARGETS=(darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 linux/386 linux/amd64)

for target in ${TARGETS[@]}; do
  rm -rf pkg # to force rebuilding of vendored packages
  export GOOS=${target%/*}
  export GOARCH=${target##*/}
  docker run -v `pwd`:/app -e GOOS=$GOOS -e GOARCH=$GOARCH gb-cross
  mv bin/polo bin/polo_"$GOOS"_"$GOARCH"
done
