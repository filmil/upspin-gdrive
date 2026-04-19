# Upspin `drive` repository

[![Build and Test](https://github.com/filmil/upspin-gdrive/actions/workflows/build-test.yml/badge.svg)](https://github.com/filmil/upspin-gdrive/actions/workflows/build-test.yml)

## GitHub Workflows

- **Build and Test**: Runs `bazel build //...` and `bazel test //...` on every push to `main` and all pull requests.

## Container Images

This project provides rules to build and push minimal Docker containers for the Upspin Drive components.

### Building Container Images

To build the container images:
```bash
bazel build //containers/...
```

To load an image into your local Docker daemon:
```bash
bazel run //containers:upspinserver_tarball
bazel run //containers:setupstorage_tarball
```

### Pushing Container Images

To push the images to Docker Hub (requires `docker login`):
```bash
bazel run //containers:upspinserver_push
bazel run //containers:setupstorage_push
```

---

# Upspin `drive` repository

Note: This repository is under construction.

This repository contains support for running Upspin on
[Google Drive](https://google.com/drive).

See the [master repository](https://github.com/upspin/upspin#readme) for more information.
