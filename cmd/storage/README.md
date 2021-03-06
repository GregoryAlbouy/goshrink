# Storage server

The package `storage` holds the code necessary to create and run a server capable of storing and serving static files.

For a production-level project, this role would be delegated to a cloud solution.

## Endpoints

### Serve an image (in practice a user's avatar)

It is a public endpoint, anybody can retrieve and view a user's avatar.

Request:

```sh
curl -X GET http://localhost:9999/storage/<user_id>/<file_uuid>.png
```

Response:

```txt
200 OK
<file>
```

### Upload an image

It is a restricted route, only our worker can access it. It needs to provide an API key in the request Authorization headers.

Only one image is stored on the server for one user at a time. If a user uploads a new image, it will overwrite the image currently stored if it exists.

Request:

```sh
curl -X POST -H "Authorization: Bearer goshrink" -H "Content-Type:multipart/form-data" -F "upload=@fixtures/sample.png" http://localhost:9998/storage/avatar
```

Response:

```txt
201 Created
```
