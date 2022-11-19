#!/bin/bash

SERVER_NAME="serverNameExample"
DEPLOY_FILE="deployments/kubernetes/${SERVER_NAME}-deployment.yml"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# Determining whether a file exists
if [ ! -f "${DEPLOY_FILE}" ];then
  echo "Deployment file file ${DEPLOY_FILE} does not exist"
  checkResult 1
fi

# Check if you are authorised to operate k8s
echo "kubectl version"
kubectl version
checkResult $?

echo "kubectl delete -f ${DEPLOY_FILE} --ignore-not-found"
kubectl delete -f ${DEPLOY_FILE} --ignore-not-found
checkResult $?

sleep 1

echo "kubectl apply -f ${DEPLOY_FILE}"
kubectl apply -f ${DEPLOY_FILE}
