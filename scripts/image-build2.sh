#!/bin/bash

# 两阶段构建镜像，优点：镜像体积最小，缺点：构建速度较慢，每次构建都产生比较大的中间镜像
serverName="serverNameExample"
# 服务的镜像名称，不能有大写字母
SERVER_NAME="project-name-example.server-name-example"
# Dockerfile文件目录
DOCKERFILE_PATH="build"
DOCKERFILE="${DOCKERFILE_PATH}/Dockerfile_build"

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

PROJECT_FILES=$(ls)
tar zcf ${serverName}.tar.gz ${PROJECT_FILES}
mv -f ${serverName}.tar.gz ${DOCKERFILE_PATH}
echo "docker build --force-rm -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}"
docker build --force-rm -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}
rm -rf ${DOCKERFILE_PATH}/${serverName}.tar.gz
