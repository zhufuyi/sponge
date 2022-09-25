package service

import (
	"context"

	pb "github.com/zhufuyi/sponge/api/userExample/v1"
	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

// nolint
func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
		pb.RegisterUserExampleServiceServer(server, NewUserExampleServiceServer()) // 把service注册到rpc服务中
	})
}

var _ pb.UserExampleServiceServer = (*userExampleService)(nil)

type userExampleService struct {
	pb.UnimplementedUserExampleServiceServer

	iDao dao.UserExampleDao
}

// NewUserExampleServiceServer 创建一个实例
func NewUserExampleServiceServer() pb.UserExampleServiceServer {
	return &userExampleService{
		iDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetRedisCli()),
		),
	}
}

// Create 创建一条记录
func (s *userExampleService) Create(ctx context.Context, req *pb.CreateUserExampleRequest) (*pb.CreateUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInvalidParams.Err()
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, req)
	if err != nil {
		logger.Warn("copier.Copy error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInternalServerError.Err()
	}

	err = s.iDao.Create(ctx, userExample)
	if err != nil {
		logger.Error("s.iDao.Create error", logger.Err(err), logger.Any("userExample", userExample))
		return nil, ecode.StatusCreateUserExample.Err()
	}

	return &pb.CreateUserExampleReply{Id: userExample.ID}, nil
}

// DeleteByID 根据id删除一条记录
func (s *userExampleService) DeleteByID(ctx context.Context, req *pb.DeleteUserExampleByIDRequest) (*pb.DeleteUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInvalidParams.Err()
	}

	err = s.iDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Error("s.iDao.DeleteByID error", logger.Err(err), logger.Any("id", req.Id))
		return nil, ecode.StatusDeleteUserExample.Err()
	}

	return &pb.DeleteUserExampleByIDReply{}, nil
}

// UpdateByID 根据id更新一条记录
func (s *userExampleService) UpdateByID(ctx context.Context, req *pb.UpdateUserExampleByIDRequest) (*pb.UpdateUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInvalidParams.Err()
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, req)
	if err != nil {
		logger.Warn("copier.Copy error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInternalServerError.Err()
	}
	userExample.ID = req.Id

	err = s.iDao.UpdateByID(ctx, userExample)
	if err != nil {
		logger.Error("s.iDao.UpdateByID error", logger.Err(err), logger.Any("userExample", userExample))
		return nil, ecode.StatusUpdateUserExample.Err()
	}

	return &pb.UpdateUserExampleByIDReply{}, nil
}

// GetByID 根据id查询一条记录
func (s *userExampleService) GetByID(ctx context.Context, req *pb.GetUserExampleByIDRequest) (*pb.GetUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record, err := s.iDao.GetByID(ctx, req.Id)
	if err != nil {
		if err.Error() == query.ErrNotFound.Error() {
			logger.Warn("s.iDao.GetByID error", logger.Err(err), logger.Any("id", req.Id))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("s.iDao.GetByID error", logger.Err(err), logger.Any("id", req.Id))
		return nil, ecode.StatusGetUserExample.Err()
	}

	data, err := covertUserExample(record)
	if err != nil {
		logger.Warn("covertUserExample error", logger.Err(err), logger.Any("record", record))
		return nil, ecode.StatusInternalServerError.Err()
	}

	return &pb.GetUserExampleByIDReply{UserExample: data}, nil
}

// ListByIDs 根据id数组获取多条记录
func (s *userExampleService) ListByIDs(ctx context.Context, req *pb.ListUserExampleByIDsRequest) (*pb.ListUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInvalidParams.Err()
	}

	records, err := s.iDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("s.iDao.GetByID error", logger.Err(err), logger.Any("ids", req.Ids))
		return nil, ecode.StatusGetUserExample.Err()
	}

	datas := []*pb.UserExample{}
	for _, record := range records {
		data, err := covertUserExample(record)
		if err != nil {
			logger.Warn("covertUserExample error", logger.Err(err), logger.Any("id", record.ID))
			continue
		}
		datas = append(datas, data)
	}

	return &pb.ListUserExampleByIDsReply{UserExamples: datas}, nil
}

// List 获取多条记录
func (s *userExampleService) List(ctx context.Context, req *pb.ListUserExampleRequest) (*pb.ListUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req))
		return nil, ecode.StatusInvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		logger.Warn("copier.Copy error", logger.Err(err), logger.Any("params", req.Params))
		return nil, ecode.StatusInternalServerError.Err()
	}
	params.Size = int(req.Params.Limit)

	records, total, err := s.iDao.GetByColumns(ctx, params)
	if err != nil {
		logger.Error("s.iDao.GetByColumns error", logger.Err(err), logger.Any("params", params))
		return nil, ecode.StatusListUserExample.Err()
	}

	datas := []*pb.UserExample{}
	for _, record := range records {
		data, err := covertUserExample(record)
		if err != nil {
			logger.Warn("covertUserExample error", logger.Err(err), logger.Any("id", record.ID))
			continue
		}
		datas = append(datas, data)
	}

	return &pb.ListUserExampleReply{
		Total:        total,
		UserExamples: datas,
	}, nil
}

func covertUserExample(record *model.UserExample) (*pb.UserExample, error) {
	value := &pb.UserExample{}
	err := copier.Copy(value, record)
	if err != nil {
		return nil, err
	}
	value.Id = record.ID
	value.CreatedAt = record.CreatedAt.UnixNano()
	value.UpdatedAt = record.UpdatedAt.UnixNano()
	return value, nil
}
