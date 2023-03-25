#!/bin/bash

# build the image for local docker, using the binaries, if you want to reduce the size of the image,
# use upx to compress the binaries before building the image.

serverName="serverNameExample_mixExample"
# image name of the service, no capital letters
SERVER_NAME="project-name-example.server-name-example"
# Dockerfile file directory
DOCKERFILE_PATH="build"
DOCKERFILE="${DOCKERFILE_PATH}/Dockerfile"

bash scripts/grpc_health_probe.sh
mv -f cmd/${serverName}/${serverName} ${DOCKERFILE_PATH}/${serverName}

# compressing binary file
#cd ${DOCKERFILE_PATH}
#upx -9 ${serverName}
#upx -9 grpc_health_probe
#cd -

mkdir -p ${DOCKERFILE_PATH}/configs && cp -f configs/${serverName}.yml ${DOCKERFILE_PATH}/configs/
echo "docker build -f ${DOCKERFILE} -t ${SERVER_NAME}:latest ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${SERVER_NAME}:latest ${DOCKERFILE_PATH}


if [ X"${serverName}" = X ];then
        exit 0
fi
rm -rf ./${DOCKERFILE_PATH}/${serverName} ${DOCKERFILE_PATH}/configs ${DOCKERFILE_PATH}/grpc_health_probe
