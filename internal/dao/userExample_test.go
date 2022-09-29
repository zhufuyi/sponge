package dao

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/internal/serverNameExample/cache"
	"github.com/zhufuyi/sponge/internal/serverNameExample/model"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func newUserExampleDao() *gotest.Dao {
	testData := &model.UserExample{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// 初始化mock cache
	c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // 为了测试mysql，禁止缓存
	c.ICache = cache.NewUserExampleCache(c.RedisClient)

	// 初始化mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = NewUserExampleDao(d.DB, c.ICache.(cache.UserExampleCache))

	return d
}

func Test_userExampleDao_Create(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*model.UserExample)

	d.SqlMock.ExpectBegin()
	d.SqlMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SqlMock.ExpectCommit()

	err := d.IDao.(UserExampleDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleDao_DeleteByID(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*model.UserExample)

	testData.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: false,
	}

	d.SqlMock.ExpectBegin()
	d.SqlMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SqlMock.ExpectCommit()

	err := d.IDao.(UserExampleDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleDao_UpdateByID(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*model.UserExample)

	d.SqlMock.ExpectBegin()
	d.SqlMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SqlMock.ExpectCommit()

	err := d.IDao.(UserExampleDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleDao_GetByID(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(UserExampleDao).GetByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = d.SqlMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleDao_GetByIDs(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(UserExampleDao).GetByIDs(d.Ctx, []uint64{testData.ID})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SqlMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleDao_GetByColumns(t *testing.T) {
	d := newUserExampleDao()
	defer d.Close()
	testData := d.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SqlMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(UserExampleDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Sort: "ignore count", // 忽略测试 select count(*)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SqlMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}
