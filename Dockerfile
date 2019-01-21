FROM golang:1.11

ENV DEBIAN_FRONTEND noninteractive
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=1
ENV GOROOT=/usr/local/go
ENV GOPATH=/go

# this sucks!
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | env INSTALL_DIRECTORY=/usr/bin bash
