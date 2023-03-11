#!/bin/bash

serviceName="serverNameExample"

mkdir -p ${serviceName}-binary/configs

cp -f deployments/binary/run.sh ${serviceName}-binary
chmod +x ${serviceName}-binary/run.sh

cp -f deployments/binary/deploy.sh ${serviceName}-binary
chmod +x ${serviceName}-binary/deploy.sh

cp -f cmd/${serviceName}_mixExample/${serviceName}_mixExample ${serviceName}-binary
cp -f configs/${serviceName}.yml ${serviceName}-binary/configs
cp -f configs/${serviceName}_cc.yml ${serviceName}-binary/configs

tar zcvf ${serviceName}-binary.tar.gz ${serviceName}-binary
rm -rf ${serviceName}-binary

echo ""
echo "package binary successfully, output file = ${serviceName}-binary.tar.gz"
