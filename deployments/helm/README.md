# Helm Chart

The Helm chart in this directory simplifies deployment of the service on Kubernetes.

## Prerequisites

The following prerequisites are required to develop and/or use the Helm chart.

### Helm

Install [Helm](https://helm.sh) according to the [docs](https://helm.sh/docs/intro/install/). The minimum required version is 3.0.

### Helm Unit Test Plugin

This repo uses using the [helm-unittest](https://github.com/rancher/helm-unittest) plugin for unit testing of Helm Charts (note that this is a Rancher fork of the plugin, since it supports Helm 3). If you wish to run these tests locally, install the plugin:

```sh
helm plugin install https://github.com/rancher/helm-unittest --version v0.1.7-rancher1
```

## Deployment

### Install Service

Build chart dependencies, and install the service:

```sh
helm dep build fuzzball/
helm install fuzzball fuzzball/
```

To connect a [Fuzzball Agent](https://github.com/sylabs/fuzzball-agent), you must expose NATS:

```sh
kubectl port-forward svc/fuzzball-nats-client 4222:4222
```

With NATS exposed, run the Agent with the correct credentials (by default, `server`:`changeme`):

```sh
fuzzball-agent -nats_uris nats://server:changeme@127.0.0.1:4222
```

To connect to the service with [`fuzzctl`](https://github.com/sylabs/fuzzctl), you must expose the Fuzzball endpoint:

```sh
kubectl port-forward svc/fuzzball 8080:8080
```

With the Fuzzball endpoint exposed, you can now run `fuzzctl` commands as usual:

```sh
fuzzctl login
fuzzctl list
...
```

### Upgrade Service

Optionally, update chart dependencies:

```sh
helm dep update fuzzball/
```

Upgrade the service:

```sh
helm upgrade fuzzball fuzzball/
```

### Uninstall Service

Delete the service:

```sh
helm delete fuzzball
```

## Helm Chart Unit Tests

You can test using the [helm-unittest](https://github.com/rancher/helm-unittest) plugin. For example, to test the `fuzzball` chart:

```sh
helm unittest fuzzball/
```
