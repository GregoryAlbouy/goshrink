# HTTP

The package `server` implements all services related to HTTP and routing.

Each entities have their dedicated handlers which are assigned to an endpoint in `routing.go`.

`request.go` and `response.go` provide reusable functions for handling requests and manipulating response objects.

The package also provides its own error definitions and methods in `error.go`.
