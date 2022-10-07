FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# 设置时区为上海
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# todo generate dockerfile code for http or grpc here
# delete the templates code start

# 添加curl，用在http服务的检查，如果用部署在k8s，可以不用安装
RUN apk add curl

# 添加grpc_health_probe，用在grpc服务的健康检查
COPY grpc_health_probe /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe

COPY configs/ /app/configs/
COPY serverNameExample /app/serverNameExample
RUN chmod +x /app/serverNameExample

# http端口，grpc服务可以忽略
EXPOSE 8080

# delete the templates code end

WORKDIR /app

CMD ["./serverNameExample", "-c", "configs/serverNameExample.yml"]
# 如果使用配置中心，serverNameExample.yml改为配置中心的配置。
#CMD ["./serverNameExample", "-c", "configs/serverNameExample.yml", "-enable-cc"]
