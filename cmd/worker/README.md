# Worker

The package `worker` defines the full worker implementation.

It connects to the queue to receive and consume message from through `pkg/queue` interface.

The image processing logic relies on `pkg/imaging` which wraps [disintegration/imaging](https://github.com/disintegration/imaging) package.

The worker makes requests to the file storage server and directly writes into the database upon successful uploads.

## CLI configuration

It is possible to configure the worker with the command line flag `w`.

It defines the width to use when rescaling the image. The aspect ratio will always be preserved. If omitted, the width is defaulted to 200 pixels.

```sh
go run $(ls -1 ./cmd/worker/*.go | grep -v _test.go) -w 300
```

> Note: this syntax using a combination of `ls -1` and `grep -v` assures that this command will run even if test files are added in the future.
