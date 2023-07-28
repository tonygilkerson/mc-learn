# basic-app

A simple deployment, service and ingress created with a custom resource using Metcontroller.

The example in this folder is based on vfarcic's [Metacontroller gist](https://gist.github.com/vfarcic/8adaf8fd6496bc99a466ba55834e1838)

> **NOTE** For this exercise run all commands from the `basic-app` folder.

## Cluster Setup

Create a cluster an install Metacontroller, see root level [README](../README.md)

```sh
# Run all commands from here
cd basic-app

# Create namespaces
kubectl create namespace production
kubectl create namespace controllers
```

## My BasicApp

Create an instance of my basic app

```sh
# Apply
kubectl apply -f crds
kubectl apply -f examples/my-basic-app.yaml

# Verify
kubectl -n production get basicapps
```

At this point `my-basic-app` is deployed into the cluster but there is no workload/implementation for it. We need to crate a controller to bring it to life. To do that we will apply `composite-controllers.yaml` which will create a `deployment`, `service` and `ingress` resources.

## BasicApp Composite Controller

```sh
# Apply
kubectl apply -f manifests/composite-controller.yaml

# Verify
kubectl get compositecontrollers
```

At this point we just created a composit controller called `basic-app` which can see `my-basic-app` but if we look at th metacontroller logs we see an issue. The composite controller is receiving events from the Kubernetes API and in response it is sending the desired state to the webhook we definde in our composite conttoller but we do not have anthing to receive the webhook request.

```sh
kubectl -n metacontroller -l app.kubernetes.io/instance=metacontroller logs | grep "not found"

# You will see the webhook failed
failed to sync BasicApp 'production/my-basic-app': sync hook failed for BasicApp production/my-basic-app: sync hook failed: http error: Post \"http://basicapp-controller.controllers/sync\"
```

## Webhook Deployments (aka controller)

We need to create a deployment that can handle the incoming webhooks and implement the reconciliation between the desired and actual state. This is our controller logic.

### Container Image

Before we get started we need to create a container image for our basic app deployment

```sh
# Use ttl.sh anonymous and ephemeral image registry for simplicity
# Guess I am to lazy to stand up a registry :-)
UUID=$(uuidgen)
IMAGE_NAME="${UUID,,}" # convert to lower case
REG="ttl.sh/${IMAGE_NAME}:2h"
podman build -t "$REG" .
podman push "$REG"

# Update the Deployment spec to point at our ephemeral image
reg="$REG" yq --inplace 'select(.kind == "Deployment") | .spec.template.spec.containers[0].image = env(reg)' manifests/deployment.yaml 
```

### Deployment

We are now ready to apply the `Deployment` and `Service` that will handle our webhooks.

```sh
# Apply
kubectl apply -f manifests/deployment.yaml
kubectl apply -f manifests/service.yaml

# Verify
kubectl -n controllers get deployments,services
```

## Verify

Now that we have our webhooks in place, let's recreate the apps to generate events and verify if things worked.

```sh
# Recreate the app just to trigger events
kubectl delete -f examples/my-basic-app.yaml
kubectl apply -f examples/my-basic-app.yaml
```

If all went well the controller received the "basicapp created" event, called the webhook and created the app. To verify let's hit the app.

```sh
# DEVTODO setup an ingress controller but until then we can use port-forward
kubectl -n production port-forward svc/basic-app 8080:80

open http://localhost:8080/
```
