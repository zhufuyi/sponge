## sgorm

`sgorm` is a library encapsulated on [gorm](gorm.io/gorm), and added features such as tracer, paging queries, etc.

Support `mysql`, `postgresql`, `sqlite`.

<br>

## Examples of use

### mysql

#### Initializing the connection

```go
    import "github.com/zhufuyi/sponge/pkg/sgorm/mysql"

    var dsn = "root:123456@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

    // case 1: connect to the database using the default settings
    gdb, err := mysql.Init(dsn)

    // case 2: customised settings to connect to the database
    db, err := mysql.Init(
        dsn,
        mysql.WithLogging(logger.Get()),  // print log
        mysql.WithLogRequestIDKey("request_id"),  // print request_id
        mysql.WithMaxIdleConns(5),
        mysql.WithMaxOpenConns(50),
        mysql.WithConnMaxLifetime(time.Minute*3),
        // mysql.WithSlowThreshold(time.Millisecond*100),  // only print logs that take longer than 100 milliseconds to execute
        // mysql.WithEnableTrace(),  // enable tracing
        // mysql.WithRWSeparation(SlavesDsn, MastersDsn...)  // read-write separation
        // mysql.WithGormPlugin(yourPlugin)  // custom gorm plugin
    )
    
    if err != nil {
        panic("mysql.Init error: " + err.Error())
    }
```

<br>

### Postgresql

```go
    import (
        "github.com/zhufuyi/sponge/pkg/sgorm/postgresql"
        "github.com/zhufuyi/sponge/pkg/utils"
    )

    func InitPostgresql() {
        opts := []postgresql.Option{
            postgresql.WithMaxIdleConns(10),
            postgresql.WithMaxOpenConns(100),
            postgresql.WithConnMaxLifetime(time.Duration(10) * time.Minute),
            postgresql.WithLogging(logger.Get()),
            postgresql.WithLogRequestIDKey("request_id"),
        }

        dsn := "root:123456@127.0.0.1:5432/test"  // or dsn := "host=127.0.0.1 user=root password=123456 dbname=account port=5432 sslmode=disable TimeZone=Asia/Shanghai"
        dsn = utils.AdaptivePostgresqlDsn(dsn)
        db, err := postgresql.Init(dsn, opts...)
        if err != nil {
            panic("postgresql.Init error: " + err.Error())
        }
    }
```

<br>

### Tidb

Tidb is mysql compatible, just use **mysql.Init**.

<br>

### Sqlite

```go
    import "github.com/zhufuyi/sponge/pkg/sgorm/sqlite"

    func InitSqlite() {
        opts := []sgorm.Option{
            sgorm.WithMaxIdleConns(10),
            sgorm.WithMaxOpenConns(100),
            sgorm.WithConnMaxLifetime(time.Duration(10) * time.Minute),
            sgorm.WithLogging(logger.Get()),
            sgorm.WithLogRequestIDKey("request_id"),
        }

        dbFile: = "test.db"
        db, err := sgorm.Init(dbFile, opts...)
        if err != nil {
            panic("sgorm.Init error: " + err.Error())
        }
    }
```

<br>

### Transaction Example

```go
    func createUser() error {
        // note that you should use tx as the database handle when you are in a transaction
        tx := db.Begin()
        defer func() {
            if err := recover(); err != nil { // rollback after a panic during transaction execution
                tx.Rollback()
                fmt.Printf("transaction failed, err = %v\n", err)
            }
        }()

        var err error
        if err = tx.Error; err != nil {
            return err
        }

        if err = tx.Where("id = ?", 1).First(table).Error; err != nil {
            tx.Rollback()
            return err
        }

        panic("mock panic")

        if err = tx.Create(&userExample{Name: "Mr Li", Age: table.Age + 2, Gender: "male"}).Error; err != nil {
            tx.Rollback()
            return err
        }

        return tx.Commit().Error
    }
```

<br>

### Model Embedding Example

```go
package model

import "github.com/zhufuyi/sponge/pkg/sgorm"

// User object fields mapping table
type User struct {
    sgorm.Model `gorm:"embedded"`

    Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
    Age    int    `gorm:"not null" json:"age"`
    Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

// TableName get table name
func (table *User) TableName() string {
    return sgorm.GetTableName(table)
}
```

<br>

### gorm User Guide

- https://gorm.io/zh_CN/docs/index.html
