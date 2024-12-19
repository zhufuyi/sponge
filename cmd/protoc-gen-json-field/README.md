## protoc-gen-json-field

Generate json code based on proto files.

<br>

### Installation

#### Installation of dependency plugins

```bash
# install protoc in linux
mkdir -p protocDir \
  && curl -L -o protocDir/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-linux-x86_64.zip \
  && unzip protocDir/protoc.zip -d protocDir\
  && mv protocDir/bin/protoc protocDir/include/ $GOROOT/bin/ \
  && rm -rf protocDir
```

#### Install protoc-gen-json-field

> go install github.com/go-dev-frame/sponge/cmd/protoc-gen-json-field@latest

<br>

### Usage

```bash
# generate json file
protoc --proto_path=. --json-field_out=. --json-field_opt=paths=source_relative demo.proto
```
