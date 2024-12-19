### 1 Build jenkins-go

Build image:

> docker build -t zhufuyi/jenkins-go:2.37 .

Push to image repository.

```shell
# login to the image repository, if it is a private image repository, specify the address
docker login -u username

docker push zhufuyi/jenkins-go:2.37
```

<br>

### 2 run jenkins-go

start the jenkins server

> docker-compose up -d

Visit `http://<address>:38080` in your browser, the first boot requires the initial admin password (obtained via the command `docker exec jenkins-go cat /var/jenkins_home/secrets/initialAdminPassword`).

<br>

### 3 Configure jenkins

After logging in to jenkins, you need to install plugins. If you are not sure which plugins you need for now, click Install Recommended Plugins.

**Create administrator account**

example: admin 123456

<br>

**Install the required plugins**

Click [Manage Jenkins] --> [Manage Plugins] -->  [Available Plugins], Install the following plugins.

```bash
# if required, install Chinese plugin
Locale

# adding parametric build plugin
Extended Choice Parameter

# adding the git parameters plugin
Git Parameter

# account Management
Role-based Authorization Strategy
```

Restart the jenkins service to enable the plugin.

<br>

**Configure global parameters**

dashboard --> System Management -->System Configuration -->Check Environment Variables

set the image's repository address

```bash
# add environment variables
PATH
/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/go/bin

# development environment image repository
DEV_REPO_HOST
dev.host.docker.com

# test environment image repository
TEST_REPO_HOST
test.host.docker.com

# production environment image repository
PROD_REPO_HOST
prod.host.docker.com
```

<br><br>

### 4 Deploy services using jenkins

A relatively simple way to configure Jenkins tasks: import existing templates when creating new tasks (fill in the name of an existing task), and then modify the git repository address. If you don't have a template, create a new task as follows.

**(1) Create a new jenkins job**

![create job](https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/createJob.jpg)

<br>

**(2)  Set up a parametric build with the parameter name `GIT_PARAMETER`**, as shown below.

![parametric construction](https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/paramSetting.jpg)

<br>

**(3) Set up pipeline information**, as shown below.

![flow line](https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/pipelineSetting.jpg)

<br>

**(4) Build the project**, click Build with Parameters on the left menu bar, and select the corresponding parameters to build, as shown in the following figure.

![runJob-dev](https://raw.githubusercontent.com/go-dev-frame/sponge/main/assets/building.jpg)

Note: Before building, modify the pinned or email notification target to facilitate viewing the build deployment results, open the Jenkinsfile file under the code repository, find the field tel_num, and fill in the cell phone number.

<br><br>

### 5 Authorization settings for image repositories

Before executing the image-push.sh script, the jenkins-go container must first be authorized by the image repository, which may not be the same for different image repositories.

**Private docker image repository authorization**

```bash
docker login <ip:port>
# account
# password
```

<br>

**harbor image repository authorization**

```bash
# (1) docker login harbor
docker login <ip:port>
# account
# password

# (2) if harbor uses private http certificates, put the license key in the certs.d file of docker, for example, the file path is as follows.
/etc/docker/certs.d/<ip>/<ip>.crt
```

<br><br>

### 6 Authorization settings for pulling image repositories

When pulling images from k8s deployment service requires authorization, you need to create an additional `Secret` for login when pulling images.

Way 1: Create Secret directly

```bash
kubectl create secret docker-registry docker-auth-secret \
    --docker-server=DOCKER_REGISTRY_SERVER \
    --docker-username=DOCKER_USER \
    --docker-password=DOCKER_PASSWORD \
    --docker-email=DOCKER_EMAIL
```

Way 2: Create a Secret in the docker host that is already logged into the image repository (recommended)

```bash
kubectl create secret generic docker-auth-secret \
    --from-file=.dockerconfigjson=/root/.docker/config.json> \
    --type=kubernetes.io/dockerconfigjson
```

<br>

Deployment, pod's resource configuration `imagePullSecrets` to specify the key

```yaml
# ......
    spec:
      containers:
        - name: server-name-example
          image: project-name-example/server-name-example:latest
# ......
      imagePullSecrets:
        - name: docker-auth-secret
```

<br>

Note: If harbor requires a private https certificate, you need to put /etc/docker/certs.d/<ip>/<ip>.crt into the same directory as the nodes (master and node) of k8s, but not if it is a public certificate.
