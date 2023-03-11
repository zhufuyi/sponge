#!/bin/bash

serviceName="serverNameExample"
cmdStr="./${serviceName} -c configs/${serviceName}.yml"

chmod +x ./${serviceName}

stopService(){
    NAME=$1

    ID=`ps -ef | grep "$NAME" | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
    if [ -n "$ID" ]; then
        for id in $ID
        do
           kill -9 $id
           echo "Stopped ${NAME} service successfully, process ID=${ID}"
        done
    fi
}

startService() {
    NAME=$1

    nohup ${cmdStr} > ${serviceName}.log 2>&1 &
    sleep 1

    ID=`ps -ef | grep "$NAME" | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
    if [ -n "$ID" ]; then
        echo "Start the ${NAME} service ...... process ID=${ID}"
    else
        echo "Failed to start ${NAME} service"
            return 1
    fi
    return 0
}


stopService ${serviceName}
if [ "$1"x != "stop"x ] ;then
  sleep 1
  startService ${serviceName}
  exit $?
  echo ""
else
  echo "Service ${serviceName} has stopped"
fi
