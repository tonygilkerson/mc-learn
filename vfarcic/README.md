# vfarcic

The example in this folder is based on Vfarcic's [Metacontroller gist](https://gist.github.com/vfarcic/8adaf8fd6496bc99a466ba55834e1838)

> **NOTE** For this exercise run all commands from the `vfarcic` folder.

## Cluster Setup

Create a cluster an install Metacontroller, see root level [README](../README.md)

```sh
# Run all commands from here
cd vfarcic

# Create namespaces
kubectl create namespace production
kubectl create namespace controllers
```

## My app

Create an instance of my app

```sh
kubectl apply -f crds.yaml
kubectl apply -f slacks.yaml
kubectl apply -f my-app.yaml
```

At this point my app is deployed into the cluster but there is no workload/implementation for it. We need to crate a controller to bring it to life. To do that we will apply `composite-controllers.yaml` which will uses Metacontroller to create a `deployment`, `service` and `ingress` composite resource.

## Composite Controller

```sh
# take a look
cat composite-controllers.yaml

# Apply
kubectl apply -f composite-controllers.yaml

# Verify
kubectl get compositecontrollers
```

At this point we have an app and Metacontroller "sees" it. The composite created is listening for events such as `kind: App` created. In response it calls a web hook with the desired state. So next we need to implement the webhooks. 

sync hook failed for App production/my-app

## Webhook Deployments (aka controller)

We need to create a deployment that can handle the incoming webhook and implement the reconciliation between the desired and actual state. This is our controller logic.

### Slack Image

Before we get started we need to create a container image for our slack deployment

```sh
# Use ttl.sh anonymous and ephemeral image registry for simplicity
# Guess I am to lazy to stand up a registry :-)
UUID=$(uuidgen)
IMAGE_NAME="${UUID,,}" # convert to lower case
REG="ttl.sh/${IMAGE_NAME}:2h"
podman build -t "$REG" .
podman push "$REG"

# Update the Deployment spec to point at our ephemeral image
reg="$REG" yq --inplace 'select(.kind == "Deployment") | .spec.template.spec.containers[0].image = env(reg)' controller-slack-deploy.yaml 

```

### Config Map

The app controller is uses a python app to service the webhook. We will store the python app in a config map so that we can mount and run it in our pod.

```sh
# Store python app in a config map
kubectl -n controllers create configmap app --from-file=app.py
```

### Deployments

We are now ready to apply the `Deployment` and `Service` that will handle our webhooks.

```sh
kubectl apply -f controller-app.yaml
kubectl apply -f controller-slack-deploy.yaml 
kubectl apply -f controller-slack-svc.yaml 
```

## Verify

```sh
# DEVTODO this might not be needed, you need to test and figure out if it is needed or not
# Recreate the app just to trigger events
kubectl delete -f my-app.yaml
kubectl apply -f my-app.yaml
```

If all went well the the controller for your app received the "app created" event, called the webhook and created the app. To verify let's hit the app.

```sh
# DEVTODO setup an ingress controller but until then we can use port-forward
kubectl -n production port-forward svc/my-app 8080:80

open http://localhost:8080/
```
