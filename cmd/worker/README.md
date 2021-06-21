# Worker

The package `worker` defines the full worker implementation.

It connects to the queue to receive and consume message from through `pkg/queue` interface.

The image processing logic relies on `pkg/imaging` which wraps [disintegration/imaging](https://github.com/disintegration/imaging) package.

The worker makes requests to the file storage server and directly writes into the database upon successful uploads.
