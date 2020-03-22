FROM golang:latest

WORKDIR /go/src/github.com/kjuvi/rubus-api

RUN go get github.com/beego/bee
