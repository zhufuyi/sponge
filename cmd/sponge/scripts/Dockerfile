FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

USER root

# set the time zone to Shanghai
#RUN apk add tzdata  \
#    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
#    && echo "Asia/Shanghai" > /etc/timezone \
#    && apk del tzdata

COPY .sponge /root/.sponge
COPY sponge /usr/local/bin/sponge
RUN chmod +x /usr/local/bin/sponge

EXPOSE 24631

WORKDIR /app

ENTRYPOINT ["sponge", "run"]
