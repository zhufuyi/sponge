package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/copier"
)

func newUserExampleHandler() *gotest.Handler {
	// todo 补充测试字段信息
	testData := &model.UserExample{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// 初始化mock cache
	c := gotest.NewCache(map[string]interface{}{"no cache": testData})
	c.ICache = cache.NewUserExampleCache(c.RedisClient)

	// 初始化mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = dao.NewUserExampleDao(d.DB, c.ICache.(cache.UserExampleCache))

	// 初始化mock handler
	h := gotest.NewHandler(d, testData)
	h.IHandler = &userExampleHandler{iDao: d.IDao.(dao.UserExampleDao)}

	testFns := []gotest.RouterInfo{
		{
			FuncName:    "Create",
			Method:      http.MethodPost,
			Path:        "/userExample",
			HandlerFunc: h.IHandler.(UserExampleHandler).Create,
		},
		{
			FuncName:    "DeleteByID",
			Method:      http.MethodDelete,
			Path:        "/userExample/:id",
			HandlerFunc: h.IHandler.(UserExampleHandler).DeleteByID,
		},
		{
			FuncName:    "UpdateByID",
			Method:      http.MethodPut,
			Path:        "/userExample/:id",
			HandlerFunc: h.IHandler.(UserExampleHandler).UpdateByID,
		},
		{
			FuncName:    "GetByID",
			Method:      http.MethodGet,
			Path:        "/userExample/:id",
			HandlerFunc: h.IHandler.(UserExampleHandler).GetByID,
		},
		{
			FuncName:    "ListByIDs",
			Method:      http.MethodPost,
			Path:        "/userExamples/ids",
			HandlerFunc: h.IHandler.(UserExampleHandler).ListByIDs,
		},
		{
			FuncName:    "List",
			Method:      http.MethodPost,
			Path:        "/userExamples",
			HandlerFunc: h.IHandler.(UserExampleHandler).List,
		},
	}

	h.GoRunHttpServer(testFns)

	return h
}

func Test_userExampleHandler_Create(t *testing.T) {
	h := newUserExampleHandler()
	defer h.Close()
	testData := &CreateUserExampleRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.UserExample))

	h.MockDao.SqlMock.ExpectBegin()
	args := h.MockDao.GetAnyArgs(h.TestData)
	h.MockDao.SqlMock.ExpectExec("INSERT INTO .*").
		WithArgs(args[:len(args)-1]...). // 根据实际参数数量修改
		WillReturnResult(sqlmock.NewResult(1, 1))
	h.MockDao.SqlMock.ExpectCommit()

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("Create"), testData)
	if err != nil {
		t.Fatal(err)
	}

	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_userExampleHandler_DeleteByID(t *testing.T) {
	h := newUserExampleHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	result := &gohttp.StdResult{}
	err := gohttp.Delete(result, h.GetRequestURL("DeleteByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_userExampleHandler_UpdateByID(t *testing.T) {
	h := newUserExampleHandler()
	defer h.Close()
	testData := &UpdateUserExampleByIDRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.UserExample))

	h.MockDao.SqlMock.ExpectBegin()
	h.MockDao.SqlMock.ExpectExec("UPDATE .*").
		WithArgs(h.MockDao.AnyTime, testData.ID). // 根据测试数据数量调整
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	h.MockDao.SqlMock.ExpectCommit()

	result := &gohttp.StdResult{}
	err := gohttp.Put(result, h.GetRequestURL("UpdateByID", testData.ID), testData)
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_userExampleHandler_GetByID(t *testing.T) {
	h := newUserExampleHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Get(result, h.GetRequestURL("GetByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_userExampleHandler_ListByIDs(t *testing.T) {
	h := newUserExampleHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SqlMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("ListByIDs"), &GetUserExamplesByIDsRequest{IDs: []uint64{testData.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}

func Test_userExampleHandler_List(t *testing.T) {
	h := newUserExampleHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	h.MockDao.SqlMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &gohttp.StdResult{}
	err := gohttp.Post(result, h.GetRequestURL("List"), &GetUserExamplesRequest{query.Params{
		Page: 0,
		Size: 10,
		Sort: "ignore count", // 忽略测试 select count(*)
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}
}
