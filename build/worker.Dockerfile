FROM golang:1.16-alpine AS builder

WORKDIR /go/src/goshrink
COPY cmd/worker cmd/worker
COPY internal internal
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY .env .env

RUN go mod tidy
RUN go get -d -v ./...
RUN cd cmd/worker && go install -v ./...

CMD ["worker"]
