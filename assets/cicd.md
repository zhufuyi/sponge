
## 6 Continuous Integration Deployment

The services created by sponge support build and deployment in [jenkins](https://www.jenkins.io/doc/), the deployment target can be docker, [k8s](https://kubernetes.io/docs/home/) , the deployment script is in the **deployments** directory, the following is an example of deployment to k8s using jenkins.

### 6.1 Building the jenkins-go platform

In order to be able to compile go code in a container, you need to build a jenkins-go image, which is already built [jenkins-go image](https://hub.docker.com/r/zhufuyi/jenkins-go/tags). If you want to build the jenkins-go image yourself, you can refer to the docker build script [Dokerfile](https://github.com/zhufuyi/sponge/blob/main/test/server/jenkins/Dockerfile)

After preparing the jenkins-go image, you also need to prepare a k8s cluster (there are many k8s cluster tutorials online), a k8s forensics file and a command line tool [kubectl](https://kubernetes.io/zh-cn/docs/tasks/tools/#kubectl) to ensure that you have permission to operate k8s in the jenkins-go container .

The jenkins-go startup script, docker-compose.yml, reads

```yaml
version: "3.7"
services:
  jenkins-go:
    image: zhufuyi/jenkins-go:2.37
    restart: always
    container_name: "jenkins-go"
    ports:
      - 38080:8080
    #- 50000:50000
    volumes:
      - $PWD/jenkins-volume:/var/jenkins_home
      # docker configuration
      - /var/run/docker.sock:/var/run/docker.sock
      - /usr/bin/docker:/usr/bin/docker
      - /root/.docker/:/root/.docker/
      # k8s api configuration directory, including config file
      - /usr/local/bin/kubectl:/usr/local/bin/kubectl
      - /root/.kube/:/root/.kube/
      # go related tools
      - /opt/go/bin/golangci-lint:/usr/local/bin/golangci-lint
```

Start the jenkis-go service.

> docker-compose up -d

Visit [http://localhost:38080](http://localhost:38080) in your browser, the first time you start it you need the admin key (execute the command to get `docker exec jenkins-go cat /var/jenkins_home/secrets/initialAdminPassword`), then install the recommended plugins and set the admin account password, then install some plugins you need to use and some custom settings.

**(1) Installation of plug-ins**

```bash
# Chinese plugin
Locale

# Add parametric build plugins
Extended Choice Parameter

# Add git parameters plugin
Git Parameter

# Account management
Role-based Authorization Strategy
```

**(2) Set Chinese **


Click [Manage Jenkins] -> [Configure System] option, find the [Locale] option, enter [zh_CN], check the following options, and finally click [Apply].


**(3) Configuration of global parameters**

dashboard --> System Administration --> System Configuration --> Check the environment variables

Set the repository address of the container image.

```bash
# Development environment image repository
DEV_REGISTRY_HOST http://localhost:27070

# Test environment image repository
TEST_REGISTRY_HOST http://localhost:28080

# Production environment image repository
PROD_REGISTRY_HOST http://localhost:29090
```

<br>

### 6.2 Creating templates

A relatively simple way to create a new task for jenkins is to import an existing template when creating a new task and then modify the git repository address. The first time you use jenkins and don't have a template yet, you can create one by following these steps.

**(1) Create a new task**, as shown in Figure 6-1.

![create job](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/createJob.jpg)
* Figure 6-1 Task creation screen*

<br>

**(2) Parameterized configuration setting**, using the parameter name `GIT_parameter`, as shown in Figure 6-2.

![parametric construction](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/paramSetting.jpg)
*Figure 6-2 Setting up the parametric build interface*

<br>

**(3) Set up the pipeline**, as shown in Figure 6-3.

![flow line](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/pipelineSetting.jpg)
*Figure 6-3 Setup pipeline screen*

<br>

**(4) Construction project**

Click **Build with Parameters** on the left menu bar, and then select the branch or tag you want to branch or tag, as shown in Figure 6-4.

![run job](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/building.jpg)
* Figure 6-4 Parametric build interface*

<br>

### 6.3 Deploying to k8s

Take the edusys service in **Chapter 3.1.2** as an example, built and deployed to k8s using jenkins.

The first build of the service requires some prep work.

(1) Upload the edusys code to the code repository.

(2) Prepare a docker image repository and make sure the docker where jenkins-go is located has permission to upload images to the image repository.

(3) Ensure that you have permission to pull images from the mirror on the k8s cluster node, and execute the command to generate the key on the logged-in docker image repository server.

```bash
kubectl create secret generic docker-auth-secret \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```

(4) Create edusys related resources at k8s.

```bash
# Switch to directory
cd deployments/kubernetes

# Create namespace, name corresponds to spong create service parameter project-name
kubectl apply -f . /*namespace.yml

# Create configmap, service
kubectl apply -f . /*configmap.yml
kubectl apply -f . /*svc.yml
```

(5) If you want to use pinned notifications to view the build deployment results, open the **Jenkinsfile** file under the code base, find the field **tel_num** and fill in the mobile number, and find **access_token** and fill in the token value.

<br>

After the prep, create a new task (name edusys) in the jenkins interface, using the template created above (name sponge), then modify the git repository, save the task and start the parametric build, the result of the construction is shown in Figure 6-5.

![run job](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/jenkins-build.jpg)
*Figure 6-5 jenkins build result interface*

<br>

Use the command `kubectl get all -n edusys` to see the status of the edusys service running in k8s.

```
NAME                             READY   STATUS    RESTARTS   AGE
pod/edusys-dm-77b4bcccc5-8xt8v   1/1     Running   0          21m

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/edusys-svc   ClusterIP   10.108.31.220   <none>        8080/TCP   27m

NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/edusys-dm   1/1     1            1           21m

NAME                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/edusys-dm-77b4bcccc5   1         1         1       21m
```

<br>

Test locally to see if it is accessible

```bash
# Proxy ports
kubectl port-forward --address=0.0.0.0 service/edusys-svc 8080:8080 -n edusys

# Requests
curl http://localhost:8080/api/v1/teacher/1
```

<br>

The services generated by sponge include a Jenkinsfile, build and upload image scripts, and k8s deployment scripts, which can be used basically without modifying the scripts, or you can modify the scripts to suit your scenario.

<br><br>
