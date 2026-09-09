FROM golang:alpine AS builder
WORKDIR /srv

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./

RUN apk add --no-cache make
RUN make build
