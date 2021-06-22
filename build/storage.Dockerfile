FROM golang:1.16-alpine AS builder

WORKDIR /go/src/goshrink

COPY cmd/storage cmd/storage
COPY internal internal
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY .env .env

RUN mkdir storage
RUN go mod tidy
RUN go get -d -v ./...
RUN cd cmd/storage && go install -v ./...

CMD ["storage"]
