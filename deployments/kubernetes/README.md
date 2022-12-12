Before deploying the service to k8s, create a Secret that pulls image permissions for k8s in a docker host that is already logged into the image repository, with the following command.

```bash
kubectl create secret generic docker-auth-secret \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```

<br>

run server:

```bash
cd deployments

kubectl apply -f ./*namespace.yml

kubectl apply -f ./
```

view the start-up status.

> kubectl get all -n project-name-example

<br>

simple test of http port

```bash
# mapping to the http port of the service on the local port
kubectl port-forward --address=0.0.0.0 service/<server-name-example-svc> 8080:8080 -n <project-name-example>
```
