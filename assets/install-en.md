
Recommended to use go version 1.22 or above, [https://go.dev/doc/install](https://go.dev/doc/install)

<br>

---

<br>

### Windows Environment

Make sure the go locale is installed before installing sponge, and add `GOBIN` to the system environment variable **path**. If it is already set, skip this step:

 ```bash 
     #Check if GOBIN directory exists 
     go env GOBIN 
    
     #If empty, GOBIN needs to be set (e.g. D:\go\bin), administrator privileges may be required 
     go env -w GOBIN=D:\go\bin 
     #Then add GOBIN directory to system path environment variable 
 ``` 

<br>

> Because sponge depends on some linux commands, git bash and make need to be installed in windows to support the linux command environment.

For installation convenience, sponge and its dependent programs have been packaged together, download address: [sponge-install.zip](https://drive.google.com/drive/folders/1T55lLXDBIQCnL5IQ-i1hWJovgLI2l0k1?usp=sharing)

After downloading the file:

1. Unzip the file, double-click **install.bat** to install, the installation process of git is always the default (if you have installed git, you can skip this step)

2. Right-click any folder (Show more options) and select **Open Git Bash here** to open the git bash terminal:

```bash
# Initialize sponge, automatically install sponge's dependency plugins
sponge init

# Check sponge version
sponge -v
```

Note: 

- When using sponge development, please use git bash terminal, do not use the system default cmd, otherwise there will be an error that cannot find the command.
- Do not open a terminal in the `GOBIN` directory (the directory where the sponge executable is located) to execute the command `sponge run`.

In addition to the above installation of sponge in windows, it also provides native installation, click to view **Installing Sponge** --> **Windows Environment** [installation documentation](https://go-sponge.com/quick-start?id=installing-sponge).

<br>

---

<br>

### Linux or macOS Environment

1. Add `GOBIN` to the system environment variable path, skip this step if already set.

```bash
# Open .bashrc file
vim ~/.bashrc

# Copy the following command to .bashrc file
export GOROOT="/opt/go"     # your go installation directory
export GOPATH=$HOME/go      # Set the directory where the "go get" command downloads third-party packages
export GOBIN=$GOPATH/bin    # Set the directory where the executable binaries are compiled by the "go install" command.
export PATH=$PATH:$GOBIN:$GOROOT/bin  # Add the GOBIN directory to the system environment variable path.

# Save .bashrc file, and make the settings take effect
source ~/.bashrc

# View the GOBIN directory, if the output is not empty, the setting is successful.
go env GOBIN
```

<br>

2. Install sponge and its dependent plugins into the `GOBIN` directory.

**✅ Install protoc**

Download protoc from: [https://github.com/protocolbuffers/protobuf/releases/tag/v25.2](https://github.com/protocolbuffers/protobuf/releases/tag/v25.2)

Download the corresponding **protoc** executable file according to the system type, and move the **protoc** executable file to the same directory as the **go** executable file.

```bash
# Install Sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# Initialize Sponge, automatically install Sponge's dependency plugins
sponge init

# Check if all plugins have been successfully installed. If any plugins fail to install, retry with the command: sponge plugins --install
sponge plugins

# Check Sponge version
sponge -v
```

> Note: Do not open the terminal in the `GOBIN` directory to execute the command `sponge run`.

<br>

---

<br>

### Docker Environment

> ⚠ Sponge UI service started by docker only supports code generation function. If you need to develop based on the generated service code, you also need to install Sponge and the required plugins locally according to the installation instructions above.

**Docker Run**

```bash
docker run -d --name sponge -p 24631:24631 zhufuyi/sponge:latest -a http://your_host_ip:24631
```

<br>

**Docker Compose**

The content of the `docker-compose.yaml` file is as follows:

```yaml
version: "3.7"

services:
  sponge:
    image: zhufuyi/sponge:latest
    container_name: sponge
    restart: always
    command: ["-a","http://your_host_ip:24631"]
    ports:
      - "24631:24631"
```

Start the service:

```bash
docker-compose up -d
```

After a successful Docker deployment, access `http://your_host_ip:24631` in your browser.
