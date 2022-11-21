## Generate one-way and two-way authentication certificates

openssl.cnf is the openssl file, the configuration has been changed to generate SAN certificates, if you use the go1.15 or higher tls package you must use SAN certificates.

gencert.sh is a certificate generation script that executes the command `bash gencert.sh` to generate both one-way and two-way authentication certificates.
