FROM golang:1.4.2-cross

WORKDIR /app
VOLUME /app

RUN go get github.com/constabulary/gb/...

CMD ["gb", "build"]
