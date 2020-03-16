# Fuzzball

[Fuzzball](https://github.com/sylabs/fuzzball-service) enables programmatic management of high performance compute resources.

## TLDR

```bash
helm install my-release .
```

## Introduction

This chart bootstraps a [Fuzzball](https://github.com/sylabs/fuzzball-service) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Helm 3.0+
- Kubernetes 1.14+

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
helm install my-release .
```

The command deploys Fuzzball on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Fuzzball chart and their default values.

| Parameter                    | Description                                           | Default                                                  |
| ---------------------------- | ----------------------------------------------------- | -------------------------------------------------------- |
| `replicaCount`               | Replica count                                         | `1`                                                      |
| `image.repository`           | Fuzzball image name                                   | `server`                                                 |
| `image.tag`                  | Fuzzball image tag                                    | `{{ .Chart.AppVersion }}`                                |
| `image.pullPolicy`           | Image pull policy                                     | `IfNotPresent`                                           |
| `image.pullSecrets`          | Specify registry secret names as an array             | `[]`                                                     |
| `nameOverride`               | String to partially override fuzzball.fullname template with a string (will prepend the release name) | `nil`    |
| `fullnameOverride`           | String to fully override fuzzball.fullname template with a string                                     | `nil`    |
| `serviceAccount.create`      | Specifies whether a ServiceAccount should be created  | `true`                                                   |
| `serviceAccount.annotations` | Annotations for the created ServiceAccount            | `{}`                                                     |
| `serviceAccount.name`        | The name of the ServiceAccount to create              | Generated from fuzzball.fullname template                |
| `podSecurityContext`         | Pod-level security context values                     | `{}`                                                     |
| `securityContext`            | Container-level security context values               | `{}`                                                     |
| `service.type`               | Kubernetes service type                               | `ClusterIP`                                              |
| `service.port`               | Kubernetes service port                               | `8080`                                                   |
| `ingress.enabled`            | Enable ingress controller resource                    | `false`                                                  |
| `ingress.annotations`        | Ingress annotations                                   | `{}`                                                     |
| `ingress.hosts[0].host`      | Hostname for the ingress                              |                                                          |
| `ingress.hosts[0].paths`     | Path within the url structure                         |                                                          |
| `ingress.tls[0].secretName`  | TLS secret to use (must be manually created)          |                                                          |
| `ingress.tls[0].hosts`       | List of FQDNs the above secret is associated with     |                                                          |
| `resources`                  | CPU/Memory resource requests/limits                   | `{}`                                                     |
| `nodeSelector`               | Node labels for pod assignment                        | `{}`                                                     |
| `tolerations`                | List of node taints to tolerate                       | `[]`                                                     |
| `affinity`                   | List of affinities                                    | `{}`                                                     |
| `mongodb.enabled`            | Enable MongoDB                                        | `true`                                                   |
| `mongodb.mongodbDatabase`    | MongoDB database to use                               | `server`                                                 |
| `mongodb.mongodbUsername`    | MongoDB username to use                               | `server`                                                 |
| `mongodb.mongodbPassword`    | MongoDB password to use                               | `changeme`                                               |
| `nats.enabled`               | Enable NATS                                           | `true`                                                   |
| `nats.auth.user`             | NATS username to use                                  | `server`                                                 |
| `nats.auth.password`         | NATS password to use                                  | `changeme`                                               |
| `redis.enabled`              | Enable Redis                                          | `true`                                                   |
| `redis.password`             | Redis password to use                                 | `changeme`                                               |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
helm install my-release --set mongodb.mongodbUsername=admin,mongodb.mongodbPassword=secretpassword .
```

The above command sets the MongoDB username and password to `admin` and `secretpassword` respectively.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
helm install my-release -f values.yaml .
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Using the Default Install

By default, this chart does not expose endpoints for [NATS](https://nats.io) (used by [Fuzzball Agents](https://github.com/sylabs/fuzzball-agent)) or the Fuzzball API (used by clients such as [fuzzctl](https://github.com/sylabs/fuzzctl)). In a production setting, NATS and the Fuzzball API would normally be exposed via a load balancer and ingress respectively. In a development setting, however, it can be simpler to use [Kubernetes port forwarding](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/).

To expose NATS (assuming release name `my-release`):

```bash
kubectl port-forward svc/my-release-nats-client 4222:4222
```

With NATS exposed, run the [Fuzzball Agent](https://github.com/sylabs/fuzzball-agent) with the correct credentials (by default, `server`:`changeme`):

```bash
fuzzball-agent -nats_uris nats://server:changeme@127.0.0.1:4222
```

To expose the Fuzzball API (assuming release name `my-release`):

```bash
kubectl port-forward svc/my-release-fuzzball 8080:8080
```

With the Fuzzball endpoint exposed, you can now run [`fuzzctl`](https://github.com/sylabs/fuzzctl) commands as usual:

```bash
fuzzctl login
fuzzctl list
...
```

## Unit Tests

This chart uses using the [helm-unittest](https://github.com/rancher/helm-unittest) plugin for unit testing of Helm Charts (note that this is a Rancher fork of the plugin, since it has a version that supports Helm 3). If you wish to run these tests locally, install the plugin:

```sh
helm plugin install https://github.com/rancher/helm-unittest --version v0.1.7-rancher1
```

You can then test using the `unittest` command:

```sh
helm unittest .
```
