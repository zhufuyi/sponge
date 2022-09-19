#!/bin/bash

# 由于GO1.15以上版本必须使用SAN证书，生成SAN证书步骤
# 复制/etc/pki/tls/openssl.cnf到当前目录，修改openssl.cnf
# (1) 取消copy_extensions = copy注释
# (2) 取消req_extensions = v3_req注释
# (3) 添加 v3_req
# [ v3_req ]
# subjectAltName = @alt_names
# (4) 添加 alt_names
# [alt_names]
# DNS.1 = localhost

# openssl.cnf文件位置
opensslCnfFile=./openssl.cnf


# 生成单向认证证书
mkdir one-way && cd one-way

openssl genrsa -out server.key 2048
openssl req -new -x509  -days 3650 -sha256 -key server.key  -out server.crt -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ${opensslCnfFile} -extensions v3_req

echo "已生成单向证书，存放在目录one-way"
cd .. 



# 生成双向认证证书
mkdir two-way && cd two-way

# 生成ca证书
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.pem -subj "/C=cn/OU=custer/O=custer/CN=localhost"

# 服务端
openssl genpkey -algorithm RSA -out server.key
openssl req -new -nodes -key server.key -out server.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ${opensslCnfFile} -extensions v3_req
openssl x509 -req -days 3650 -in server.csr -out server.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile  ${opensslCnfFile} -extensions v3_req

# 客户端
openssl genpkey -algorithm RSA -out client.key
openssl req -new -nodes -key client.key -out client.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ${opensslCnfFile} -extensions v3_req
openssl x509 -req -days 3650 -in client.csr -out client.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ${opensslCnfFile} -extensions v3_req

echo "已生成双向证书，存放在目录two-way"
cd ..

