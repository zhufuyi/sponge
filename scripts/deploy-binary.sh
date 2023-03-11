#!/bin/bash

# chkconfig: - 85 15
# description: serverNameExample

serverName="serverNameExample_mixExample"
cmdStr="cmd/${serverName}/${serverName}"


function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

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

    if [ -f "${NAME}" ] ;then
        rm "${NAME}"
    fi
    sleep 0.2
    go build -o ${cmdStr} cmd/${NAME}/main.go
    checkResult $?

    nohup ${cmdStr} > ${NAME}.log 2>&1 &
    sleep 1

    ID=`ps -ef | grep "$NAME" | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
    if [ -n "$ID" ]; then
        echo "Start the ${NAME} service successfully, process ID=${ID}"
    else
        echo "Failed to start ${NAME} service"
		    return 1
    fi
	return 0
}


stopService ${serverName}
if [ "$1"x != "stop"x ] ;then
  sleep 1
  startService ${serverName}
  checkResult $?
else
  echo "Service ${serverName} has stopped"
fi
