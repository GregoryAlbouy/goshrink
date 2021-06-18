# HTTP

The package `server` implements all services related to HTTP and routing.

Each entities have their dedicated handlers. They also have a helper function for registering said handlers into the server's router.

`request.go` and `response.go` provide reusable functions for handling requests and manipulating response objects.

The package also provides its own error definitions and methods in `error.go`.
