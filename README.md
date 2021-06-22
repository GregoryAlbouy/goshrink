# Shrinker <!-- omit in toc -->

Shrinker is an application that uses a message broker to offload image processing tasks from a web API to a dedicated worker.

Authenticated users can make request to the API server to upload image file as their avatar. The images are resized to a smaller scale and saved to a file storage to be viewable by anyone.

We currently support images uploaded in `.png` or `.jpeg` format.

## Table of contents <!-- omit in toc -->

- [Getting started](#getting-started)
  - [Installation and dependencies](#installation-and-dependencies)
  - [Run the project with Docker Compose](#run-the-project-with-docker-compose)
  - [Run in development mode](#run-in-development-mode)
  - [Quick manual testing with provided scripts](#quick-manual-testing-with-provided-scripts)
- [Infrastructure](#infrastructure)
  - [The API server](#the-api-server)
  - [The message broker](#the-message-broker)
  - [The worker](#the-worker)
- [Control flow](#control-flow)
- [Architecture](#architecture)
  - [Folder structure](#folder-structure)
  - [`internal`](#internal)
  - [`pkg`](#pkg)
  - [`cmd`](#cmd)
  - [Miscellaneous](#miscellaneous)
- [Further documentation](#further-documentation)
- [Acknowledgement and contributors](#acknowledgement-and-contributors)

## Getting started

### Installation and dependencies

You must be able to run Docker and Docker Compose on your local machine to run this app. Refer to the [Get Docker](https://docs.docker.com/get-docker/) and [Install Docker Compose](https://docs.docker.com/compose/install/) docs for their installation.

You must provide a `.env` file inside the root directory.
For a quick start, you can use the values from the provided example:

```sh
echo "$(cat .env.example)" >> .env
```

The storage server expects to store and serve images from `./storage`. You must create this directory at the root of the project:

```sh
mkdir storage
```

### Run the project with Docker Compose

The app is fully dockerized. To start all instances required, simply run:

```sh
make docker-up
```

To stop run:

```sh
make docker-down
```

To restart run:

```sh
make docker-restart
```

You can also run each entity individually in Docker:

```sh
make [entity]
```

### Run in development mode

Before running in dev mode, make sure you installed Go module dependencies.

```sh
go get -u
```

You might want to use the Go command to run each part individually. For this, we also provide scripts to run them using the `go run` command.

```sh
make local-[entity]
# alias to go run ./cmd/entity/main.go
```

Note that the MySQL and RabbitMQ instances are always run through docker-compose using `make rabbitmq` and `make mysql`.

> Refer to [Makefile](/Makefile) for a reference of our scripts.

### Quick manual testing with provided scripts

In order to showcase the main point of this project, we provide some handy scripts in the for quick manual tests.

These scripts are only concerned with login you in as our test user and making an upload request as them.

Create a new user (`email` can be omitted as it is not relevant to our usecase).

```sh
./scripts/post-user <username> <password> [<email>]
# 201 Created
```

Login in as that user.

```sh
./scripts/post-login <username> <password>
# {"access_token": "<token>"}
```

Upload an avatar (`filepath` is defaulted to `fixtures/sample.png` if not specified).

```sh
./scripts/post-avatar <token> [<filepath>]
# 202 Accepted
```

You can now get the updated user and open in your browser the URL to see the processed image.

```sh
curl http://localhost:9999/users/<username>
# {
#   "id":1,
#   "username":"amdin",
#   "email":"admin@goshrink.com",
#   "avatar_url":"http://localhost:9998/storage/9f31c631-a868-4727-94f5-ccd30f0e3db7.png"
# }
```

## Infrastructure

We aim to mimic what a cloud based infrastructure would be, to minimize the overhead of swapping to such a deployment strategy later.

![infrastrucute schema](docs/infrastructure.svg)

Each of the 5 entites above runs in a docker container.

The 3 main applicative entities are the **API server**, the **worker** and the **message broker**.

### The API server

The API simply answers to client requests and interacts with the database and the message broker when resquests call for it.

### The message broker

The message broker is in charge of forwarding API messages to the worker for it to start its job.

It uses RabbitMQ.

### The worker

The worker runs its job then makes resquests to store files and can write in the database.

It leverages two customs modules, `queue` package built ontop of [streadway/amqp](https://github.com/streadway/amqp) and `imaging` which wraps [disintegration/imaging](https://github.com/disintegration/imaging) for image processing.

## Control flow

The main work flow of this project is the request for a user to upload an avatar.

```txt
POST /api/avatar

content type: multipart/form-data
body: <image_file>
```

![image upload flowchart](docs/control_flow.svg)

## Architecture

### Folder structure

The main functional packages for the project are:

```txt
.
├── cmd
│   ├── server
│   ├── storage
│   └── worker
├── internal
│   ├── database
│   └── http
└── pkg
    ├── imaging
    └── queue
```

### `internal`

The main application code, it can be viewed as private. Domain related type definitions are at its root and packages within are grouped by dependencies (database, http, etc.).

### `pkg`

Reusable library code. It is safe to use by external applications. It does not import any types from `internal` and does not rely on it to work. There is no domain logic inside this directory.

### `cmd`

Main applications for the project. Children direcories are each equivalent to one executable and are independant from one another. A `main` function may need to import and call code from `internal` or `pkg`. If so, it will ties all dependencies together.

### Miscellaneous

Other directories are self explanatory.

## Further documentation

Sub-packages have their own documentation when it is relevant. You may refer to these docs for further explanation: Notably:

- [http](internal/http/README.md) (comes with handy curl commands)
- [database](internal/database/README.md)
- [queue](pkg/queue/README.md)
- [storage](cmd/storage/README.md)
- [worker](cmd/worker/README.md)

## Acknowledgement and contributors

This project is part of an school assignement. Our team members are:

- Gregory Albouy
- Thomas Moreira
