package handler

import (
	"context"

	serverNameExampleV1 "github.com/go-dev-frame/sponge/api/serverNameExample/v1"
	"github.com/go-dev-frame/sponge/internal/service"
)

var _ serverNameExampleV1.{{.TableNameCamel}}Logicer = (*{{.TableNameCamelFCL}}Handler)(nil)

type {{.TableNameCamelFCL}}Handler struct {
	server serverNameExampleV1.{{.TableNameCamel}}Server
}

// New{{.TableNameCamel}}Handler create a handler
func New{{.TableNameCamel}}Handler() serverNameExampleV1.{{.TableNameCamel}}Logicer {
	return &{{.TableNameCamelFCL}}Handler{
		server: service.New{{.TableNameCamel}}Server(),
	}
}

// Create a record
func (h *{{.TableNameCamelFCL}}Handler) Create(ctx context.Context, req *serverNameExampleV1.Create{{.TableNameCamel}}Request) (*serverNameExampleV1.Create{{.TableNameCamel}}Reply, error) {
	return h.server.Create(ctx, req)
}

// DeleteBy{{.ColumnNameCamel}} delete a record by {{.ColumnNameCamelFCL}}
func (h *{{.TableNameCamelFCL}}Handler) DeleteBy{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNameCamel}}Request) (*serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply, error) {
	return h.server.DeleteBy{{.ColumnNameCamel}}(ctx, req)
}

// UpdateBy{{.ColumnNameCamel}} update a record by {{.ColumnNameCamelFCL}}
func (h *{{.TableNameCamelFCL}}Handler) UpdateBy{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Request) (*serverNameExampleV1.Update{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply, error) {
	return h.server.UpdateBy{{.ColumnNameCamel}}(ctx, req)
}

// GetBy{{.ColumnNameCamel}} get a record by {{.ColumnNameCamelFCL}}
func (h *{{.TableNameCamelFCL}}Handler) GetBy{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.Get{{.TableNameCamel}}By{{.ColumnNameCamel}}Request) (*serverNameExampleV1.Get{{.TableNameCamel}}By{{.ColumnNameCamel}}Reply, error) {
	return h.server.GetBy{{.ColumnNameCamel}}(ctx, req)
}

// List of records by query parameters
func (h *{{.TableNameCamelFCL}}Handler) List(ctx context.Context, req *serverNameExampleV1.List{{.TableNameCamel}}Request) (*serverNameExampleV1.List{{.TableNameCamel}}Reply, error) {
	return h.server.List(ctx, req)
}

// DeleteBy{{.ColumnNamePluralCamel}} delete records by batch {{.ColumnNameCamelFCL}}
func (h *{{.TableNameCamelFCL}}Handler) DeleteBy{{.ColumnNamePluralCamel}}(ctx context.Context, req *serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Request) (*serverNameExampleV1.Delete{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Reply, error) {
	return h.server.DeleteBy{{.ColumnNamePluralCamel}}(ctx, req)
}

// GetByCondition get a record by condition
func (h *{{.TableNameCamelFCL}}Handler) GetByCondition(ctx context.Context, req *serverNameExampleV1.Get{{.TableNameCamel}}ByConditionRequest) (*serverNameExampleV1.Get{{.TableNameCamel}}ByConditionReply, error) {
	return h.server.GetByCondition(ctx, req)
}

// ListBy{{.ColumnNamePluralCamel}} list of records by batch {{.ColumnNameCamelFCL}}
func (h *{{.TableNameCamelFCL}}Handler) ListBy{{.ColumnNamePluralCamel}}(ctx context.Context, req *serverNameExampleV1.List{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Request) (*serverNameExampleV1.List{{.TableNameCamel}}By{{.ColumnNamePluralCamel}}Reply, error) {
	return h.server.ListBy{{.ColumnNamePluralCamel}}(ctx, req)
}

// ListByLast{{.ColumnNameCamel}} get records by last {{.ColumnNameCamelFCL}}
func (h *{{.TableNameCamelFCL}}Handler) ListByLast{{.ColumnNameCamel}}(ctx context.Context, req *serverNameExampleV1.List{{.TableNameCamel}}ByLast{{.ColumnNameCamel}}Request) (*serverNameExampleV1.List{{.TableNameCamel}}ByLast{{.ColumnNameCamel}}Reply, error) {
	return h.server.ListByLast{{.ColumnNameCamel}}(ctx, req)
}
