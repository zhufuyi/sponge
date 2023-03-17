### Install sponge in linux or macOS

#### (1) Install go, requires version 1.16 or above

Download go address: [https://go.dev/dl/](https://go.dev/dl/)

Check the go version after installation

```bash
go version
```

<br>

#### (2) Install protoc, requires v3.20 or above

Download protoc at: [https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3)

Add the protoc binary file to the system path.

Check the protoc version after installation

```bash
protoc --version
```

<br>

#### (3) Install sponge

```bash
# Install sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# Initialize sponge
sponge init

# Check if the plugins are installed successfully, if not, execute the command to retry sponge tools --install
sponge tools

# Check sponge version after installation
sponge -v
```

<br>
<br>
<br>

### Install sponge in windows

#### (1) Install go, requires version 1.16 or above

Download go address: [https://go.dev/dl/](https://go.dev/dl/)

Check the go version after installation

```bash
go version
```

<br>

#### (2) Install protoc, v3.20 or higher

Download protoc at: [https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3)

Add the protoc binary file to the system path.

Check the protoc version after installation

```bash
protoc --version
```

<br>

#### (3) Install linux command environment on windows

**install mingw64**

Download mingw64 at: [https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z](https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/8.1.0/threads-posix/seh/x86_64-8.1.0-release-posix-seh-rt_v6-rev0.7z)

After downloading and extracting to the `D:\Program Files\mingw64` directory, modify the system environment variable PATH and add `D:\Program Files\mingw64\bin`.

<br>

**Install the make command**

Switch to the `D:\Program Files\mingw64\bin` directory, find the `mingw32-make.exe` executable file, copy it and rename it to `make.exe`.

Check the version after installation

```bash
make -v
```

<br>

**Install cmder**

Download cmder at: [https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip](https://github.com/cmderdev/cmder/releases/download/v1.3.20/cmder.zip)

After downloading and extracting to the `D:\Program Files\cmder` directory, modify the system environment variable PATH and add `D:\Program Files\cmder`.

Open the `Cmder.exe` terminal and check if common linux commands are supported.

```bash
ls --version
make --version
cp --version
chmod --version
rm --version
```
<br>

#### (4) Install sponge

Open a `cmder`(not the cmd that comes with windows) terminal and execute the command to install sponge.

```bash
# Install sponge
go install github.com/zhufuyi/sponge/cmd/sponge@latest

# Initialize sponge
sponge init

# Check if the plugins are installed successfully, if not, execute the command to retry sponge tools --install
sponge tools

# Check sponge version after installation
sponge -v
```
