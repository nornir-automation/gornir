FROM golang:1.12-alpine
RUN apk add git
ENV GO111MODULE=on
WORKDIR /go/src/github.com/nornir-automation/gornir
ADD . .
RUN go mod download
