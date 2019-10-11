FROM golang:1.12-alpine

RUN apk update && apk add --no-cache git

WORKDIR /opt/code/
ADD ./ /opt/code/

RUN go get
RUN go build -o bin/twelve-factor main.go

ENTRYPOINT ["./bin/twelve-factor"]