#!/bin/bash

# As SAN certificates are required for GO version 1.15 and above, steps to generate a SAN certificate
# copy /etc/pki/tls/openssl.cnf to the current directory and modify openssl.cnf
# (1) uncomment copy_extensions = copy
# (2) uncomment req_extensions = v3_req
# (3) add v3_req
  # [ v3_req ]
  # subjectAltName = @alt_names
# (4) add alt_names
  # [alt_names]
  # DNS.1 = localhost

# openssl.cnf file
opensslCnfFile=./openssl.cnf


# generating one-way authentication certificates
mkdir one-way && cd one-way

openssl genrsa -out server.key 2048
openssl req -new -x509  -days 3650 -sha256 -key server.key  -out server.crt -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ${opensslCnfFile} -extensions v3_req

echo "A one-way certificate has been generated and is stored in the directory 'one-way'"
cd .. 



# generate two-way authentication certificates
mkdir two-way && cd two-way

# generate a ca certificate
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.pem -subj "/C=cn/OU=custer/O=custer/CN=localhost"

# server
openssl genpkey -algorithm RSA -out server.key
openssl req -new -nodes -key server.key -out server.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ${opensslCnfFile} -extensions v3_req
openssl x509 -req -days 3650 -in server.csr -out server.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile  ${opensslCnfFile} -extensions v3_req

# client
openssl genpkey -algorithm RSA -out client.key
openssl req -new -nodes -key client.key -out client.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ${opensslCnfFile} -extensions v3_req
openssl x509 -req -days 3650 -in client.csr -out client.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ${opensslCnfFile} -extensions v3_req

echo "A two-way certificate has been generated and is stored in the directory 'two-way'"
cd ..
