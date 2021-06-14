# Image

The package `image` defines the image processing implementation.

In practice, it is injected into the worker as a dependency to pass it a job.

The pacakge is clearly split from our main application as it acts as a dependency or a plug-in (i.e. we could swap `image` package for a third party image processing library).
