#!/bin/bash

dockerComposeFilePath="deployments/docker-compose"

mkdir -p ${dockerComposeFilePath}/configs
cp configs/serverNameExample.yml ${dockerComposeFilePath}/configs

# shellcheck disable=SC2164
cd ${dockerComposeFilePath}

docker-compose up -d
