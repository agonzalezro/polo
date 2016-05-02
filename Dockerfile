FROM golang:1.6.2-alpine

ENV APP $GOPATH/src/github.com/agonzalezro/polo
RUN mkdir -p $APP
WORKDIR $APP

ADD glide.yaml $APP/glide.yaml
ADD glide.lock $APP/glide.lock
RUN apk add --no-cache git \
    && go get -u github.com/Masterminds/glide/... \
    && go get -u github.com/jteeuwen/go-bindata/... \
    && glide install \
    && apk del git

ADD . $APP
RUN apk add --no-cache make \
    && make \
    && apk del make

ENTRYPOINT ["bin/polo"]
CMD ["--help"]
