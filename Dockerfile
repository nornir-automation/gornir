FROM golang:1.12-alpine
RUN apk add git
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
WORKDIR /go/src/github.com/nornir-automation/gornir
ADD . .
RUN go mod download > /dev/null
