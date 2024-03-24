
要求使用go 1.20以上版本： [https://studygolang.com/dl](https://studygolang.com/dl)

> 注：如果不能科学上网，获取github的库可能会遇到超时失败问题，建议设置为国内代理，执行命令 **go env -w GOPROXY=https://goproxy.cn,direct**

<br>

---

<br>

### Windows环境

> 因为sponge依赖一些linux命令，因此在windows环境中需要安装git bash、make来支持linux命令环境。

为了安装方便，已经把sponge及其依赖的程序打包在一起，下载地址(选择一个下载即可)：

- 百度云：[**sponge-install.zip**](https://pan.baidu.com/s/1fiTiMROkiIIzAdj2bk93CA?pwd=prys)。
- 蓝奏云：[**sponge安装文件**](https://wwm.lanzoue.com/b049fldpi) 密码:5rq9，共下载4个文件，安装前先看`安装说明.txt`文件。

下载文件后：

(1) 解压文件，双击 **install.bat** 进行安装，安装git过程一直默认即可(如果已经安装过git，可以跳过安装git这个步骤)。

(2) 在任意文件夹下右键(显示更多选项)，选择【Open Git Bash here】打开git bash终端：

```bash
# 初始化sponge，自动安装sponge依赖插件
sponge init

# 查看sponge版本
sponge -v
```

> 注： 使用sponge开发时，请使用git bash终端，不要使用系统默认的cmd，否则会出现找不到命令的错误。

在windows除了上面安装sponge方式，还提供了原生安装，点击查看【🏷安装 sponge】 --> 【windows环境】[安装文档](https://go-sponge.com/zh-cn/quick-start?id=%f0%9f%8f%b7%e5%ae%89%e8%a3%85-sponge)。

<br>

---

<br>

### Linux或MacOS环境

(1) 把`GOBIN`添加到系统环境变量**path**，如果已经设置过可以跳过此步骤。

```bash
# 打开 .bashrc 文件
vim ~/.bashrc

# 复制下面命令到.bashrc
export GOROOT="/opt/go"     # 你的go安装目录
export GOPATH=$HOME/go      # 设置 go get 命令下载第三方包的目录
export GOBIN=$GOPATH/bin    # 设置 go install 命令编译后生成可执行文件的存放目录
export PATH=$PATH:$GOBIN:$GOROOT/bin   # 把GOBIN目录添加到系统环境变量path

# 保存 .bashrc 文件后，使设置生效
source ~/.bashrc

# 查看GOBIN目录
go env GOBIN
```

<br>

(2) 把sponge及其依赖的插件安装到 `GOBIN` 目录。

**✅ 安装 protoc**

下载protoc地址： [https://github.com/protocolbuffers/protobuf/releases/tag/v25.2](https://github.com/protocolbuffers/protobuf/releases/tag/v25.2)

根据系统类型下载对应的 **protoc** 可执行文件，把 **protoc** 可执行文件移动到`GOBIN`目录下。

```bash
# 安装sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# 初始化sponge，自动安装sponge依赖插件
sponge init

# 查看插件是否都安装成功，如果发现有插件没有安装成功，执行命令重试 sponge plugins --install
sponge plugins

# 查看sponge版本
sponge -v
```

<br>

---

<br>

### Docker环境

> ⚠ 使用docker启动的sponge UI服务，只支持在界面操作来生成代码功能，如果需要在生成的服务代码基础上进行开发，还是需要根据上面的安装说明，在本地安装sponge和依赖插件。

**方式一：Docker启动**

```bash
docker run -d --name sponge -p 24631:24631 zhufuyi/sponge:latest -a http://你的宿主机ip:24631
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
    command: ["-a","http://你的宿主机ip:24631"]
    ports:
      - "24631:24631"
```

启动服务：

```bash
docker-compose up -d
```

在docker部署成功后，在浏览器访问 `http://你的宿主机ip:24631`。
