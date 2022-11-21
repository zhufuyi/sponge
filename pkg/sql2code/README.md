## sql2code

Generate code for different purposes according to sql, support generating json, gorm model, update parameter, request parameter code, sql can be obtained from parameter, file, db three ways, priority from high to low.

<br>

### Example of use

Main setting parameters.

```go
type Args struct {
	SQL string // DDL sql

	DDLFile string // DDL file

	DBDsn   string // connecting to mysql's dsn
	DBTable string // table name

	Package        string // specify the package name (only valid for model types)
	GormType       bool   // gorm type
	JSONTag        bool   // does it include a json tag
	JSONNamedType  int    // json naming type, 0: consistent with the column name, other values indicate a hump
	IsEmbed        bool   // is gorm.Model embedded
	CodeType       string // specify the different types of code to be generated, namely model (default), json, dao, handler, proto
}
```

<br>

Generated code example.

```go
    // generate gorm model code
    code, err := sql2code.GenerateOne(&sql2code.Args{
        SQL: sqlData,  // source from sql text
        // DDLFile: "user.sql", // source from sql file
        // DBDsn: "root:123456@(127.0.0.1:3306)/account"
        // DBTable "user"
        GormType: true,
        JSONTag: true,
        IsEmbed: true,
        CodeType: "model"
    })

      // generate json, model, dao, handler code
      codes, err := sql2code.Generate(&sql2code.Args{
          SQL: sqlData,  // source from sql text
          // DDLFile: "user.sql", // source from sql file
          // DBDsn: "root:123456@(127.0.0.1:3306)/account"
          // DBTable "user"
          GormType: true,
          JSONTag: true,
          IsEmbed: true,
          CodeType: "dao"
      })
```
