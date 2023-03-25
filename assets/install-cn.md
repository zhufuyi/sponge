### 在linux或macOS上安装sponge

#### (1) 安装go，要求1.16版本以上

下载go地址： [https://studygolang.com/dl](https://studygolang.com/dl)

安装完后查看go版本

```bash
go version
```

<br>

#### (2) 安装 protoc，要求v3.20以上版本

下载protoc地址： [https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3)

把 protoc 二进制文件添加到系统path下。

安装完后查看protoc版本

```bash
protoc --version
```

<br>

#### (3) 安装 sponge

```bash
# 安装sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# 初始化sponge
sponge init

# 查看插件是否都安装成功，如果有安装不成功，执行命令重试 sponge tools --install
sponge tools

# 安装完后查看sponge版本
sponge -v
```

💡  如果想更新到最新版本，执行命令 `sponge upgrade`。

<br>
<br>
<br>

### 在windows上安装sponge

#### (1) 安装go，要求1.16版本以上

下载go地址： [https://studygolang.com/dl](https://studygolang.com/dl)

安装完后查看go版本

```bash
go version
```

<br>

#### (2) 安装 protoc，v3.20以上版本

下载protoc地址： [https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3)

把 protoc 二进制文件添加到系统path下。

安装完后查看protoc版本

```bash
protoc --version
```

<br>

#### (3) 在windows上安装支持linux命令环境

**安装 mingw64**

下载mingw64地址： [https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z)

下载后解压到`D:\Program Files\mingw64`目录下，修改系统环境变量PATH，新增`D:\Program Files\mingw64\bin`。

<br>

**安装 make 命令**

切换到`D:\Program Files\mingw64\bin`目录，找到`mingw32-make.exe`可执行文件，复制并改名为`make.exe`。

安装完后查看版本

```bash
make -v
```

<br>

**安装 cmder**

下载cmder地址： [https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip](https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip)

下载后解压到`D:\Program Files\cmder`目录下，修改系统环境变量PATH，新增`D:\Program Files\cmder`。

打开`Cmder.exe`终端，检查是否支持常用的linux命令。
```bash
ls --version
make --version
cp --version
chmod --version
rm --version
```

<br>

#### (4) 安装 sponge

打开`cmder.exe`终端(不是windows自带的cmd)，执行命令安装sponge：

```bash
# 安装sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# 初始化sponge
sponge init

# 查看插件是否都安装成功，如果有安装不成功，执行命令重试 sponge tools --install
sponge tools

# 安装完后查看sponge版本
sponge -v
```

💡  如果想更新到最新版本，执行命令 `sponge upgrade`。
