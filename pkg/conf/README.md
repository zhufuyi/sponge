## conf

解析yaml、json、toml配置文件到go struct。

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
