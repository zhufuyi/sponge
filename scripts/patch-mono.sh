#!/bin/bash

goModFile="go.mod"
thirdPartyProtoDir="third_party"
genServerType=$1

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

if [ ! -f "../$goModFile" ]; then
    sponge patch copy-go-mod -f
    checkResult $?
    mv -f go.mod ..
    mv -f go.sum ..
fi

if [ "$genServerType"x != "http"x ]; then
    if [ ! -d "../$thirdPartyProtoDir" ]; then
        sponge patch copy-third-party-proto
        checkResult $?
        mv -f $thirdPartyProtoDir ..
    fi
fi

if [ "$genServerType"x = "grpc"x ]; then
    if [ ! -d "../api/types" ]; then
        sponge patch gen-types-pb --out=.
        checkResult $?
        mv -f api/types ../api
        rmdir api
    fi
fi
