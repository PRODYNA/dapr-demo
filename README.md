# demo
DAPR Demo application

Configure and lauch minikube

```
minikube config set memory 8000
minikube config set cpus 4
minikube config set kubernetes-version 1.23.3
minikube start
```

Once minikube is running enable the ingress controller

```
minikube addons enable ingress-controller
```

Read out the IP address of minikube with

```
minikube ip
```

and write those entries to your /etc/hosts

```
<minikube-ip> dapr.minikube servicea.minikube serviceb.minikube servicec.minikube
```

In the subdirectory terraform run terraform

```
terraform apply
```

This will install 

* DAPR in the namespace dapr-system
* Service A, B and C in the namespace Backend

