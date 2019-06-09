FROM golang:1.12-alpine
ENV GO111MODULE=on
WORKDIR /go/src/github.com/nornir-automation/gornir
ADD . .
RUN apk add git
RUN go mod download
