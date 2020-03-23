# Fuzzball Service

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
[![CI Workflow](https://github.com/sylabs/fuzzball-service/workflows/ci/badge.svg)](https://github.com/sylabs/fuzzball-service/actions)
[![Dependabot](https://api.dependabot.com/badges/status?host=github&repo=sylabs/fuzzball-service&identifier=232618704)](https://app.dependabot.com/accounts/sylabs/repos/232618704)

Fuzzball enables programmatic management of high performance compute resources.

## Quick Start

Ensure that you have one of the two most recent minor versions of Go installed as per the [installation instructions](https://golang.org/doc/install).

Install [Mage](https://magefile.org) as per the [installation instructions](https://magefile.org/#installation).

To run the server, you'll need MongoDB, NATS, and Redis endpoints to point it to. If you don't have these already, you can start them with Docker easy enough:

```sh
docker run -d -p 27017:27017 mongo
docker run -d -p 4222:4222 nats
docker run -d -p 6379:6379 redis
```

Finally, run the server:

```sh
mage run
```

## Testing

### Unit Tests

Unit tests can be run like so:

```sh
mage unittest
```

### Integration Tests

To run integration tests, you'll need MongoDB, NATS, and Redis endpoints to point it to. If you don't have these already, you can start them with Docker easy enough:

```sh
docker run -d -p 27017:27017 mongo
docker run -d -p 4222:4222 nats
docker run -d -p 6379:6379 redis
```

Integration tests can then be run like so:

```sh
mage test
```

## License

This project is licensed under a 3-clause BSD license found in the [license file](LICENSE.md).
