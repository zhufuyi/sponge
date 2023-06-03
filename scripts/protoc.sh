#!/bin/bash

# the directory where the proto files are located
protoBasePath="api"
allProtoFiles=""

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# add the import of useless packages from the generated *.pb.go code here
function deleteUnusedPkg() {
  file=$1
  sed -i "s#_ \"github.com/envoyproxy/protoc-gen-validate/validate\"##g" ${file}
  sed -i "s#_ \"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options\"##g" ${file}
  sed -i "s#_ \"github.com/srikrsna/protoc-gen-gotag/tagger\"##g" ${file}
  sed -i "s#_ \"google.golang.org/genproto/googleapis/api/annotations\"##g" ${file}
}

function listProtoFiles(){
    cd $1
    items=$(ls)

    for item in $items; do
        if [ -d "$item" ]; then
            listProtoFiles $item
        else
            if [ "${item#*.}"x = "proto"x ];then
              file=$(pwd)/${item}
              protoFile="${protoBasePath}${file#*${protoBasePath}}"
              allProtoFiles="${allProtoFiles} ${protoFile}"
            fi
        fi
    done
    cd ..
}

function handlePbGoFiles(){
    cd $1
    items=$(ls)

    for item in $items; do
        if [ -d "$item" ]; then
            handlePbGoFiles $item
        else
            if [ "${item#*.}"x = "pb.go"x ];then
              deleteUnusedPkg $item
            fi
        fi
    done
    cd ..
}

function generateByAllProto(){
  # get all proto file paths
  listProtoFiles $protoBasePath
  if [ "$allProtoFiles"x = x ];then
    echo "Error: not found protobuf file in path $protoBasePath"
    exit 1
  fi

  # generate files *_pb.go
  protoc --proto_path=. --proto_path=./third_party \
    --go_out=. --go_opt=paths=source_relative \
    $allProtoFiles

  checkResult $?
  # todo generate grpc files here
  # delete the templates code start
  # generate files *_grpc_pb.go
  protoc --proto_path=. --proto_path=./third_party \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $allProtoFiles

  checkResult $?
  # delete the templates code end

  # generate the file *_pb.validate.go
  protoc --proto_path=. --proto_path=./third_party \
    --validate_out=lang=go:. --validate_opt=paths=source_relative \
    $allProtoFiles

  checkResult $?

  # embed the tag field into *_pb.go
  protoc --proto_path=. --proto_path=./third_party \
    --gotag_out=:. --gotag_opt=paths=source_relative \
    $allProtoFiles

  checkResult $?
}

function generateBySpecifiedProto(){
  # get the proto file of the serverNameExample server
  allProtoFiles=""
  listProtoFiles ${protoBasePath}/serverNameExample
  cd ..
  specifiedProtoFiles=$allProtoFiles
  # todo generate router code for gin here
  # delete the templates code start 2

  # generate the swagger document and merge all files into docs/apis.swagger.json
  protoc --proto_path=. --proto_path=./third_party \
    --openapiv2_out=. --openapiv2_opt=logtostderr=true --openapiv2_opt=allow_merge=true --openapiv2_opt=merge_file_name=docs/apis.json \
    $specifiedProtoFiles

  checkResult $?

  sponge web swagger --file=docs/apis.swagger.json
  checkResult $?

  echo ""
  echo "run server and see docs by http://localhost:8080/apis/swagger/index.html"
  echo ""

  # A total of four files are generated: the registration route file **router.pb.go (saved in the same directory as the protobuf file),
  # the injection route file *_service.pb.go (saved in internal/routers by default), the logic code template file *_logic.go (saved in internal/service by default),
  # and the return error code template file *_http.go (saved in internal/ecode by default). internal/service),
  # return error code template file *_http.go (default path in internal/ecode)
  protoc --proto_path=. --proto_path=./third_party \
    --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
    --go-gin_opt=moduleName=github.com/zhufuyi/sponge --go-gin_opt=serverName=serverNameExample \
    $specifiedProtoFiles

  checkResult $?
  # delete the templates code end 2
}

# generate pb.go by all proto files
generateByAllProto

# generate pb.go by specified proto files
generateBySpecifiedProto

# delete unused packages in pb.go
handlePbGoFiles $protoBasePath

go mod tidy
checkResult $?

echo "exec protoc command successfully."
