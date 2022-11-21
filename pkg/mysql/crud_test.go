package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var table = &userExample{}

type userExample struct {
	Model `gorm:"embedded"`

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

func newUserExampleDao() *gotest.Dao {
	testData := &userExample{Name: "ZhangSan", Age: 20, Gender: "male"}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// init mock dao
	d := gotest.NewDao(nil, testData)

	return d
}

func TestTableName(t *testing.T) {
	t.Logf("table name = %s", TableName(&userExample{}))
}

func TestCreate(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := Create(d.Ctx, d.DB, testData)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.Name).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := Delete(d.Ctx, d.DB, table, "name = ?", testData.Name)
	assert.NoError(t, err)
}

func TestDeleteByID(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := Delete(d.Ctx, d.DB, table, "id = ?", testData.ID)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(sqlmock.AnyArg(), d.AnyTime, testData.Name).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := Update(d.Ctx, d.DB, table, "age", gorm.Expr("age  + ?", 1), "name = ?", testData.Name)
	assert.NoError(t, err)
}

func TestUpdates(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(sqlmock.AnyArg(), d.AnyTime, testData.Gender).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	update := KV{"age": gorm.Expr("age  + ?", 1)}
	err := Updates(d.Ctx, d.DB, table, update, "gender = ?", testData.Gender)
	assert.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "age", "gender"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt, testData.Name, testData.Age, testData.Gender)

	d.SQLMock.ExpectQuery("SELECT .*").WithArgs(testData.ID).WillReturnRows(rows)

	err := GetByID(d.Ctx, d.DB, table, testData.ID)
	assert.NoError(t, err)

	t.Logf("%+v", table)
}

func TestGet(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "age", "gender"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt, testData.Name, testData.Age, testData.Gender)

	d.SQLMock.ExpectQuery("SELECT .*").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows) // adjusted for number of fields

	err := Get(d.Ctx, d.DB, table, "name = ?", testData.Name)
	assert.NoError(t, err)

	t.Logf("%+v", table)
}

func TestList(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "age", "gender"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt, testData.Name, testData.Age, testData.Gender)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	page := query.NewPage(0, 10, "")
	tables := []userExample{}
	err := List(d.Ctx, d.DB, &tables, page, "")
	assert.NoError(t, err)

	for _, user := range tables {
		t.Logf("%+v", user)
	}
}

func TestCount(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "age", "gender"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt, testData.Name, testData.Age, testData.Gender)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	count, err := Count(d.Ctx, d.DB, table, "id > ?", 0)
	assert.NotNil(t, err)

	t.Logf("count=%d", count)
}

func TestTx(t *testing.T) {
	err := createUser()
	if err != nil {
		t.Fatal(err)
	}
}

func createUser() error {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*userExample)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "age", "gender"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt, testData.Name, testData.Age, testData.Gender)
	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectQuery("SELECT .*").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows) // adjusted for number of fields
	d.SQLMock.ExpectCommit()

	// note that you should use tx as the database handle when you are in a transaction
	tx := d.DB.Begin()
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

	if err = tx.WithContext(d.Ctx).Where("id = ?", testData.ID).First(table).Error; err != nil {
		tx.Rollback()
		return err
	}

	panic("mock panic")

	if err = tx.WithContext(d.Ctx).Create(&userExample{Name: "lisi", Age: table.Age + 2, Gender: "male"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
