FROM golang:1.16-alpine AS builder

WORKDIR /go/src/goshrink

COPY cmd/server cmd/server
COPY internal internal
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY .env .env

RUN go mod tidy
RUN go get -d -v ./...
RUN cd cmd/server && go install -v ./...

CMD ["server"]

# # TODO: multistage build
# FROM scratch
# COPY --from=builder ["/go/bin/server", "/"]
# CMD ["/server"]
