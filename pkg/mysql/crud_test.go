package mysql

import (
	"context"
	"fmt"
	"testing"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"gorm.io/gorm"
)

var table = &userExample{}
var db = initDB()
var ctx = context.Background()

type userExample struct {
	Model `gorm:"embedded"`

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

func initDB() *gorm.DB {
	db, err := Init(dsn, WithLog())
	if err != nil {
		panic(err)
	}

	return db
}

func TestTableName(t *testing.T) {
	t.Logf("table name = %s", TableName(table))
}

func TestCreate(t *testing.T) {
	user := &userExample{Name: "姜维", Age: 20, Gender: "男"}
	err := Create(ctx, db, user)
	if err != nil {
		t.Error(err)
	}

	if user.ID == 0 {
		t.Error("insert failed")
		return
	}

	t.Logf("id =%d", user.ID)
}

func TestDelete(t *testing.T) {
	err := Delete(ctx, db, table, "name = ?", "姜维")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteByID(t *testing.T) {
	err := Delete(ctx, db, table, "id = ?", 25)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	err := Update(ctx, db, table, "age", gorm.Expr("age  + ?", 1), "name = ?", "姜维")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdates(t *testing.T) {
	update := KV{"age": gorm.Expr("age  + ?", 1)}
	err := Updates(ctx, db, table, update, "gender = ?", "女")
	if err != nil {
		t.Error(err)
	}
}

func TestGetByID(t *testing.T) {
	table := &userExample{}
	err := GetByID(ctx, db, table, 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", table)
}

func TestGet(t *testing.T) {
	table := &userExample{}
	err := Get(ctx, db, table, "name = ?", "刘备")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", table)
}

func TestList(t *testing.T) {
	page := query.NewPage(0, 10, "-name")
	tables := []userExample{}
	err := List(ctx, db, &tables, page, "")
	if err != nil {
		t.Error(err)
		return
	}
	for _, user := range tables {
		t.Logf("%+v", user)
	}
}

func TestCount(t *testing.T) {
	count, err := Count(ctx, db, table, "id > ?", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("count=%d", count)
}

// 事务
func TestTx(t *testing.T) {
	err := createUser()
	if err != nil {
		t.Fatal(err)
	}
}

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

	if err = tx.WithContext(ctx).Where("id = ?", 1).First(table).Error; err != nil {
		tx.Rollback()
		return err
	}

	panic("发生了异常")

	if err = tx.WithContext(ctx).Create(&userExample{Name: "lisi", Age: table.Age + 2, Gender: "男"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
