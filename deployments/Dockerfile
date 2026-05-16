FROM golang:alpine AS builder
WORKDIR /srv
COPY . .
RUN apk add --no-cache make
RUN make build
