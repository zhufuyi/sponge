启动nacos服务，在浏览器打开ui界面 `http://<ip>:8848/nacos/index.html`

- 名称空间: 区分不同服务
- 组: 用来区分不同环境
- 配置集：具体服务配置

创建一个开发环境服务配置流程：

点击左边菜单栏【名称空间】--> 点击【新建名称空间】，示例：ID为空，名称空间名:serverName，描述:xxx服务，然后点击左边菜单栏【配置管理】-->【配置列表】-->【serverName】--> 【+】新建配置，输入Data ID: serverName.yml，Group: dev，配置格式: yaml，天界配置内容。
