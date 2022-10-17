#!/bin/bash

HOST_ADDR=$1

if [ "X${HOST_ADDR}" = "X" ];then
  HOST_ADDR=$(cat cmd/serverNameExample/main.go | grep "@host" | awk '{print $3}')
  HOST_ADDR=$(echo  ${HOST_ADDR} | cut -d ':' -f 1)
else
    sed -i "s/@host .*:8080/@host ${HOST_ADDR}:8080/g" cmd/serverNameExample/main.go
fi

swag init -g cmd/serverNameExample/main.go

echo ""
echo "run server and see docs by http://${HOST_ADDR}:8080/swagger/index.html"
echo ""
