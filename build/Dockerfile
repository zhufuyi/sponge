# 需要先打包代码`tar zcf sponge.tar.gz $(ls)`，并移动到和Dokerfile在同一个目录下
# 编译go代码
FROM golang:1.18-alpine as build
COPY . /go/src/sponge
WORKDIR /go/src/sponge
RUN tar zxf sponge.tar.gz
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /sponge cmd/sponge/main.go

# 构建grpc_health_probe
RUN go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.12
RUN cd $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-health-probe@v0.4.12 \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /grpc_health_probe

# 使用二进制文件构建镜像
FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# 设置时区为上海
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 添加curl，用在http服务的检查，如果用部署在k8s，可以不用安装
RUN apk add curl

# 添加grpc_health_probe，用在grpc服务的健康检查
COPY --from=build /grpc_health_probe /bin/grpc_health_probe

COPY --from=build /sponge /app/sponge
COPY --from=build /go/src/sponge/config/conf.yml /app/config/conf.yml

WORKDIR /app
CMD ["./sponge", "-c", "config/conf.yml"]

