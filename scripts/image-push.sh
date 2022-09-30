#!/bin/bash

# todo 填写自己的私有镜像仓库host，https://index.docker.io/v1是docker官方镜像仓库
IMAGE_REPO_HOST="https://index.docker.io/v1"

# 镜像名称不能有大写字母
PROJECT_NAME="project-name-example"
SERVER_NAME="server-name-example"

# 私人仓库地址，通过第一个参数传进来
#REPO_HOST="ip或域名"
REPO_HOST=$1
if [ "X${REPO_HOST}" = "X" ];then
    echo "param 'repo host' cannot be empty, example: ./image-push.sh github.com v1.0.0"
    exit 1
fi

# 版本tag，通过第二参数传进来，如果为空，默认为latest
TAG=$2
if [ "X${TAG}" = "X" ];then
    TAG="latest"
fi

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}


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
docker push ${REPO_HOST}/$PROJECT_NAME/$SERVER_NAME:${TAG}
checkResult $?
echo "docker push ${REPO_HOST}/$PROJECT_NAME/$SERVER_NAME:${TAG} success."

sleep 1

# 删除镜像
docker rmi -f ${REPO_HOST}/$PROJECT_NAME/$SERVER_NAME:${TAG}
checkResult $?
echo "docker rmi -f ${REPO_HOST}/$PROJECT_NAME/$SERVER_NAME:${TAG} success."
