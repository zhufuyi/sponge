package dao

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newUserExampleDao() *gotest.Dao {
	testData := &model.UserExample{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// 初始化mock cache
	//c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // 为了测试mysql，禁止缓存
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewUserExampleCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

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

	// zero id error
	err = d.IDao.(UserExampleDao).DeleteByID(d.Ctx, 0)
	assert.Error(t, err)
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

	// zero id error
	err = d.IDao.(UserExampleDao).UpdateByID(d.Ctx, &model.UserExample{})
	assert.Error(t, err)
	// delete the templates code start
	// update error
	testData = &model.UserExample{
		Name:     "foo",
		Password: "f447b20a7fcbf53a5d5be013ea0b15af",
		Email:    "foo@bar.com",
		Phone:    "16000000001",
		Avatar:   "http://foo/1.jpg",
		Age:      10,
		Gender:   1,
		Status:   1,
		LoginAt:  time.Now().Unix(),
	}
	testData.ID = 1
	err = d.IDao.(UserExampleDao).UpdateByID(d.Ctx, testData)
	assert.Error(t, err)
	// delete the templates code end
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

	_, err := d.IDao.(UserExampleDao).GetByID(d.Ctx, testData.ID) // notfound
	if err != nil {
		t.Fatal(err)
	}

	err = d.SqlMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// notfound error
	d.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(2).
		WillReturnRows(rows)
	_, err = d.IDao.(UserExampleDao).GetByID(d.Ctx, 2)
	assert.Error(t, err)

	d.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(3, 4).
		WillReturnRows(rows)
	_, err = d.IDao.(UserExampleDao).GetByID(d.Ctx, 4)
	assert.Error(t, err)
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

	_, err = d.IDao.(UserExampleDao).GetByIDs(d.Ctx, []uint64{111})
	assert.Error(t, err)

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

	// err test
	_, _, err = d.IDao.(UserExampleDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Columns: []query.Column{
			{
				Name:  "id",
				Exp:   "<",
				Value: 0,
			},
		},
	})
	assert.Error(t, err)

	// error test
	dao := &userExampleDao{}
	_, _, err = dao.GetByColumns(context.Background(), &query.Params{Columns: []query.Column{{}}})
	t.Log(err)
}
