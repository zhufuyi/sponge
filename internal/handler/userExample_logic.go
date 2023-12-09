package handler

import (
	"context"
	"errors"
	"strings"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"

	"github.com/jinzhu/copier"
)

var _ serverNameExampleV1.UserExampleLogicer = (*userExamplePbHandler)(nil)

type userExamplePbHandler struct {
	userExampleDao dao.UserExampleDao
}

// NewUserExamplePbHandler create a handler
func NewUserExamplePbHandler() serverNameExampleV1.UserExampleLogicer {
	return &userExamplePbHandler{
		userExampleDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetCacheType()),
		),
	}
}

// Create a record
func (h *userExamplePbHandler) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, req)
	if err != nil {
		return nil, ecode.ErrCreateUserExample.Err()
	}

	err = h.userExampleDao.Create(ctx, userExample)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("userExample", userExample), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &serverNameExampleV1.CreateUserExampleReply{Id: userExample.ID}, nil
}

// DeleteByID delete a record by id
func (h *userExamplePbHandler) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	err = h.userExampleDao.DeleteByID(ctx, req.Id)
	if err != nil {
		logger.Warn("DeleteByID error", logger.Err(err), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDReply{}, nil
}

// DeleteByIDs delete records by batch id
func (h *userExamplePbHandler) DeleteByIDs(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDsRequest) (*serverNameExampleV1.DeleteUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	err = h.userExampleDao.DeleteByIDs(ctx, req.Ids)
	if err != nil {
		logger.Warn("DeleteByIDs error", logger.Err(err), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &serverNameExampleV1.DeleteUserExampleByIDsReply{}, nil
}

// UpdateByID update a record by id
func (h *userExamplePbHandler) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, req)
	if err != nil {
		return nil, ecode.ErrUpdateByIDUserExample.Err()
	}
	userExample.ID = req.Id

	err = h.userExampleDao.UpdateByID(ctx, userExample)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("userExample", userExample), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	return &serverNameExampleV1.UpdateUserExampleByIDReply{}, nil
}

// GetByID get a record by id
func (h *userExamplePbHandler) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	record, err := h.userExampleDao.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("id", req.Id), middleware.CtxRequestIDField(ctx))
			return nil, ecode.NotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("id", req.Id), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	data, err := convertUserExamplePb(record)
	if err != nil {
		logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", record), middleware.CtxRequestIDField(ctx))
		return nil, ecode.ErrGetByIDUserExample.Err()
	}

	return &serverNameExampleV1.GetUserExampleByIDReply{
		UserExample: data,
	}, nil
}

// GetByCondition get a record by condition
func (h *userExamplePbHandler) GetByCondition(ctx context.Context, req *serverNameExampleV1.GetUserExampleByConditionRequest) (*serverNameExampleV1.GetUserExampleByConditionReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	conditions := &query.Conditions{}
	for _, v := range req.Conditions.GetColumns() {
		column := query.Column{}
		_ = copier.Copy(&column, v)
		conditions.Columns = append(conditions.Columns, column)
	}
	err = conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error", logger.Err(err), logger.Any("conditions", conditions), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	record, err := h.userExampleDao.GetByCondition(ctx, conditions)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByID error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
			return nil, ecode.NotFound.Err()
		}
		logger.Error("GetByID error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	data, err := convertUserExamplePb(record)
	if err != nil {
		logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", record), middleware.CtxRequestIDField(ctx))
		return nil, ecode.ErrGetByIDUserExample.Err()
	}

	return &serverNameExampleV1.GetUserExampleByConditionReply{
		UserExample: data,
	}, nil
}

// ListByIDs list of records by batch id
func (h *userExamplePbHandler) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	userExampleMap, err := h.userExampleDao.GetByIDs(ctx, req.Ids)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("ids", req.Ids), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, id := range req.Ids {
		if v, ok := userExampleMap[id]; ok {
			record, err := convertUserExamplePb(v)
			if err != nil {
				logger.Warn("convertUserExample error", logger.Err(err), logger.Any("userExample", v), middleware.CtxRequestIDField(ctx))
				return nil, ecode.InternalServerError.Err()
			}
			userExamples = append(userExamples, record)
		}
	}

	return &serverNameExampleV1.ListUserExampleByIDsReply{
		UserExamples: userExamples,
	}, nil
}

// ListByLastID get records by last id
func (h *userExamplePbHandler) ListByLastID(ctx context.Context, req *serverNameExampleV1.ListUserExampleByLastIDRequest) (*serverNameExampleV1.ListUserExampleByLastIDReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	records, err := h.userExampleDao.GetByLastID(ctx, req.LastID, int(req.Limit), req.Sort)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, record := range records {
		data, err := convertUserExamplePb(record)
		if err != nil {
			logger.Warn("convertUserExample error", logger.Err(err), logger.Any("id", record.ID), middleware.CtxRequestIDField(ctx))
			continue
		}
		userExamples = append(userExamples, data)
	}

	return &serverNameExampleV1.ListUserExampleByLastIDReply{
		UserExamples: userExamples,
	}, nil
}

// List of records by query parameters
func (h *userExamplePbHandler) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InvalidParams.Err()
	}

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.ErrListUserExample.Err()
	}
	params.Size = int(req.Params.Limit)

	records, total, err := h.userExampleDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), middleware.CtxRequestIDField(ctx))
			return nil, ecode.InvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), middleware.CtxRequestIDField(ctx))
		return nil, ecode.InternalServerError.Err()
	}

	userExamples := []*serverNameExampleV1.UserExample{}
	for _, record := range records {
		data, err := convertUserExamplePb(record)
		if err != nil {
			logger.Warn("convertUserExample error", logger.Err(err), logger.Any("id", record.ID), middleware.CtxRequestIDField(ctx))
			continue
		}
		userExamples = append(userExamples, data)
	}

	return &serverNameExampleV1.ListUserExampleReply{
		Total:        total,
		UserExamples: userExamples,
	}, nil
}

func convertUserExamplePb(record *model.UserExample) (*serverNameExampleV1.UserExample, error) {
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
