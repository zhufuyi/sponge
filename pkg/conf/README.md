## conf

解析yaml、json、toml配置文件到go struct，结合[goctl](https://github.com/zhufuyi/goctl)工具自动生成config.go到指定目录，例如：

> goctl covert yaml --file=test.yaml --tags=json --out=/yourProjectName/config。

<br>

### 安装

> go get -u github.com/zhufuyi/pkg/conf

<br>

### 使用示例

具体示例看[config](internal/config)。
