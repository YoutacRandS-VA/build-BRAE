[![Go Reference](https://pkg.go.dev/badge/golang.org/x/build/cmd/gopherbot.svg)](https://pkg.go.dev/golang.org/x/build/cmd/gopherbot)

# golang.org/x/build/cmd/gopherbot

The gopherbot command runs Go's gopherbot role account on GitHub and Gerrit.

## Development

To test changes locally, run:

```sh
$ go run . -dry-run
```

To connect gopherbot to development instances of, e.g. devapp, modify the
source code to point at those instances.

## Development with Docker

```
make docker-image
docker volume create golang-maintner
docker run -v golang-maintner:/.cache/golang-maintner \
    -it --rm gcr.io/go-dashboard-dev/gopherbot --dry-run
```

## Deployment

See the documentation on [deployment](../../doc/deployment.md).