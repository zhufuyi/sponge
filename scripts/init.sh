#!/bin/bash

goModFile="go.mod"
thirdPartyProtoDir="third_party"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

if [ ! -f "../$goModFile" ]; then
    sponge patch copy-go-mod
    checkResult $?
    mv -f go.mod ..
    mv -f go.sum ..
fi

if [ ! -d "../$thirdPartyProtoDir" ]; then
    sponge patch copy-third-party-proto
    checkResult $?
    mv -f $thirdPartyProtoDir ..
fi
