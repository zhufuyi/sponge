FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# 设置时区为上海
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# todo 根据服务类型(http或grpc)选择健康检查工具
# 添加curl，用在http服务的检查
RUN apk add curl

# 添加grpc_health_probe，用在grpc服务的健康检查
#COPY grpc_health_probe /bin/grpc_health_probe
#RUN chmod +x /bin/grpc_health_probe

COPY configs/ /app/configs/
COPY serverNameExample /app/serverNameExample
RUN chmod +x /app/serverNameExample

WORKDIR /app

# http端口，grpc服务可以忽略
EXPOSE 8080

# todo 选择是否使用配置中心配置
CMD ["./serverNameExample", "-c", "configs/serverNameExample.yml"]
# 使用配置中心，serverNameExample.yml是配置中心的配置。
#CMD ["./serverNameExample", "-c", "configs/serverNameExample.yml", "-enable-cc"]
