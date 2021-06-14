# Shrinker

Shrinker is a web api using a message queue in order to offload heavy computing tasks (namely image processing) to a worker.

## Running the project

```sh
# server
go run ./cmd/server/main.go

# worker
go run ./cmd/worker/main.go
```

## Infrastructure

![infrastrucute schema](docs/infrastructure.svg)

## Control flow

The control flow of the update of a user's avatar is as follow:

```txt
POST /avatars
```

![control flow schema](docs/control_flow.svg)

## Project structure

The main functional packages are structured this way:

```txt
.
├── cmd
│   ├── server
│   └── worker
├── internal
│   ├── database
│   ├── http
│   └── storage
└── pkg
    ├── queue
    └── image
```

`internal` package holds the definitions of the business entities at its root.
