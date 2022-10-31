## jy2struct

一个生成go struct代码库，支持json和yaml。

<br>

### 使用示例

主要设置参数：

```go
type Args struct {
	Format    string // 文档格式，json或yaml
	Data      string // json或yaml内容
	InputFile string // 文件
	Name      string // 结构体名称
	SubStruct bool   // 子结构体是否分开
	Tags      string // 添加额外tag，多个tag用逗号分隔
}
```

<br>

转换示例：

```go
    // json转struct
    code, err := jy2struct.Covert(&jy2struct.Args{
        Format: "json",
        // InputFile: "user.json", // 来源于json文件
        SubStruct: true,
    })

    // json转struct
    code, err := jy2struct.Covert(&jy2struct.Args{
        Format: "yaml",
        // InputFile: "user.yaml", // 来源于yaml文件
        SubStruct: true,
    })
```
