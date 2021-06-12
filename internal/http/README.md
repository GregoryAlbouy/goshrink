# HTTP

The package `server` implements all services related to HTTP and routing.

Each entities have their dedicated handlers which are assigned to an endpoint in `routing.go`.

`reponse.go` and request.go provide reusable functions for manipulating request and response objects.

The package also provides its own error definitions and methods in `error.go`.
