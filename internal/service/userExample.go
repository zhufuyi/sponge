package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/logger"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"
)

func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
		serverNameExampleV1.RegisterUserExampleServer(server, NewUserExampleServer()) // register service to the rpc service
	})
}

var _ serverNameExampleV1.UserExampleServer = (*userExample)(nil)
var _ time.Time

type userExample struct {
	serverNameExampleV1.UnimplementedUserExampleServer

	iDao dao.UserExampleDao
}

// NewUserExampleServer create a new service
func NewUserExampleServer() serverNameExampleV1.UserExampleServer {
	return &userExample{
		iDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetCacheType()),
		),
	}
}

// Create a record
func (s *userExample) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	record := &model.UserExample{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusCreateUserExample.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = s.iDao.Create(ctx, record)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.CreateUserExampleReply{Id: record.ID}, nil
}

// DeleteByID delete a record by id
func (s *userExample) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	err = s.iDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDReply{}, nil
}

// UpdateByID update a record by id
func (s *userExample) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	record := &model.UserExample{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusUpdateByIDUserExample.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	record.ID = req.Id

	err = s.iDao.UpdateByID(ctx, record)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.UpdateUserExampleByIDReply{}, nil
}

// GetByID get a record by id
func (s *userExample) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	record, err := s.iDao.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", req.Id), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convertUserExample(record)
	if err != nil {
		logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetByIDUserExample.Err()
	}

	return &serverNameExampleV1.GetUserExampleByIDReply{UserExample: data}, nil
}

// List of records by query parameters
func (s *userExample) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.StatusListUserExample.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	params.Limit = int(req.Params.Limit)

	records, total, err := s.iDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusInvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, record := range records {
		data, err := convertUserExample(record)
		if err != nil {
			logger.Warn("convertUserExample error", logger.Err(err), logger.Any("id", record.ID), interceptor.ServerCtxRequestIDField(ctx))
			continue
		}
		userExamples = append(userExamples, data)
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
	// Note: if copier.Copy cannot assign a value to a field, add it here, e.g. CreatedAt, UpdatedAt
	value.Id = record.ID
	// todo generate the conversion createdAt and updatedAt code here
	// delete the templates code start
	value.CreatedAt = record.CreatedAt.Format(time.RFC3339)
	value.UpdatedAt = record.UpdatedAt.Format(time.RFC3339)
	// delete the templates code end
	return value, nil
}
