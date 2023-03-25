#!/bin/bash

# build the docker image using the binaries, if you want to reduce the size of the image,
# use upx to compress the binaries before building the image.

serverName="serverNameExample_mixExample"
# image name of the service, no capital letters
SERVER_NAME="project-name-example.server-name-example"
# Dockerfile file directory
DOCKERFILE_PATH="build"
DOCKERFILE="${DOCKERFILE_PATH}/Dockerfile"

# image repo address, REPO_HOST="ip or domain", passed in via the first parameter
REPO_HOST=$1
if [ "X${REPO_HOST}" = "X" ];then
        echo "param 'repo host' cannot be empty, example: ./image-build.sh hub.docker.com v1.0.0"
        exit 1
fi
# the version tag, which defaults to latest if empty, is passed in via the second parameter
TAG=$2
if [ "X${TAG}" = "X" ];then
        TAG="latest"
fi
# image name and tag
IMAGE_NAME_TAG="${REPO_HOST}/${SERVER_NAME}:${TAG}"

# binary executable files
BIN_FILE="cmd/${serverName}/${serverName}"
# configuration file directory
CONFIG_PATH="configs"

# only grpc use start
bash scripts/grpc_health_probe.sh
# only grpc use end
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct go build -o ${BIN_FILE} cmd/${serverName}/*.go
mv -f ${BIN_FILE} ${DOCKERFILE_PATH}
mkdir -p ${DOCKERFILE_PATH}/${CONFIG_PATH} && cp -f ${CONFIG_PATH}/${serverName}.yml ${DOCKERFILE_PATH}/${CONFIG_PATH}

# compressing binary file
#cd ${DOCKERFILE_PATH}
#upx -9 ${serverName}
#upx -9 grpc_health_probe
#cd -

echo "docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}

if [ X"${serverName}" = X ];then
        exit 0
fi
rm -rf ./${DOCKERFILE_PATH}/${serverName} ${DOCKERFILE_PATH}/configs ${DOCKERFILE_PATH}/grpc_health_probe
