FROM golang:latest

# runtime dependency
RUN go get github.com/beego/bee github.com/xiorcale/rubus-api

WORKDIR /go/src/github.com/xiorcale/rubus-api
