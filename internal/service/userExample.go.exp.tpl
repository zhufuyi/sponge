package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/grpc/interceptor"
	"github.com/go-dev-frame/sponge/pkg/logger"

	serverNameExampleV1 "github.com/go-dev-frame/sponge/api/serverNameExample/v1"
	"github.com/go-dev-frame/sponge/internal/cache"
	"github.com/go-dev-frame/sponge/internal/dao"
	"github.com/go-dev-frame/sponge/internal/database"
	"github.com/go-dev-frame/sponge/internal/ecode"
	"github.com/go-dev-frame/sponge/internal/model"
)

func init() {
	registerFns = append(registerFns, func(server *grpc.Server) {
		serverNameExampleV1.Register{{.TableNameCamel}}Server(server, New{{.TableNameCamel}}Server()) // register service to the rpc service
	})
}

var _ serverNameExampleV1.{{.TableNameCamel}}Server = (*{{.TableNameCamelFCL}})(nil)
var _ time.Time

type {{.TableNameCamelFCL}} struct {
	serverNameExampleV1.Unimplemented{{.TableNameCamel}}Server

	iDao dao.{{.TableNameCamel}}Dao
}

// New{{.TableNameCamel}}Server create a new service
func New{{.TableNameCamel}}Server() serverNameExampleV1.{{.TableNameCamel}}Server {
	return &{{.TableNameCamelFCL}}{
		iDao: dao.New{{.TableNameCamel}}Dao(
			database.GetDB(), // todo show db driver name here
			cache.New{{.TableNameCamel}}Cache(database.GetCacheType()),
		),
	}
}

// Create a record
func (s *{{.TableNameCamelFCL}}) Create(ctx context.Context, req *serverNameExampleV1.Create{{.TableNameCamel}}Request) (*serverNameExampleV1.Create{{.TableNameCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	record := &model.{{.TableNameCamel}}{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusCreate{{.TableNameCamel}}.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	err = s.iDao.Create(ctx, record)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("{{.TableNameCamelFCL}}", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.Create{{.TableNameCamel}}Reply{ {{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}: record.{{.ColumnNameCamel}} }, nil
}

// DeleteBy{{.ColumnNameCamel}} delete a record by {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) DeleteBy{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNameCamel}}Request) (*serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	err = s.iDao.DeleteBy{{.ColumnNameCamel}}(ctx, req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}})
	if err != nil {
		logger.Error("DeleteBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply{}, nil
}

// UpdateBy{{.ColumnNameCamel}} update a record by {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) UpdateBy{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Request) (*serverNameExampleV1.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	record := &model.{{.TableNameCamel}}{}
	err = copier.Copy(record, req)
	if err != nil {
		return nil, ecode.StatusUpdateBy{{.ColumnNameCamel}}{{.TableNameCamel}}.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	record.{{.ColumnNameCamel}} = req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}

	err = s.iDao.UpdateBy{{.ColumnNameCamel}}(ctx, record)
	if err != nil {
		logger.Error("UpdateBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.TableNameCamelFCL}}", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply{}, nil
}

// GetBy{{.ColumnNameCamel}} get a record by {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) GetBy{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.Get{{.TableNameCamel}}By{{.ColumnNameCamel}}Request) (*serverNameExampleV1.Get{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	record, err := s.iDao.GetBy{{.ColumnNameCamel}}(ctx, req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}})
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("GetBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convert{{.TableNameCamel}}(record)
	if err != nil {
		logger.Warn("convert{{.TableNameCamel}} error", logger.Err(err), logger.Any("{{.TableNameCamelFCL}}", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetBy{{.ColumnNameCamel}}{{.TableNameCamel}}.Err()
	}

	return &serverNameExampleV1.Get{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply{ {{.TableNameCamel}}: data}, nil
}

// List of records by query parameters
func (s *{{.TableNameCamelFCL}}) List(ctx context.Context, req *serverNameExampleV1.List{{.TableNameCamel}}Request) (*serverNameExampleV1.List{{.TableNameCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	params := &query.Params{}
	err = copier.Copy(params, req.Params)
	if err != nil {
		return nil, ecode.StatusList{{.TableNameCamel}}.Err()
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	records, total, err := s.iDao.GetByColumns(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "query params error:") {
			logger.Warn("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusInvalidParams.Err()
		}
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("params", params), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	{{.TableNamePluralCamelFCL}} := []*serverNameExampleV1.{{.TableNameCamel}}{}
	for _, record := range records {
		data, err := convert{{.TableNameCamel}}(record)
		if err != nil {
			logger.Warn("convert{{.TableNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", record.{{.ColumnNameCamel}}), interceptor.ServerCtxRequestIDField(ctx))
			continue
		}
		{{.TableNamePluralCamelFCL}} = append({{.TableNamePluralCamelFCL}}, data)
	}

	return &serverNameExampleV1.List{{.TableNameCamel}}Reply{
		Total:        total,
		{{.TableNamePluralCamel}}: {{.TableNamePluralCamelFCL}},
	}, nil
}

// DeleteBy{{.ColumnNamePluralCamel}} delete records by batch {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) DeleteBy{{.ColumnNamePluralCamel}}(ctx context.Context, req *serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Request) (*serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	err = s.iDao.DeleteBy{{.ColumnNamePluralCamel}}(ctx, req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}s)
	if err != nil {
		logger.Error("DeleteBy{{.ColumnNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNamePluralCamelFCL}}", req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}s), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	return &serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Reply{}, nil
}

// GetByCondition get a record by {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) GetByCondition(ctx context.Context, req *serverNameExampleV1.Get{{.TableNameCamel}}ByConditionRequest) (*serverNameExampleV1.Get{{.TableNameCamel}}ByConditionReply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	conditions := &query.Conditions{}
	for _, v := range req.Conditions.GetColumns() {
		column := query.Column{}
		_ = copier.Copy(&column, v)
		conditions.Columns = append(conditions.Columns, column)
	}
	err = conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error", logger.Err(err), logger.Any("conditions", conditions), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}

	record, err := s.iDao.GetByCondition(ctx, conditions)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByCondition error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
			return nil, ecode.StatusNotFound.Err()
		}
		logger.Error("GetByCondition error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	data, err := convert{{.TableNameCamel}}(record)
	if err != nil {
		logger.Warn("convert{{.TableNameCamel}} error", logger.Err(err), logger.Any("{{.TableNameCamelFCL}}", record), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusGetByCondition{{.TableNameCamel}}.Err()
	}

	return &serverNameExampleV1.Get{{.TableNameCamel}}ByConditionReply{
		{{.TableNameCamel}}: data,
	}, nil
}

// ListBy{{.ColumnNamePluralCamel}} list of records by batch {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) ListBy{{.ColumnNamePluralCamel}}(ctx context.Context, req *serverNameExampleV1.List{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Request) (*serverNameExampleV1.List{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	ctx = interceptor.WrapServerCtx(ctx)

	{{.TableNameCamelFCL}}Map, err := s.iDao.GetBy{{.ColumnNamePluralCamel}}(ctx, req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}s)
	if err != nil {
		logger.Error("GetBy{{.ColumnNamePluralCamel}} error", logger.Err(err), logger.Any("{{.ColumnNamePluralCamelFCL}}", req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}s), interceptor.ServerCtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	{{.TableNamePluralCamelFCL}} := []*serverNameExampleV1.{{.TableNameCamel}}{}
	for _, {{.ColumnNameCamelFCL}} := range req.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}}s {
		if v, ok := {{.TableNameCamelFCL}}Map[{{.ColumnNameCamelFCL}}]; ok {
			record, err := convert{{.TableNameCamel}}(v)
			if err != nil {
				logger.Warn("convert{{.TableNameCamel}} error", logger.Err(err), logger.Any("{{.TableNameCamelFCL}}", v), interceptor.ServerCtxRequestIDField(ctx))
				return nil, ecode.StatusInternalServerError.ToRPCErr()
			}
			{{.TableNamePluralCamelFCL}} = append({{.TableNamePluralCamelFCL}}, record)
		}
	}

	return &serverNameExampleV1.List{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Reply{ {{.TableNamePluralCamel}}: {{.TableNamePluralCamelFCL}} }, nil
}

// ListByLast{{.ColumnNameCamel}} list {{.TableNameCamelFCL}} by last {{.ColumnNameCamelFCL}}
func (s *{{.TableNameCamelFCL}}) ListByLast{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.List{{.TableNameCamel}}ByLast{{.ColumnNameCamel}}Request) (*serverNameExampleV1.List{{.TableNameCamel}}ByLast{{.ColumnNameCamel}}Reply, error) {
	err := req.Validate()
	if err != nil {
		logger.Warn("req.Validate error", logger.Err(err), logger.Any("req", req), interceptor.CtxRequestIDField(ctx))
		return nil, ecode.StatusInvalidParams.Err()
	}
	{{if .IsStringType}}if req.Last{{.ColumnNameCamel}} == "" {
		req.Last{{.ColumnNameCamel}} = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	}{{else}}
	if req.Last{{.ColumnNameCamel}} == 0 {
		req.Last{{.ColumnNameCamel}} = math.MaxInt32
	}{{end}}
	if req.Limit == 0 {
		req.Limit = 10
	}

	records, err := s.iDao.GetByLast{{.ColumnNameCamel}}(ctx, req.Last{{.ColumnNameCamel}}, int(req.Limit), req.Sort)
	if err != nil {
		logger.Error("ListByLast{{.ColumnNameCamel}} error", logger.Err(err), interceptor.CtxRequestIDField(ctx))
		return nil, ecode.StatusInternalServerError.ToRPCErr()
	}

	{{.TableNamePluralCamelFCL}} := []*serverNameExampleV1.{{.TableNameCamel}}{}
	for _, record := range records {
		data, err := convert{{.TableNameCamel}}(record)
		if err != nil {
			logger.Warn("convert{{.TableNameCamel}} error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", record.{{.ColumnNameCamel}}), interceptor.ServerCtxRequestIDField(ctx))
			continue
		}
		{{.TableNamePluralCamelFCL}} = append({{.TableNamePluralCamelFCL}}, data)
	}

	return &serverNameExampleV1.List{{.TableNameCamel}}ByLast{{.ColumnNameCamel}}Reply{
		{{.TableNamePluralCamel}}: {{.TableNamePluralCamelFCL}},
	}, nil
}

func convert{{.TableNameCamel}}(record *model.{{.TableNameCamel}}) (*serverNameExampleV1.{{.TableNameCamel}}, error) {
	value := &serverNameExampleV1.{{.TableNameCamel}}{}
	err := copier.Copy(value, record)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here, e.g. CreatedAt, UpdatedAt
	value.{{if .IsStandardPrimaryKey}}Id{{else}}{{.ColumnNameCamel}}{{end}} = record.{{.ColumnNameCamel}}
	// todo generate the conversion createdAt and updatedAt code here
	// delete the templates code start
	value.CreatedAt = record.CreatedAt.Format(time.RFC3339)
	value.UpdatedAt = record.UpdatedAt.Format(time.RFC3339)
	// delete the templates code end
	return value, nil
}
