package handler

import (
	"time"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"
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
// @Param data body CreateUserExampleRequest true "userExample信息"
// @Success 200 {object} Result{}
// @Router /api/v1/userExample [post]
func (h *userExampleHandler) Create(c *gin.Context) {
	form := &CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	err = h.iDao.Create(c.Request.Context(), userExample)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrCreateUserExample)
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
// @Success 200 {object} Result{}
// @Router /api/v1/userExample/{id} [delete]
func (h *userExampleHandler) DeleteByID(c *gin.Context) {
	_, id, isAbout := getUserExampleIDFromPath(c)
	if isAbout {
		return
	}

	err := h.iDao.DeleteByID(c.Request.Context(), id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrDeleteUserExample)
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
// @Param data body UpdateUserExampleByIDRequest true "userExample信息"
// @Success 200 {object} Result{}
// @Router /api/v1/userExample/{id} [put]
func (h *userExampleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbout := getUserExampleIDFromPath(c)
	if isAbout {
		return
	}

	form := &UpdateUserExampleByIDRequest{}
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
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	err = h.iDao.UpdateByID(c.Request.Context(), userExample)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrUpdateUserExample)
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
// @Success 200 {object} Result{}
// @Router /api/v1/userExample/{id} [get]
func (h *userExampleHandler) GetByID(c *gin.Context) {
	idstr, id, isAbout := getUserExampleIDFromPath(c)
	if isAbout {
		return
	}

	userExample, err := h.iDao.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == query.ErrNotFound.Error() {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
			response.Error(c, ecode.ErrGetUserExample)
		}
		return
	}

	data := &GetUserExampleByIDRespond{}
	err = copier.Copy(data, userExample)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("id", id), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InternalServerError)
		return
	}
	data.ID = idstr

	response.Success(c, gin.H{"userExample": data})
}

// List 通过post获取多条记录
// @Summary 获取userExample列表
// @Description 使用post请求获取userExample列表
// @Tags userExample
// @accept json
// @Produce json
// @Param data body Params true "查询条件"
// @Success 200 {object} Result{}
// @Router /api/v1/userExamples [post]
func (h *userExampleHandler) List(c *gin.Context) {
	form := &GetUserExamplesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userExamples, total, err := h.iDao.GetByColumns(c.Request.Context(), &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.ErrListUserExample)
		return
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form), utils.FieldRequestIDFromContext(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	response.Success(c, gin.H{
		"userExamples": data,
		"total":        total,
	})
}

// ----------------------------------- 定义请求参数和返回结果 -----------------------------
// todo generate the request and response struct to here
// delete the templates code start

// CreateUserExampleRequest 创建请求参数，所有字段是必须的，并且满足binding规则
// binding使用说明 https://github.com/go-playground/validator
type CreateUserExampleRequest struct {
	Name     string `json:"name" binding:"min=2"`         // 名称
	Email    string `json:"email" binding:"email"`        // 邮件
	Password string `json:"password" binding:"md5"`       // 密码
	Phone    string `form:"phone" binding:"e164"`         // 手机号码，必须在前加'+86'
	Avatar   string `form:"avatar" binding:"min=5"`       // 头像
	Age      int    `form:"age" binding:"gt=0,lt=120"`    // 年龄
	Gender   int    `form:"gender" binding:"gte=0,lte=2"` // 性别，1:男，2:女
}

// UpdateUserExampleByIDRequest 更新请求参数，所有字段不是必须的，字段为非零值更新
type UpdateUserExampleByIDRequest struct {
	ID       uint64 `json:"id" binding:"-"`      // id
	Name     string `json:"name" binding:""`     // 名称
	Email    string `json:"email" binding:""`    // 邮件
	Password string `json:"password" binding:""` // 密码
	Phone    string `form:"phone" binding:""`    // 手机号码，必须在前加'+86'
	Avatar   string `form:"avatar" binding:""`   // 头像
	Age      int    `form:"age" binding:""`      // 年龄
	Gender   int    `form:"gender" binding:""`   // 性别，1:男，2:女
}

// GetUserExampleByIDRespond 返回数据
type GetUserExampleByIDRespond struct {
	ID        string    `json:"id"`         // id
	Name      string    `json:"name"`       // 名称
	Email     string    `json:"email"`      // 邮件
	Phone     string    `json:"phone"`      // 手机号码
	Avatar    string    `json:"avatar"`     // 头像
	Age       int       `json:"age"`        // 年龄
	Gender    int       `json:"gender"`     // 性别，1:男，2:女
	Status    int       `json:"status"`     // 账号状态
	LoginAt   int64     `json:"login_at"`   // 登录时间戳
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// delete the templates code end

// GetUserExamplesRequest query params
type GetUserExamplesRequest struct {
	query.Params
}

// ListUserExamplesRespond list data
type ListUserExamplesRespond []struct {
	GetUserExampleByIDRespond
}

// ------------------------------- 除了handler的辅助函数 -----------------------------

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

func convertUserExamples(fromValues []*model.UserExample) ([]*GetUserExampleByIDRespond, error) {
	toValues := []*GetUserExampleByIDRespond{}
	for _, v := range fromValues {
		data := &GetUserExampleByIDRespond{}
		err := copier.Copy(data, v)
		if err != nil {
			return nil, err
		}
		data.ID = utils.Uint64ToStr(v.ID)
		toValues = append(toValues, data)
	}

	return toValues, nil
}
