## mysql客户端

在[gorm](gorm.io/gorm)基础上封装的库，添加了链路跟踪，分页查询等功能。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/mysql

<br>

## 使用示例

### 初始化连接示例

```go
    var dsn = "root:123456@(192.168.1.6:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

    // (1) 使用默认设置连接数据库
    db, err := mysql.Init(dsn)

    // (2) 自定义设置连接数据库
	db, err := Init(
		dsn,
		//WithLog(), // 打印所有日志
		WithSlowThreshold(time.Millisecond*100), // 只打印执行时间超过100毫秒的日志
		WithEnableTrace(),                       // 开启链路跟踪
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
	)
```

<br>

### model示例

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

### 事务示例

```go
func createUser() error {
	// 注意，当你在一个事务中应使用 tx 作为数据库句柄
	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil { // 在事务执行过程发生panic后回滚
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

	panic("发生了异常")

	if err = tx.Create(&userExample{Name: "lisi", Age: table.Age + 2, Gender: "男"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
```
<br>

更多使用查看gorm的使用指南

- https://gorm.io/zh_CN/docs/index.html
- https://learnku.com/docs/gorm/v2
