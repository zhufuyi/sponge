部署服务到k8s前，在已经登录镜像仓库的docker主机中，创建一个为k8s拉取镜像权限的Secret，命令如下：

```bash
kubectl create secret generic docker-auth-secret \
    --from-file=.dockerconfigjson=/root/.docker/config.json> \
    --type=kubernetes.io/dockerconfigjson
```

<br>

启动服务：

> kubectl apply -f ./

查看启动状态：

> kubectl get all -n project-name-example

<br>

简单测试http端口

```bash
# 在本机端口映射到服务的http端口
kubectl port-forward --address=0.0.0.0 service/server-name-example-svc 8080:8080 -n project-name-example
```
