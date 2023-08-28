#!/bin/bash

HOST_ADDR=$1

# change host addr
if [ "X${HOST_ADDR}" = "X" ];then
  HOST_ADDR=$(cat cmd/serverNameExample_mixExample/main.go | grep "@host" | awk '{print $3}')
  HOST_ADDR=$(echo  ${HOST_ADDR} | cut -d ':' -f 1)
else
    sed -i "s/@host .*:8080/@host ${HOST_ADDR}:8080/g" cmd/serverNameExample_mixExample/main.go
fi

swag init -g cmd/serverNameExample_mixExample/main.go

colorStart='\e[1;36m'
underLine='\e[4m'
markEnd='\e[0m'

echo ""
echo -e "execute the command ${colorStart}make run${markEnd} and then visit ${underLine}http://${HOST_ADDR}:8080/swagger/index.html${markEnd} in your browser."
echo ""
