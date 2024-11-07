package handler

import (
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/dao"
	"github.com/zhufuyi/sponge/internal/database"
	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/internal/types"
)

var _ {{.TableNameCamel}}Handler = (*{{.TableNameCamelFCL}}Handler)(nil)

// {{.TableNameCamel}}Handler defining the handler interface
type {{.TableNameCamel}}Handler interface {
	Create(c *gin.Context)
	DeleteBy{{.ColumnNameCamel}}(c *gin.Context)
	UpdateBy{{.ColumnNameCamel}}(c *gin.Context)
	GetBy{{.ColumnNameCamel}}(c *gin.Context)
	List(c *gin.Context)

	DeleteBy{{.ColumnNamePluralCamel}}(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListBy{{.ColumnNamePluralCamel}}(c *gin.Context)
	ListByLast{{.ColumnNameCamel}}(c *gin.Context)
}

type {{.TableNameCamelFCL}}Handler struct {
	iDao dao.{{.TableNameCamel}}Dao
}

// New{{.TableNameCamel}}Handler creating the handler interface
func New{{.TableNameCamel}}Handler() {{.TableNameCamel}}Handler {
	return &{{.TableNameCamelFCL}}Handler{
		iDao: dao.New{{.TableNameCamel}}Dao(
			database.GetDB(),
			cache.New{{.TableNameCamel}}Cache(database.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create {{.TableNameCamelFCL}}
// @Description submit information to create {{.TableNameCamelFCL}}
// @Tags {{.TableNameCamelFCL}}
// @accept json
// @Produce json
// @Param data body types.Create{{.TableNameCamel}}Request true "{{.TableNameCamelFCL}} information"
// @Success 200 {object} types.Create{{.TableNameCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}} [post]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) Create(c *gin.Context) {
	form := &types.Create{{.TableNameCamel}}Request{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	{{.TableNameCamelFCL}} := &model.{{.TableNameCamel}}{}
	err = copier.Copy({{.TableNameCamelFCL}}, form)
	if err != nil {
		response.Error(c, ecode.ErrCreate{{.TableNameCamel}})
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, {{.TableNameCamelFCL}})
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"{{.ColumnNameCamelFCL}}": {{.TableNameCamelFCL}}.{{.ColumnNameCamel}}})
}

// DeleteBy{{.ColumnNameCamel}} delete a record by {{.ColumnNameCamelFCL}}
// @Summary delete {{.TableNameCamelFCL}}
// @Description delete {{.TableNameCamelFCL}} by {{.ColumnNameCamelFCL}}
// @Tags {{.TableNameCamelFCL}}
// @accept json
// @Produce json
// @Param {{.ColumnNameCamelFCL}} path string true "{{.ColumnNameCamelFCL}}"
// @Success 200 {object} types.Delete{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/{{{.ColumnNameCamelFCL}}} [delete]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) DeleteBy{{.ColumnNameCamel}}(c *gin.Context) {
	{{.ColumnNameCamelFCL}}, isAbort := get{{.TableNameCamel}}{{.ColumnNameCamel}}FromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteBy{{.ColumnNameCamel}}(ctx, {{.ColumnNameCamelFCL}})
	if err != nil {
		logger.Error("DeleteBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", {{.ColumnNameCamelFCL}}), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateBy{{.ColumnNameCamel}} update information by {{.ColumnNameCamelFCL}}
// @Summary update {{.TableNameCamelFCL}}
// @Description update {{.TableNameCamelFCL}} information by {{.ColumnNameCamelFCL}}
// @Tags {{.TableNameCamelFCL}}
// @accept json
// @Produce json
// @Param {{.ColumnNameCamelFCL}} path string true "{{.ColumnNameCamelFCL}}"
// @Param data body types.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Request true "{{.TableNameCamelFCL}} information"
// @Success 200 {object} types.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/{{{.ColumnNameCamelFCL}}} [put]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) UpdateBy{{.ColumnNameCamel}}(c *gin.Context) {
	{{.ColumnNameCamelFCL}}, isAbort := get{{.TableNameCamel}}{{.ColumnNameCamel}}FromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Request{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.{{.ColumnNameCamel}} = {{.ColumnNameCamelFCL}}

	{{.TableNameCamelFCL}} := &model.{{.TableNameCamel}}{}
	err = copier.Copy({{.TableNameCamelFCL}}, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateBy{{.ColumnNameCamel}}{{.TableNameCamel}})
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateBy{{.ColumnNameCamel}}(ctx, {{.TableNameCamelFCL}})
	if err != nil {
		logger.Error("UpdateBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetBy{{.ColumnNameCamel}} get a record by {{.ColumnNameCamelFCL}}
// @Summary get {{.TableNameCamelFCL}} detail
// @Description get {{.TableNameCamelFCL}} detail by {{.ColumnNameCamelFCL}}
// @Tags {{.TableNameCamelFCL}}
// @Param {{.ColumnNameCamelFCL}} path string true "{{.ColumnNameCamelFCL}}"
// @Accept json
// @Produce json
// @Success 200 {object} types.Get{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/{{{.ColumnNameCamelFCL}}} [get]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) GetBy{{.ColumnNameCamel}}(c *gin.Context) {
	{{.ColumnNameCamelFCL}}, isAbort := get{{.TableNameCamel}}{{.ColumnNameCamel}}FromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	{{.TableNameCamelFCL}}, err := h.iDao.GetBy{{.ColumnNameCamel}}(ctx, {{.ColumnNameCamelFCL}})
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetBy{{.ColumnNameCamel}} not found", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", {{.ColumnNameCamelFCL}}), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", {{.ColumnNameCamelFCL}}), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.{{.TableNameCamel}}ObjDetail{}
	err = copier.Copy(data, {{.TableNameCamelFCL}})
	if err != nil {
		response.Error(c, ecode.ErrGetBy{{.ColumnNameCamel}}{{.TableNameCamel}})
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"{{.TableNameCamelFCL}}": data})
}

// List of records by query parameters
// @Summary list of {{.TableNamePluralCamelFCL}} by query parameters
// @Description list of {{.TableNamePluralCamelFCL}} by paging and conditions
// @Tags {{.TableNameCamelFCL}}
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.List{{.TableNamePluralCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/list [post]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) List(c *gin.Context) {
	form := &types.List{{.TableNamePluralCamel}}Request{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	{{.TableNamePluralCamelFCL}}, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convert{{.TableNamePluralCamel}}({{.TableNamePluralCamelFCL}})
	if err != nil {
		response.Error(c, ecode.ErrList{{.TableNameCamel}})
		return
	}

	response.Success(c, gin.H{
		"{{.TableNamePluralCamelFCL}}": data,
		"total":        total,
	})
}

// DeleteBy{{.ColumnNamePluralCamel}} delete records by batch {{.ColumnNameCamelFCL}}
// @Summary delete {{.TableNamePluralCamelFCL}}
// @Description delete {{.TableNamePluralCamelFCL}} by batch {{.ColumnNameCamelFCL}}
// @Tags {{.TableNameCamelFCL}}
// @Param data body types.Delete{{.TableNamePluralCamel}}By{{.ColumnNamePluralCamel}}Request true "{{.ColumnNameCamelFCL}} array"
// @Accept json
// @Produce json
// @Success 200 {object} types.Delete{{.TableNamePluralCamel}}By{{.ColumnNamePluralCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/delete/{{.ColumnNamePluralCamelFCL}} [post]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) DeleteBy{{.ColumnNamePluralCamel}}(c *gin.Context) {
	form := &types.Delete{{.TableNamePluralCamel}}By{{.ColumnNamePluralCamel}}Request{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.DeleteBy{{.ColumnNamePluralCamel}}(ctx, form.{{.ColumnNamePluralCamel}})
	if err != nil {
		logger.Error("GetBy{{.ColumnNamePluralCamel}} error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByCondition get a record by condition
// @Summary get {{.TableNameCamelFCL}} by condition
// @Description get {{.TableNameCamelFCL}} by condition
// @Tags {{.TableNameCamelFCL}}
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.Get{{.TableNameCamel}}ByConditionReply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/condition [post]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) GetByCondition(c *gin.Context) {
	form := &types.Get{{.TableNameCamel}}ByConditionRequest{}
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

	ctx := middleware.WrapCtx(c)
	{{.TableNameCamelFCL}}, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.{{.TableNameCamel}}ObjDetail{}
	err = copier.Copy(data, {{.TableNameCamelFCL}})
	if err != nil {
		response.Error(c, ecode.ErrGetBy{{.ColumnNameCamel}}{{.TableNameCamel}})
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"{{.TableNameCamelFCL}}": data})
}

// ListBy{{.ColumnNamePluralCamel}} list of records by batch {{.ColumnNameCamelFCL}}
// @Summary list of {{.TableNamePluralCamelFCL}} by batch {{.ColumnNameCamelFCL}}
// @Description list of {{.TableNamePluralCamelFCL}} by batch {{.ColumnNameCamelFCL}}
// @Tags {{.TableNameCamelFCL}}
// @Param data body types.List{{.TableNamePluralCamel}}By{{.ColumnNamePluralCamel}}Request true "{{.ColumnNameCamelFCL}} array"
// @Accept json
// @Produce json
// @Success 200 {object} types.List{{.TableNamePluralCamel}}By{{.ColumnNamePluralCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/list/{{.ColumnNamePluralCamelFCL}} [post]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) ListBy{{.ColumnNamePluralCamel}}(c *gin.Context) {
	form := &types.List{{.TableNamePluralCamel}}By{{.ColumnNamePluralCamel}}Request{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	{{.TableNameCamelFCL}}Map, err := h.iDao.GetBy{{.ColumnNamePluralCamel}}(ctx, form.{{.ColumnNamePluralCamel}})
	if err != nil {
		logger.Error("GetBy{{.ColumnNamePluralCamel}} error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	{{.TableNamePluralCamelFCL}} := []*types.{{.TableNameCamel}}ObjDetail{}
	for _, {{.ColumnNameCamelFCL}} := range form.{{.ColumnNamePluralCamel}} {
		if v, ok := {{.TableNameCamelFCL}}Map[{{.ColumnNameCamelFCL}}]; ok {
			record, err := convert{{.TableNameCamel}}(v)
			if err != nil {
				response.Error(c, ecode.ErrList{{.TableNameCamel}})
				return
			}
			{{.TableNamePluralCamelFCL}} = append({{.TableNamePluralCamelFCL}}, record)
		}
	}

	response.Success(c, gin.H{
		"{{.TableNamePluralCamelFCL}}": {{.TableNamePluralCamelFCL}},
	})
}

// ListByLast{{.ColumnNameCamel}} get records by last {{.ColumnNameCamelFCL}} and limit
// @Summary list of {{.TableNamePluralCamelFCL}} by last {{.ColumnNameCamelFCL}} and limit
// @Description list of {{.TableNamePluralCamelFCL}} by last {{.ColumnNameCamelFCL}} and limit
// @Tags {{.TableNameCamelFCL}}
// @accept json
// @Produce json
// @Param last{{.ColumnNameCamel}} query {{if .IsStringType}}string{{else}}int{{end}} false "last {{.ColumnNameCamelFCL}}, default is {{if .IsStringType}}Max string {{.ColumnNameCamelFCL}}" default(){{else}}MaxInt32" default(0){{end}}
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-{{.ColumnName}})
// @Success 200 {object} types.List{{.TableNamePluralCamel}}Reply{}
// @Router /api/v1/{{.TableNameCamelFCL}}/list [get]
// @Security BearerAuth
func (h *{{.TableNameCamelFCL}}Handler) ListByLast{{.ColumnNameCamel}}(c *gin.Context) {
	{{if .IsStringType}}last{{.ColumnNameCamel}} := c.Query("last{{.ColumnNameCamel}}")
	if last{{.ColumnNameCamel}} == "" {
		last{{.ColumnNameCamel}} = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	}{{else}}last{{.ColumnNameCamel}} := utils.StrTo{{.GoTypeFCU}}(c.Query("last{{.ColumnNameCamel}}"))
	if last{{.ColumnNameCamel}} == 0 {
		last{{.ColumnNameCamel}} = math.MaxInt32
	}{{end}}
	limit := utils.StrToInt(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	sort := c.Query("sort")

	ctx := middleware.WrapCtx(c)
	{{.TableNamePluralCamelFCL}}, err := h.iDao.GetByLast{{.ColumnNameCamel}}(ctx, last{{.ColumnNameCamel}}, limit, sort)
	if err != nil {
		logger.Error("GetByLast{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("last{{.ColumnNameCamel}}", last{{.ColumnNameCamel}}), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convert{{.TableNamePluralCamel}}({{.TableNamePluralCamelFCL}})
	if err != nil {
		response.Error(c, ecode.ErrListByLast{{.ColumnNameCamel}}{{.TableNameCamel}})
		return
	}

	response.Success(c, gin.H{
		"{{.TableNamePluralCamelFCL}}": data,
	})
}

func get{{.TableNameCamel}}{{.ColumnNameCamel}}FromPath(c *gin.Context) ({{.GoType}}, bool) {
	{{.ColumnNameCamelFCL}}Str := c.Param("{{.ColumnNameCamelFCL}}")
{{if .IsStringType}}
	if {{.ColumnNameCamelFCL}}Str == "" {
		logger.Warn("{{.ColumnNameCamelFCL}} is empty", middleware.GCtxRequestIDField(c))
		return "", true
	}
	return {{.ColumnNameCamelFCL}}Str, false
{{else}}
	{{.ColumnNameCamelFCL}}, err := utils.StrTo{{.GoTypeFCU}}E({{.ColumnNameCamelFCL}}Str)
	if err != nil || {{.ColumnNameCamelFCL}}Str == "" {
		logger.Warn("StrTo{{.GoTypeFCU}}E error: ", logger.String("{{.ColumnNameCamelFCL}}Str", {{.ColumnNameCamelFCL}}Str), middleware.GCtxRequestIDField(c))
		return 0, true
	}
	return {{.ColumnNameCamelFCL}}, false
{{end}}
}

func convert{{.TableNameCamel}}({{.TableNameCamelFCL}} *model.{{.TableNameCamel}}) (*types.{{.TableNameCamel}}ObjDetail, error) {
	data := &types.{{.TableNameCamel}}ObjDetail{}
	err := copier.Copy(data, {{.TableNameCamelFCL}})
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convert{{.TableNamePluralCamel}}(fromValues []*model.{{.TableNameCamel}}) ([]*types.{{.TableNameCamel}}ObjDetail, error) {
	toValues := []*types.{{.TableNameCamel}}ObjDetail{}
	for _, v := range fromValues {
		data, err := convert{{.TableNameCamel}}(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
