package service

import (
	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/model"
	"testing"
	"time"

	pb "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/api/types"
	"github.com/zhufuyi/sponge/pkg/gotest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func newUserExampleService() *gotest.Service {
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

	// 初始化mock service
	s := gotest.NewService(d, testData)
	pb.RegisterUserExampleServiceServer(s.Server, &userExampleService{
		UnimplementedUserExampleServiceServer: pb.UnimplementedUserExampleServiceServer{},
		iDao:                                  d.IDao.(dao.UserExampleDao),
	})

	s.GoGrpcServer()

	time.Sleep(time.Millisecond * 100)

	s.IServiceClient = pb.NewUserExampleServiceClient(s.GetClientConn())

	return s
}

func Test_userExampleService_Create(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := &pb.CreateUserExampleRequest{}
	_ = copier.Copy(testData, s.TestData.(*model.UserExample))

	s.MockDao.SqlMock.ExpectBegin()
	args := s.MockDao.GetAnyArgs(s.TestData)
	s.MockDao.SqlMock.ExpectExec("INSERT INTO .*").
		WithArgs(args[:len(args)-1]...). // 根据实际参数数量修改
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDao.SqlMock.ExpectCommit()

	reply, err := s.IServiceClient.(pb.UserExampleServiceClient).Create(s.Ctx, testData)
	//assert.NoError(t, err)

	t.Log(err, reply.String())
}

func Test_userExampleService_DeleteByID(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := &pb.DeleteUserExampleByIDRequest{
		Id: s.TestData.(*model.UserExample).ID,
	}

	reply, err := s.IServiceClient.(pb.UserExampleServiceClient).DeleteByID(s.Ctx, testData)
	assert.NoError(t, err)

	t.Log(reply.String())
}

func Test_userExampleService_UpdateByID(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &pb.UpdateUserExampleByIDRequest{}
	_ = copier.Copy(testData, s.TestData.(*model.UserExample))
	testData.Id = data.ID

	s.MockDao.SqlMock.ExpectBegin()
	s.MockDao.SqlMock.ExpectExec("UPDATE .*").
		WithArgs(s.MockDao.AnyTime, testData.Id). // 根据测试数据数量调整
		WillReturnResult(sqlmock.NewResult(int64(testData.Id), 1))
	s.MockDao.SqlMock.ExpectCommit()

	reply, err := s.IServiceClient.(pb.UserExampleServiceClient).UpdateByID(s.Ctx, testData)
	assert.NoError(t, err)

	t.Log(reply.String())
}

func Test_userExampleService_GetByID(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &pb.GetUserExampleByIDRequest{
		Id: data.ID,
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(data.ID, data.CreatedAt, data.UpdatedAt)

	s.MockDao.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(testData.Id).
		WillReturnRows(rows)

	reply, err := s.IServiceClient.(pb.UserExampleServiceClient).GetByID(s.Ctx, testData)
	assert.NoError(t, err)

	t.Log(reply.String())
}

func Test_userExampleService_ListByIDs(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	data := s.TestData.(*model.UserExample)
	testData := &pb.ListUserExampleByIDsRequest{
		Ids: []uint64{data.ID},
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(data.ID, data.CreatedAt, data.UpdatedAt)

	s.MockDao.SqlMock.ExpectQuery("SELECT .*").
		WithArgs(data.ID).
		WillReturnRows(rows)

	reply, err := s.IServiceClient.(pb.UserExampleServiceClient).ListByIDs(s.Ctx, testData)
	assert.NoError(t, err)

	t.Log(reply.String())
}

func Test_userExampleService_List(t *testing.T) {
	s := newUserExampleService()
	defer s.Close()
	testData := s.TestData.(*model.UserExample)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	s.MockDao.SqlMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	reply, err := s.IServiceClient.(pb.UserExampleServiceClient).List(s.Ctx, &pb.ListUserExampleRequest{
		Params: &types.Params{
			Page:  0,
			Limit: 10,
			Sort:  "ignore count", // 忽略测试 select count(*)
		},
	})
	assert.NoError(t, err)

	t.Log(reply.String())
}

func Test_covertUserExample(t *testing.T) {
	testData := &model.UserExample{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	data, err := covertUserExample(testData)
	assert.NoError(t, err)

	t.Logf("%+v", data)
}
