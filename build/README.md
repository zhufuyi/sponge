
- `Dockerfile`: 直接复制已编译的二进制文件构建出来的镜像。
  - 优点：构建速度快
  - 缺点：镜像体积被两阶段构建大一倍。
- `Dockerfile_build`: 两阶段构建镜像。
  - 优点：镜像体积最小
  - 缺点：构建速度较慢，每次构建都产生比较大的中间镜像，需要定时执行命令`docker rmi $(docker images | grep "<none>" | awk '{print $3}')`清理。
