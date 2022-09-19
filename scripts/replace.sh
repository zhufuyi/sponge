#!/bin/bash

basePath=$1
oldName=$2
newName=$3


if [ ! -n "$basePath" ]; then
    basePath="."
    echo "error: missing parameter basePath "
    echo "usage: ./replace.sh <basePath> <oldName> <newName>"
    echo "   eg: ./replace.sh ./ old-name new-name"
    exit 1
fi

if [ ! -n "$oldName" ]; then
    basePath="."
    echo "error: missing parameter oldName"
    echo "usage: ./replace.sh <basePath> <oldName> <newName>"
    echo "   eg: ./replace.sh ./ old-name new-name"
    exit 1
fi

if [ ! -n "$newName" ]; then
    basePath="."
    echo "error: missing parameter newName"
    echo "usage: ./replace.sh <basePath> <oldName> <newName>"
    echo "   eg: ./replace.sh ./ old-name new-name"
    exit 1
fi


function listFiles(){
    cd $1
    items=$(ls)

    for item in $items
    do  
        if [ -d "$item" ]; then
            listFiles $item
        else
            if [ "${item#*.}" = "pb.validate.go" ];then
                # 修改文件内容
                #echo "change file '$item' content: $oldName-->$newName"
                sed -i "s/$oldName/$newName/g" $item
            fi  
        fi  
    done 
    cd ..
}

listFiles $basePath
