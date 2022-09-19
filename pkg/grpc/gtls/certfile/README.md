## 生成单向和双向认证证书

openssl.cnf是openssl的文件，配置已经改为生成SAN证书，如果使用go1.15以上的tls包必须使用SAN证书。

gencert.sh是生成证书脚本，执行命令 `bash gencert.sh` 同时生成生成单向和双向认证证书。
