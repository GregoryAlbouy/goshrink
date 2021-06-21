# HTTP

The package `http` implements all services related to HTTP and routing.

Each entities have their dedicated handlers. They also have a helper function for registering said handlers into the server's router.

`request.go` and `response.go` provide reusable functions for handling requests and manipulating response objects.

The package also provides its own error definitions and methods in `error.go`.

## Routes

> Note: for quick testing, curl commands are provided.

### Retrieve a user

Request:

```sh
curl -X GET http://localhost:9999/users/1
```

Response:

```json
{
  "id": 1,
  "username": "Bret",
  "email": "Sincere@april.biz",
  "avatar_url": ""
}
```

### Create a user

Request:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"username": "user", "email": "user@mail.com", "password": "pkEfkV39Bs"}' http://localhost:9999/users
```

Response:

```txt
201 Created
```

### Login

Request:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"username": "user", "password": "pkEfkV39Bs"}' http://localhost:9999/login
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQxODMzMjAsImp0aSI6IjEifQ.gl2EQWCxszEHXkQkR4HIAhyhIxwufowdlaNJSOtIYek"
}
```

### Upload an avatar

Request:

```sh
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQxODMzMjAsImp0aSI6IjEifQ.gl2EQWCxszEHXkQkR4HIAhyhIxwufowdlaNJSOtIYek" -H "Content-Type:multipart/form-data" -F "upload=@fixtures/sample.jpeg" http://localhost:9999/users/1/avatar

```

Response:

```txt
202 Accepted
```
