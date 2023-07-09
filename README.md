# mc-learn

This project is a place form me to play around and learn about [Metacontroller](https://metacontroller.github.io/metacontroller/intro.html).

## Setup

In most cases you will need to complete this setup in order to create a cluster and install Metacontroller before going into the exercise subfolders.

```sh
podman machine start
kind create cluster
```

To install Metacontroller you can use [Helm](https://metacontroller.github.io/metacontroller/guide/install.html#install-metacontroller-using-helm) or [Kustomize](https://metacontroller.github.io/metacontroller/guide/install.html#install-metacontroller-using-kustomize).

```sh
kubectl create ns metacontroller

helm -n metacontroller install metacontroller oci://ghcr.io/metacontroller/metacontroller-helm --version=v4.10.4
```
