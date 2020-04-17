FROM golang:latest

# runtime dependency
RUN go get github.com/beego/bee github.com/kjuvi/rubus-api

WORKDIR /go/src/github.com/kjuvi/rubus-api
