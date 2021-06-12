# Image

The package `image` defines the image processing implementation.

In practice, it is injected into the worker as a dependency to pass it a job.

It should not import any types from internal. Indeed, `pkg/` folder can be viewed as publicly available and we would not want to leak our internals.

It also clearly slipts our dependencies and act as a plug-in (i.e. we could swap `image` for a thrid party image processing library).
