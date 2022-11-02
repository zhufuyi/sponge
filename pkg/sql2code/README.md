## sql2code

根据sql生成不同用途代码，支持生成json、gorm model、update参数、request 参数代码，sql可以从参数、文件、db三种方式获取，优先从高到低。

<br>

### 使用示例

主要设置参数

```go
type Args struct {
	SQL string // DDL sql

	DDLFile string // 读取文件的DDL sql

	DBDsn   string // 从db获取表的DDL sql
	DBTable string

	Package        string // 生成字段的包名(只有model类型有效)
	GormType       bool   // gorm type
	JSONTag        bool   // 是否包括json tag
	JSONNamedType  int    // json命名类型，0:和列名一致，其他值表示驼峰
	IsEmbed        bool   // 是否嵌入gorm.Model
	CodeType       string // 指定生成代码用途，支持4中类型，分别是 model(默认), json, dao, handler
}
```

<br>

生成代码示例：

```go
    // 生成gorm model 代码
    code, err := sql2code.GenerateOne(&sql2code.Args{
        SQL: sqlData,  // 来源于sql语句
        // DDLFile: "user.sql", // 来源于sql文件
        // DBDsn: "root:123456@(127.0.0.1:3306)/account"
        // DBTable "user"
        GormType: true,
        JSONTag: true,
        IsEmbed: true,
        CodeType: "model"
    })

      // 生成json、model、dao、handler代码
      codes, err := sql2code.Generate(&sql2code.Args{
          SQL: sqlData,  // 来源于sql语句
          // DDLFile: "user.sql", // 来源于sql文件
          // DBDsn: "root:123456@(127.0.0.1:3306)/account"
          // DBTable "user"
          GormType: true,
          JSONTag: true,
          IsEmbed: true,
          CodeType: "model"
      })
```
