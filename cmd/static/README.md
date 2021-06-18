# Static file server

The package `server` holds the code necessary to create and run a server capable of storing and serving static files.

For a production-level project, this role would be delegated to a cloud solution.

## Endpoints

### Serve an image (in practice a user's avatar)

It is a public endpoint, anybody can retrieve and view a user's avatar.

Request:

```sh
curl -X GET http://localhost:9999/static/{UUID}.{png,jpeg}
```

Response:

```txt
200 OK
<file>
```

### Upload an image

It a restricted route, only our worker can access it. It need to provide an API key in the request Authorization headers.

Request:

```sh
curl -X POST -H "Authorization: Bearer goshrink" -H "Content-Type:multipart/form-data" -F "upload=@fixtures/sample.png" http://localhost:8000/static/avatar
```

Response:

```txt
201 Created
```
