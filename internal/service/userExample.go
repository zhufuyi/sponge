package service

import (
	"context"
	"errors"
	"strings"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

// nolint
func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
		serverNameExampleV1.RegisterUserExampleServiceServer(server, NewUserExampleServiceServer()) // register service to the rpc service
	})
}

var _ serverNameExampleV1.UserExampleServiceServer = (*userExampleService)(nil)

type userExampleService struct {
	serverNameExampleV1.UnimplementedUserExampleServiceServer

	iDao dao.UserExampleDao
}

// NewUserExampleServiceServer create a new service
func NewUserExampleServiceServer() serverNameExampleV1.UserExampleServiceServer {
	return &userExampleService{
		iDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetCacheType()),
		),
	}
}

// Create a record
func (s *userExampleService) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, req)
	if err != nil {
		return nil, ecode.StatusCreateUserExample.Err()
	}

	err = s.iDao.Create(ctx, userExample)
	if err != nil {
		logger.Error("s.iDao.Create error", logger.Err(err), logger.Any("userExample", userExample), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.CreateUserExampleReply{Id: userExample.ID}, nil
}

// DeleteByID delete a record by id
func (s *userExampleService) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	err = s.iDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Error("s.iDao.DeleteByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDReply{}, nil
}

// DeleteByIDs delete records by batch id
func (s *userExampleService) DeleteByIDs(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDsRequest) (*serverNameExampleV1.DeleteUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	err = s.iDao.DeleteByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("s.iDao.DeleteByID error", logger.Err(err), logger.Any("ids", req.Ids), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDsReply{}, nil
}

// UpdateByID update a record by id
func (s *userExampleService) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, req)
	if err != nil {
		return nil, ecode.StatusUpdateUserExample.Err()
	}
	userExample.ID = req.Id

	err = s.iDao.UpdateByID(ctx, userExample)
	if err != nil {
		logger.Error("s.iDao.UpdateByID error", logger.Err(err), logger.Any("userExample", userExample), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.UpdateUserExampleByIDReply{}, nil
}

// GetByID get a record by id
func (s *userExampleService) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record, err := s.iDao.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("s.iDao.GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("s.iDao.GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convertUserExample(record)
	if err != nil {
		logger.Warn("convertUserExample error", logger.Err(err), logger.Any("record", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetUserExample.Err()
	}

	return &serverNameExampleV1.GetUserExampleByIDReply{UserExample: data}, nil
}

// ListByIDs list of records by batch id
func (s *userExampleService) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	userExampleMap, err := s.iDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("s.iDao.GetByID error", logger.Err(err), logger.Any("ids", req.Ids), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, id := range req.Ids {
		if v, ok := userExampleMap[id]; ok {
			record, err := convertUserExample(v)
			if err != nil {
				logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", v), interceptor.ServerCtxRequestIDField(ctx))
				return nil, ecode.StatusInternalServerError.ToRPCErr()
			}
			userExamples = append(userExamples, record)
		}
	}

	return &serverNameExampleV1.ListUserExampleByIDsReply{UserExamples: userExamples}, nil
}

// List of records by query parameters
func (s *userExampleService) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.StatusListUserExample.Err()
	}
	params.Size = int(req.Params.Limit)

	records, total, err := s.iDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("s.iDao.GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusInvalidParams.Err()
		}
		logger.Error("s.iDao.GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, record := range records {
		userExample, err := convertUserExample(record)
		if err != nil {
			logger.Warn("convertUserExample error", logger.Err(err), logger.Any("id", record.ID), interceptor.ServerCtxRequestIDField(ctx))
			continue
		}
		userExamples = append(userExamples, userExample)
	}

	return &serverNameExampleV1.ListUserExampleReply{
		Total:        total,
		UserExamples: userExamples,
	}, nil
}

func convertUserExample(record *model.UserExample) (*serverNameExampleV1.UserExample, error) {
	value := &serverNameExampleV1.UserExample{}
	err := copier.Copy(value, record)
	if err != nil {
		return nil, err
	}
	value.Id = record.ID
	value.CreatedAt = record.CreatedAt.Unix()
	value.UpdatedAt = record.UpdatedAt.Unix()
	return value, nil
}
