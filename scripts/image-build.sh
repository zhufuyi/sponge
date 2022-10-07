#!/bin/bash

# 直接复制已编译的二进制文件构建出来的镜像。 优点：构建速度快，缺点：镜像体积被两阶段构建大一倍。

serverName="serverNameExample"
# 服务的镜像名称，不能有大写字母
SERVER_NAME="project-name-example.server-name-example"
# Dockerfile文件目录
DOCKERFILE_PATH="build"
DOCKERFILE="${DOCKERFILE_PATH}/Dockerfile"

# 镜像仓库地址，REPO_HOST="ip或域名"，通过第一个参数传进来
REPO_HOST=$1
if [ "X${REPO_HOST}" = "X" ];then
        echo "param 'repo host' cannot be empty, example: ./image-build.sh hub.docker.com v1.0.0"
        exit 1
fi
# 版本tag，如果为空，默认为latest，通过第二个参数传进来
TAG=$2
if [ "X${TAG}" = "X" ];then
        TAG="latest"
fi
# 镜像名称和tag
IMAGE_NAME_TAG="${REPO_HOST}/${SERVER_NAME}:${TAG}"

# 二进制执行文件
BIN_FILE="cmd/${serverName}/${serverName}"
# 配置文件目录
CONFIG_PATH="configs"

# only grpc use start
bash scripts/grpc_health_probe.sh
# only grpc use end
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn,direct go build -gcflags "all=-N -l" -o ${BIN_FILE} cmd/${serverName}/*.go
mv -f ${BIN_FILE} ${DOCKERFILE_PATH}
mkdir -p ${DOCKERFILE_PATH}/${CONFIG_PATH} && cp -f ${CONFIG_PATH}/${serverName}.yml ${DOCKERFILE_PATH}/${CONFIG_PATH}
echo "docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}
rm -rf ./${DOCKERFILE_PATH}/${serverName} ${DOCKERFILE_PATH}/configs ${DOCKERFILE_PATH}/grpc_health_probe
