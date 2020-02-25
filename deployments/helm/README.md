# Helm Chart

The Helm chart in this directory simplifies deployment of the service on Kubernetes.

## Prerequisites

The following prerequisites are required to develop and/or use the Helm chart.

### Helm

Install [Helm](https://helm.sh) according to the [docs](https://helm.sh/docs/intro/install/). The minimum required version is 3.0.

### Helm Unit Test Plugin

This repo uses using [helm-unittest](https://github.com/lrills/helm-unittest) plugin for unit testing of Helm Charts. If you wish to run these tests locally, install the plugin:

```sh
helm plugin install https://github.com/lrills/helm-unittest
```

## Deployment

### Install Service

Build chart dependencies, and install the service:

```sh
helm dep build fuzzball/
helm install <name> fuzzball/
```

### Upgrade Service

Optionally, update chart dependencies:

```sh
helm dep update fuzzball/
```

Upgrade the service:

```sh
helm upgrade <name> fuzzball/
```

### Uninstall Service

Delete the service:

```sh
helm delete <name>
```

## Helm Chart Unit Tests

You can test using the [helm-unittest](https://github.com/lrills/helm-unittest) plugin. For example, to test the `fuzzball` chart:

```sh
helm unittest fuzzball/
```
