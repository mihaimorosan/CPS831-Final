FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /CPS831-Final
RUN mkdir -p /CPS831-Final/service

WORKDIR /CPS831-Final

COPY ./service/service.pb.go /CPS831-Final/service
COPY ./server.go   /CPS831-Final

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o CPS831-Final .
CMD ./CPS831-Final