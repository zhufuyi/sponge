## conf

解析yaml、json、toml配置文件到go struct，结合[goctl](https://github.com/zhufuyi/goctl)工具自动生成config.go到指定目录，例如：

> goctl covert yaml --file=test.yaml --tags=json --out=yourServerName/config。

<br>

### 使用示例

```go
    // 方式一：无监听配置文件
	conf := &App{}
	err := Parse("test.yml", conf)

    // 方式二：开启监听配置文件
	conf := &App{}
	fs := []func(){
		func() {
			fmt.Println("监听到配置文件有更新")
		},
	}
	err := Parse("test.yml", conf, fs...)
```
