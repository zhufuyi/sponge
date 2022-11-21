## mysql

A library wrapped on top of [gorm](gorm.io/gorm), with added features such as link tracing, paging queries, etc.

<br>

## Example of use

### Initializing the connection

```go
    var dsn = "root:123456@(192.168.1.6:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

    // (1) connect to the database using the default settings
    db, err := mysql.Init(dsn)

    // (2) customised settings to connect to the database
	db, err := Init(
		dsn,
		//WithLog(), // print all logs
		WithSlowThreshold(time.Millisecond*100), // only print logs that take longer than 100 milliseconds to execute
		WithEnableTrace(),                       // enable tracing
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
	)
```

<br>

### Model

```go
package model

import (
	"github.com/zhufuyi/sponge/pkg/mysql"
)

// UserExample object fields mapping table
type UserExample struct {
	mysql.Model `gorm:"embedded"`

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

// TableName get table name
func (table *UserExample) TableName() string {
	return mysql.GetTableName(table)
}
```

<br>

### Transaction

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

## gorm User Guide

- https://gorm.io/zh_CN/docs/index.html
- https://learnku.com/docs/gorm/v2
