#!/bin/bash

serverName="serverNameExample"
# 镜像名称，不能有大写字母
SERVER_NAME="server-name-example"
PROJECT_NAME="project-name-example"
# 二进制执行程序
BIN_FILE="./cmd/${serverName}/${serverName}"
# Dockerfile文件位置
DOCKERFILE_PATH="build"
# 配置文件位置
CONFIG_PATH="configs"

# 私人仓库地址，通过第一个参数传进来，可以灵活构建成目标仓库地址的镜像
#REPO_HOST="ip或域名"
REPO_HOST=$1
if [ "X${REPO_HOST}" = "X" ];then
        echo "param 'repo host' cannot be empty, example: ./image-build.sh github.com v1.0.0"
        exit 1
fi

# 版本tag，通过第二个参数传进来，如果为空，默认为latest
TAG=$2
if [ "X${TAG}" = "X" ];then
        TAG="latest"
fi

# todo 根据服务类型(http或grpc)是否使用健康检查工具
bash scripts/grpc_health_probe.sh
cp -f /tmp/grpc_health_probe ${DOCKERFILE_PATH}

mv -f ${BIN_FILE} ${DOCKERFILE_PATH}
mkdir -p ${DOCKERFILE_PATH}/configs
cp -f ${CONFIG_PATH}/${serverName}.yml ${DOCKERFILE_PATH}/configs

echo "docker build -t ${REPO_HOST}/${PROJECT_NAME}/${SERVER_NAME}:${TAG} ${DOCKERFILE_PATH}"
docker build -t ${REPO_HOST}/${PROJECT_NAME}/${SERVER_NAME}:${TAG} ${DOCKERFILE_PATH}

rm -rf ./${DOCKERFILE_PATH}/${BIN_FILE} ${DOCKERFILE_PATH}/configs ${DOCKERFILE_PATH}/grpc_health_probe
