## ggorm

`ggorm` library wrapped in [gorm](gorm.io/gorm), with added features such as tracer, paging queries, etc.

Support `mysql`, `postgresql`, `sqlite`.

<br>

## Examples of use

### mysql

#### Initializing the connection

```go
    import (
        "github.com/go-dev-frame/sponge/pkg/ggorm"
    )

    var dsn = "root:123456@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

    // (1) connect to the database using the default settings
    db, err := ggorm.InitMysql(dsn)

    // (2) customised settings to connect to the database
    db, err := ggorm.InitMysql(
        dsn,
        ggorm.WithLogging(logger.Get()),  // print log
        ggorm.WithLogRequestIDKey("request_id"),  // print request_id
        ggorm.WithMaxIdleConns(5),
        ggorm.WithMaxOpenConns(50),
        ggorm.WithConnMaxLifetime(time.Minute*3),
        // ggorm.WithSlowThreshold(time.Millisecond*100),  // only print logs that take longer than 100 milliseconds to execute
        // ggorm.WithEnableTrace(),  // enable tracing
        // ggorm.WithRWSeparation(SlavesDsn, MastersDsn...)  // read-write separation
        // ggorm.WithGormPlugin(yourPlugin)  // custom gorm plugin
    )
```

<br>

#### Model

```go
package model

import (
	"github.com/go-dev-frame/sponge/pkg/ggorm"
)

// UserExample object fields mapping table
type UserExample struct {
	ggorm.Model `gorm:"embedded"`

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

// TableName get table name
func (table *UserExample) TableName() string {
	return ggorm.GetTableName(table)
}
```

<br>

#### Transaction

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

### Postgresql

```go
import (
   "github.com/go-dev-frame/sponge/pkg/ggorm"
   "github.com/go-dev-frame/sponge/pkg/utils"
)

func InitSqlite() {
	opts := []ggorm.Option{
		ggorm.WithMaxIdleConns(10),
		ggorm.WithMaxOpenConns(100),
		ggorm.WithConnMaxLifetime(time.Duration(10) * time.Minute),
		ggorm.WithLogging(logger.Get()),
		ggorm.WithLogRequestIDKey("request_id"),
	}

	dsn := "root:123456@127.0.0.1:5432/test"
	dsn = utils.AdaptivePostgresqlDsn(dsn)
	db, err := ggorm.InitPostgresql(dsn, opts...)
	if err != nil {
		panic("ggorm.InitPostgresql error: " + err.Error())
	}
}
```

<br>

### Tidb

Tidb is mysql compatible, just use **InitMysql**.

<br>

### Sqlite

```go
import (
   "github.com/go-dev-frame/sponge/pkg/ggorm"
)

func InitSqlite() {
	opts := []ggorm.Option{
		ggorm.WithMaxIdleConns(10),
		ggorm.WithMaxOpenConns(100),
		ggorm.WithConnMaxLifetime(time.Duration(10) * time.Minute),
		ggorm.WithLogging(logger.Get()),
		ggorm.WithLogRequestIDKey("request_id"),
	}

	dbFile: = "test.db"
	db, err := ggorm.InitSqlite(dbFile, opts...)
	if err != nil {
		panic("ggorm.InitSqlite error: " + err.Error())
	}
}
```

<br>

### gorm User Guide

- https://gorm.io/zh_CN/docs/index.html
