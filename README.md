# Compute Service

The Sylabs Compute Service enables programmatic management of high performance compute resources.

## Quick Start

Configure your go environment to pull private go modules, by forcing `go get` to use `git+ssh` instead of `https`. This lets the go compiler pull private dependencies using your machine's ssh keys.

```sh
git config --global url."ssh://git@github.com/sylabs".insteadOf "https://github.com/sylabs"
```

If using Go 1.13, the `go` command defaults to downloading modules from the public Go module mirror, and validating downloaded modules against the public Go checksum database. Since private Sylabs projects are not availble in the public mirror nor the public checksum database, we must tell Go about this. One way to do this is to set `GOPRIVATE` in the Go environment:

```sh
go env -w GOPRIVATE=github.com/sylabs
```

Install `go-bindata` to help generate the GraphQL schema later:

```sh
go get -u github.com/jteeuwen/go-bindata/...
```

In order for go to execute this binary the path in `go env GOPATH` needs to be included in your `PATH`.

To run the server, you'll need MongoDB and NATs endpoints to point it to. If you don't have these already, you can start them with Docker easy enough:

```sh
docker run -d -p 27017:27017 mongo
docker run -d -p 4222:4222 nats
```

Finally, start the server:

```sh
$ go generate ./... && go run ./cmd/server/
INFO[0000] starting                                      name="Compute Server" org=Sylabs
INFO[0000] connecting to database
INFO[0000] database ready                                took=7.375268ms
INFO[0000] connecting to messaging system
INFO[0000] messaging system ready                        took=5.21203ms
INFO[0000] listening                                     addr="[::]:8080"
```

## Testing

### Unit Tests

Unit tests can be run like so:

```sh
go test ./...
```

### Integration Tests

To run integration tests, you'll need MongoDB and NATs endpoints to point it to. If you don't have these already, you can start them with Docker easy enough:

```sh
docker run -d -p 27017:27017 mongo
docker run -d -p 4222:4222 nats
```

Integration tests can then be run like so:

```sh
go test -tags=integration ./...
```
