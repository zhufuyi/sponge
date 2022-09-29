```bash
# todo 复制服务配置文件到configs目录下。
cd deployments/docker-compose
mkdir configs
cp ../../configs/serverNameExample.yml configs/

tree
#    ├── configs
#    │         └── serverNameExample.yml
#    └── docker-compose.yml
```

启动服务：

> docker-compose up

