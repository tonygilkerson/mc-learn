# Hello World

## Setup

```sh
# Create a namespace
kubectl create ns hello
kubectl ns hello           # IMPORTANT dont skip this

# Apply a custom resource
kubectl apply -f hello-world/crd.yaml

# Apply a custom controller
kubectl apply -f hello-world/controller.yaml

# Apply the webhook
kubectl apply -f hello-world/webhook.yaml

# Create a config map that contains our Python script
kubectl create configmap hello-controller --from-file=hello-world/sync.py
```

## Deploy

Edit the `hello.yaml` and enter your name, then create an instance of the "hello" custom resource.

```sh
kubectl apply -f hello-world/hello.yaml
```

## Results

It might take awhile for the Python image to download but after awhile all should be successful. Look around.

The first thing we see it the `tony` pod created by our hello-controller. It is called `tony` because that is the resource name set in the `hello.yaml` file, your pod will have a different name if you edited this file.

```sh
$ kubectl get pods
NAME                                READY   STATUS      RESTARTS   AGE
hello-controller-7bcfc6ff69-mvs7k   1/1     Running     0          14m
tony                                0/1     Completed   0          7m8s
```

If we describe our pod we can see its parent resource. If we look at the parent resource we see it contains our desired state.

```sh
# Find the pods parent
$ kubectl describe pod/tony | grep "Controlled By"
Controlled By:  HelloWorld/tony

# List parents
$ kubectl get helloworlds
NAME   AGE
tony   22m

# Look at our parent
$ kubectl get HelloWorld/tony -oyaml
apiVersion: example.com/v1
kind: HelloWorld
metadata:
  name: tony
  namespace: hello
spec:
  who: Tony Gilkerson
```

Looking at the logs we see it printed the name as expected.

```sh
$ kubectl get helloworlds
NAME   AGE
tony   7m56s
```

If we delete our pod we will notice it comes right back because its parent still exists and want to bring about the desired state.

```sh
$ kubectl delete pod/tony; sleep 1;  kubectl get pods 
pod "tony" deleted
NAME                                READY   STATUS              RESTARTS   AGE
hello-controller-7bcfc6ff69-mvs7k   1/1     Running             0          34m
tony                                0/1     ContainerCreating   0          1s
```

If we delete the parent the pod goes away and everything is cleaned up nicely. The controller still exists, we will delete that in the next section.

```sh
$ kubectl delete HelloWorld/tony; sleep 1;  kubectl get pods 
helloworld.example.com "tony" deleted
NAME                                READY   STATUS    RESTARTS   AGE
hello-controller-7bcfc6ff69-mvs7k   1/1     Running   0          36m
```

## Cleanup

```sh
kubectl delete cm hello-controller
kubectl delete -f hello-world
kubectl delete ns hello
```
