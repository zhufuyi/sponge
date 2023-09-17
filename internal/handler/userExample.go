package handler

import (
	"context"
	"errors"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/internal/types"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

var _ UserExampleHandler = (*userExampleHandler)(nil)

// UserExampleHandler defining the handler interface
type UserExampleHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	List(c *gin.Context)
}

type userExampleHandler struct {
	iDao dao.UserExampleDao
}

// NewUserExampleHandler creating the handler interface
func NewUserExampleHandler() UserExampleHandler {
	return &userExampleHandler{
		iDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create userExample
// @Description submit information to create userExample
// @Tags userExample
// @accept json
// @Produce json
// @Param data body types.CreateUserExampleRequest true "userExample information"
// @Success 200 {object} types.CreateUserExampleRespond{}
// @Router /api/v1/userExample [post]
func (h *userExampleHandler) Create(c *gin.Context) {
	form := &types.CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUserExample)
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	err = h.iDao.Create(ctx, userExample)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": userExample.ID})
}

// DeleteByID delete a record by id
// @Summary delete userExample
// @Description delete userExample by id
// @Tags userExample
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUserExampleByIDRespond{}
// @Router /api/v1/userExample/{id} [delete]
func (h *userExampleHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserExampleIDFromPath(c)
	if isAbort {
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// DeleteByIDs delete records by batch id
// @Summary delete userExamples
// @Description delete userExamples by batch id
// @Tags userExample
// @Param data body types.DeleteUserExamplesByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteUserExamplesByIDsRespond{}
// @Router /api/v1/userExample/delete/ids [post]
func (h *userExampleHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUserExamplesByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	err = h.iDao.DeleteByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update userExample
// @Description update userExample information by id
// @Tags userExample
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUserExampleByIDRequest true "userExample information"
// @Success 200 {object} types.UpdateUserExampleByIDRespond{}
// @Router /api/v1/userExample/{id} [put]
func (h *userExampleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUserExampleIDFromPath(c)
	if isAbort {
		return
	}

	form := &types.UpdateUserExampleByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUserExample)
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	err = h.iDao.UpdateByID(ctx, userExample)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get userExample detail
// @Description get userExample detail by id
// @Tags userExample
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserExampleByIDRespond{}
// @Router /api/v1/userExample/{id} [get]
func (h *userExampleHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getUserExampleIDFromPath(c)
	if isAbort {
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	userExample, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UserExampleObjDetail{}
	err = copier.Copy(data, userExample)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUserExample)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"userExample": data})
}

// GetByCondition get a record by condition
// @Summary get userExample by condition
// @Description get userExample by condition
// @Tags userExample
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserExampleByConditionRespond{}
// @Router /api/v1/userExample/condition [post]
func (h *userExampleHandler) GetByCondition(c *gin.Context) {
	form := &types.GetUserExampleByConditionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = form.Conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	userExample, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UserExampleObjDetail{}
	err = copier.Copy(data, userExample)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUserExample)
		return
	}
	data.ID = utils.Uint64ToStr(userExample.ID)

	response.Success(c, gin.H{"userExample": data})
}

// ListByIDs list of records by batch id
// @Summary list of userExamples by batch id
// @Description list of userExamples by batch id
// @Tags userExample
// @Param data body types.ListUserExamplesByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListUserExamplesByIDsRespond{}
// @Router /api/v1/userExample/list/ids [post]
func (h *userExampleHandler) ListByIDs(c *gin.Context) {
	form := &types.ListUserExamplesByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	userExampleMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	userExamples := []*types.UserExampleObjDetail{}
	for _, id := range form.IDs {
		if v, ok := userExampleMap[id]; ok {
			record, err := convertUserExample(v)
			if err != nil {
				response.Error(c, ecode.ErrListUserExample)
				return
			}
			userExamples = append(userExamples, record)
		}
	}

	response.Success(c, gin.H{
		"userExamples": userExamples,
	})
}

// List of records by query parameters
// @Summary list of userExamples by query parameters
// @Description list of userExamples by paging and conditions
// @Tags userExample
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserExamplesRespond{}
// @Router /api/v1/userExample/list [post]
func (h *userExampleHandler) List(c *gin.Context) {
	form := &types.ListUserExamplesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := context.WithValue(c.Request.Context(), middleware.ContextRequestIDKey, middleware.GCtxRequestID(c)) //nolint
	userExamples, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		response.Error(c, ecode.ErrListUserExample)
		return
	}

	response.Success(c, gin.H{
		"userExamples": data,
		"total":        total,
	})
}

func getUserExampleIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return "", 0, true
	}

	return idStr, id, false
}

func convertUserExample(userExample *model.UserExample) (*types.UserExampleObjDetail, error) {
	data := &types.UserExampleObjDetail{}
	err := copier.Copy(data, userExample)
	if err != nil {
		return nil, err
	}
	data.ID = utils.Uint64ToStr(userExample.ID)
	return data, nil
}

func convertUserExamples(fromValues []*model.UserExample) ([]*types.UserExampleObjDetail, error) {
	toValues := []*types.UserExampleObjDetail{}
	for _, v := range fromValues {
		data, err := convertUserExample(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
