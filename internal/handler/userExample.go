package handler

import (
	"errors"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/internal/types"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

var _ UserExampleHandler = (*userExampleHandler)(nil)

// UserExampleHandler 定义handler接口
type UserExampleHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	ListByIDs(c *gin.Context)
	List(c *gin.Context)
}

type userExampleHandler struct {
	iDao dao.UserExampleDao
}

// NewUserExampleHandler 创建handler接口
func NewUserExampleHandler() UserExampleHandler {
	return &userExampleHandler{
		iDao: dao.NewUserExampleDao(
			model.GetDB(),
			cache.NewUserExampleCache(model.GetRedisCli()),
		),
	}
}

// Create 创建
// @Summary 创建userExample
// @Description 提交信息创建userExample
// @Tags userExample
// @accept json
// @Produce json
// @Param data body types.CreateUserExampleRequest true "userExample信息"
// @Success 200 {object} types.Result{}
// @Router /api/v1/userExample [post]
func (h *userExampleHandler) Create(c *gin.Context) {
	form := &types.CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		logger.Warn("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrCreateUserExample)
		return
	}

	err = h.iDao.Create(c.Request.Context(), userExample)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": userExample.ID})
}

// DeleteByID 根据id删除一条记录
// @Summary 删除userExample
// @Description 根据id删除userExample
// @Tags userExample
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.Result{}
// @Router /api/v1/userExample/{id} [delete]
func (h *userExampleHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserExampleIDFromPath(c)
	if isAbort {
		return
	}

	err := h.iDao.DeleteByID(c.Request.Context(), id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID 根据id更新信息
// @Summary 更新userExample信息
// @Description 根据id更新userExample信息
// @Tags userExample
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUserExampleByIDRequest true "userExample信息"
// @Success 200 {object} types.Result{}
// @Router /api/v1/userExample/{id} [put]
func (h *userExampleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUserExampleIDFromPath(c)
	if isAbort {
		return
	}

	form := &types.UpdateUserExampleByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		logger.Warn("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrUpdateUserExample)
		return
	}

	err = h.iDao.UpdateByID(c.Request.Context(), userExample)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID 根据id获取一条记录
// @Summary 获取userExample详情
// @Description 根据id获取userExample详情
// @Tags userExample
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/userExample/{id} [get]
func (h *userExampleHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getUserExampleIDFromPath(c)
	if isAbort {
		return
	}

	userExample, err := h.iDao.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.GetUserExampleByIDRespond{}
	err = copier.Copy(data, userExample)
	if err != nil {
		logger.Warn("Copy error", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrGetUserExample)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"userExample": data})
}

// ListByIDs 根据多个id获取多条记录
// @Summary 根据多个id获取userExample列表
// @Description 使用post请求，根据多个id获取userExample列表
// @Tags userExample
// @Param data body types.GetUserExamplesByIDsRequest true "id 数组"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/userExamples/ids [post]
func (h *userExampleHandler) ListByIDs(c *gin.Context) {
	form := &types.GetUserExamplesByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userExamples, err := h.iDao.GetByIDs(c.Request.Context(), form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())

		return
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		logger.Warn("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrListUserExample)
		return
	}

	response.Success(c, gin.H{
		"userExamples": data,
	})
}

// List 通过post获取多条记录
// @Summary 获取userExample列表
// @Description 使用post请求获取userExample列表
// @Tags userExample
// @accept json
// @Produce json
// @Param data body types.Params true "查询条件"
// @Success 200 {object} types.Result{}
// @Router /api/v1/userExamples [post]
func (h *userExampleHandler) List(c *gin.Context) {
	form := &types.GetUserExamplesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userExamples, total, err := h.iDao.GetByColumns(c.Request.Context(), &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		logger.Warn("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
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
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return "", 0, true
	}

	return idStr, id, false
}

func convertUserExamples(fromValues []*model.UserExample) ([]*types.GetUserExampleByIDRespond, error) {
	toValues := []*types.GetUserExampleByIDRespond{}
	for _, v := range fromValues {
		data := &types.GetUserExampleByIDRespond{}
		err := copier.Copy(data, v)
		if err != nil {
			return nil, err
		}
		data.ID = utils.Uint64ToStr(v.ID)
		toValues = append(toValues, data)
	}

	return toValues, nil
}
