### 安装依赖

> 安装sponge之前需要先安装`go`和`protoc`两个依赖。

**✅ 安装 go**

下载go地址： [https://studygolang.com/dl](https://studygolang.com/dl)

> 要求1.16以上版本。

查看go版本 `go version`

<br>

**✅ 安装 protoc**

下载protoc地址： [https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3)

> 要求v3.20以上版本，把 protoc 二进制文件所在目录添加到系统环境变量path。

查看protoc版本: `protoc --version`

<br>

安装完go和protoc之后，接下来安装sponge，支持在windows、mac、linux环境安装。

> 如果不能科学上网，安装sponge时，获取github的库会遇到超时失败问题，建议设置为国内代理，执行命令 **go env -w GOPROXY=https://goproxy.cn,direct**

<br>
<br>

### linux或macOS环境

把`GOPATH/bin`添加到系统path，如果已经设置过跳过此步骤。

```bash
# 查看是否设置了GOPATH
go env GOPATH

# 如果为空，设置GOPATH目录
go env -w GOPATH=你的目录(示例~/golang/package)

# 把GOPATH下的bin目录(示例~/golang/package/bin) 添加到系统path
echo 'export PATH=$PATH:你的GOPATH/bin目录' >> ~/.bashrc
source ~/.bashrc
```

<br>

执行命令安装sponge，sponge和依赖插件将安装到 `GOPATH/bin` 目录下。

```bash
# 安装sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# 初始化sponge，自动安装sponge依赖插件
sponge init

# 查看插件是否都安装成功，如果发现有插件没有安装成功，执行命令重试 sponge tools --install
sponge tools

# 查看sponge版本
sponge -v
```

<br>
<br>

### Windows环境

> 在windows环境中需要安装mingw64、make、cmder来支持linux命令环境才可以使用sponge。

**✅ 安装 mingw64**

下载mingw64地址： [x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z)

下载后解压到`D:\Program Files\mingw64`目录，把linux常用命令所在的目录`D:\Program Files\mingw64\bin`添加系统环境变量PATH。

<br>

**✅ 安装 make 命令**

切换到`D:\Program Files\mingw64\bin`目录，找到`mingw32-make.exe`可执行文件，复制并改名为`make.exe`。

查看make版本：`make -v`

<br>

**✅ 安装 cmder**

下载cmder地址： [cmder-v1.3.20.zip](https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip)

下载后解压到`D:\Program Files\cmder`目录下，并把目录`D:\Program Files\cmder`添加到系统环境变量path。

对cmder进行简单的配置：

- **解决输入命令时的空格问题**，打开cmder界面，按下组合键win+alt+p进入设置界面，在左上角搜索`Monospace`，取消勾选，保存退出。
- **配置右键启动cmder**，按下组合键`win+x`，再按字母`a`进入有管理权限的终端，执行命令`Cmder.exe /REGISTER ALL`。 随便在一个文件夹里按下鼠标右键，选择`Cmder Here`即可打开cmder界面。

> ⚠ 在windows环境使用sponge开发项目，为了避免找不到linux命令错误，请使用cmder，不要用系统自带的cmd终端、Goland和VS Code下的终端。

打开`cmder.exe`终端，检查是否支持常用的linux命令。

```bash
ls --version
make --version
cp --version
chmod --version
rm --version
sed --version
```

<br>

**✅ 安装 sponge**

打开`cmder.exe`终端(不是windows自带的cmd)。

把`GOPATH/bin`添加到系统path，如果已经设置过跳过此步骤。

```bash
# 查看是否设置了GOPATH
go env GOPATH

# 如果为空，设置GOPATH目录，有可能需要使用管理员权限执行此命令
go env -w GOPATH=你的目录(示例D:\golang\package)

# 把GOPATH下的bin目录(示例D:\golang\package\bin)添加到系统path
```

<br>

执行命令安装sponge，sponge和依赖插件将安装到`GOPATH/bin`目录下

```bash
# 安装sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# 初始化sponge，自动安装sponge依赖插件
sponge init

# 查看插件是否都安装成功，如果发现有插件没有安装成功，执行命令重试 sponge tools --install
sponge tools

# 查看sponge版本
sponge -v
```

<br>

### 在docker安装sponge

> ⚠ 使用docker安装的sponge只是sponge ui界面服务，如果需要在生成的服务代码基础上进行开发，还是需要根据上面的安装说明，在本地安装sponge和依赖插件的二进制文件。

**方式一：Docker启动**

```bash
docker run -d --name sponge -p 24631:24631 zhufuyi/sponge:latest -l -a http://你的宿主机ip:24631
```

<br>

**方式二：docker-compose启动**

docker-compose.yaml 文件内容如下：

```yaml
version: "3.7"

services:
  sponge:
    image: zhufuyi/sponge:latest
    container_name: sponge
    restart: always
    command: ["-l","-a","http://你的宿主机ip:24631"]
    ports:
      - "24631:24631"
```

```bash
# 启动服务
docker-compose up -d

```

在docker部署成功后，在浏览器访问 `http://你的宿主机ip:24631`。
