#!/bin/bash

serviceName="serverNameExample"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# determine if the startup service script run.sh exists
runFile="~/app/${serviceName}/run.sh"
if [ ! -f "$runFile" ]; then
  # if it does not exist, copy the entire directory
  mkdir -p ~/app
  cp -rf /tmp/${serviceName}-binary ~/app/
  checkResult $?
  rm -rf /tmp/${serviceName}-binary*
else
  # replace only the binary file if it exists
  cp -f ${serviceName}-binary/${serviceName} ~/app/${serviceName}-binary/${serviceName}
  checkResult $?
  rm -rf /tmp/${serviceName}-binary*
fi

# running service
cd ~/app/${serviceName}-binary
chmod +x run.sh
./run.sh
checkResult $?

echo "server directory is ~/app/${serviceName}-binary"
