#### 启动 elasticsearch 服务

切换到elasticsearch目录，启动服务

> docker-compose up -d

<br>

#### 启动jaeger 服务

切换到jaeger目录，打开.env环境变量，填写elasticsearch的url、登录账号、密码

启动jaeger

> docker-compose up -d

查看是否正常

> docker-compose ps

