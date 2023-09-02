#!/bin/bash

dockerComposeFilePath="deployments/docker-compose"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

mkdir -p ${dockerComposeFilePath}/configs
if [ ! -f "${dockerComposeFilePath}/configs/serverNameExample.yml" ];then
  cp configs/serverNameExample.yml ${dockerComposeFilePath}/configs
fi

# shellcheck disable=SC2164
cd ${dockerComposeFilePath}

docker-compose down
checkResult $?

docker-compose up -d
checkResult $?

colorCyan='\e[1;36m'
highBright='\e[1m'
markEnd='\e[0m'

echo ""
echo -e "run service successfully, if you want to stop the service, go into the ${highBright}${dockerComposeFilePath}${markEnd} directory and execute the command ${colorCyan}docker-compose down${markEnd}."
echo ""
