#!/bin/bash

# 镜像名称不能有大写字母
SERVER_NAME="project-name-example.server-name-example"

# 镜像仓库地址，通过第一个参数传进来
#REPO_HOST="ip或域名"
REPO_HOST=$1
if [ "X${REPO_HOST}" = "X" ];then
    echo "param 'repo host' cannot be empty, example: ./image-push.sh hub.docker.com v1.0.0"
    exit 1
fi

# 版本tag，通过第二参数传进来，如果为空，默认为latest
TAG=$2
if [ "X${TAG}" = "X" ];then
    TAG="latest"
fi
# 镜像名称和tag
IMAGE_NAME_TAG="${REPO_HOST}/${SERVER_NAME}:${TAG}"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# 镜像仓库host，https://index.docker.io/v1是docker官方镜像仓库
IMAGE_REPO_HOST="image-repo-host"
# 检查是否授权登录docker
function checkLogin() {
  loginStatus=$(cat /root/.docker/config.json | grep "${IMAGE_REPO_HOST}")
  if [ "X${loginStatus}" = "X" ];then
      echo "docker未登录镜像仓库"
      checkResult 1
  fi
}

checkLogin

# 上传镜像
echo "docker push ${IMAGE_NAME_TAG}"
docker push ${IMAGE_NAME_TAG}
checkResult $?
echo "docker push image success."

sleep 1

# 删除镜像
echo "docker rmi -f ${IMAGE_NAME_TAG}"
docker rmi -f ${IMAGE_NAME_TAG}
checkResult $?
echo "docker remove image success."
