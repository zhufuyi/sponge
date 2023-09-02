#!/bin/bash

# image name, prohibit uppercase letters in names.
IMAGE_NAME="project-name-example/server-name-example"

# image repo address, passed in via the first parameter
REPO_HOST=$1
if [ "X${REPO_HOST}" = "X" ];then
    echo "param 'repo host' cannot be empty, example: ./image-push.sh hub.docker.com v1.0.0"
    exit 1
fi

# version tag, passed in via the second parameter, if empty, defaults to latest
TAG=$2
if [ "X${TAG}" = "X" ];then
    TAG="latest"
fi
# image name and tag
IMAGE_NAME_TAG="${REPO_HOST}/${IMAGE_NAME}:${TAG}"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# image repository host, https://index.docker.io/v1 is the official docker image repository
IMAGE_REPO_HOST="image-repo-host"
# check if you are authorized to log into docker
function checkLogin() {
  loginStatus=$(cat /root/.docker/config.json | grep "${IMAGE_REPO_HOST}")
  if [ "X${loginStatus}" = "X" ];then
      echo "docker is not logged into the image repository"
      checkResult 1
  fi
}

checkLogin

# push image to image repository
echo "docker push ${IMAGE_NAME_TAG}"
docker push ${IMAGE_NAME_TAG}
checkResult $?
echo "docker push image success."

sleep 1

# delete image
echo "docker rmi -f ${IMAGE_NAME_TAG}"
docker rmi -f ${IMAGE_NAME_TAG}
checkResult $?
echo "docker remove image success."
