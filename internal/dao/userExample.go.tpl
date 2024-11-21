package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/sgorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/zhufuyi/sponge/internal/cache"
	"github.com/zhufuyi/sponge/internal/database"
	"github.com/zhufuyi/sponge/internal/model"
)

var _ {{.TableNameCamel}}Dao = (*{{.TableNameCamelFCL}}Dao)(nil)

// {{.TableNameCamel}}Dao defining the dao interface
type {{.TableNameCamel}}Dao interface {
	Create(ctx context.Context, table *model.{{.TableNameCamel}}) error
	DeleteBy{{.ColumnNameCamel}}(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error
	UpdateBy{{.ColumnNameCamel}}(ctx context.Context, table *model.{{.TableNameCamel}}) error
	GetBy{{.ColumnNameCamel}}(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) (*model.{{.TableNameCamel}}, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.{{.TableNameCamel}}, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.{{.TableNameCamel}}) ({{.GoType}}, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, {{.ColumnNameCamelFCL}} {{.GoType}}) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.{{.TableNameCamel}}) error
}

type {{.TableNameCamelFCL}}Dao struct {
	db    *gorm.DB
	cache cache.{{.TableNameCamel}}Cache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

// New{{.TableNameCamel}}Dao creating the dao interface
func New{{.TableNameCamel}}Dao(db *gorm.DB, xCache cache.{{.TableNameCamel}}Cache) {{.TableNameCamel}}Dao {
	if xCache == nil {
		return &{{.TableNameCamelFCL}}Dao{db: db}
	}
	return &{{.TableNameCamelFCL}}Dao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *{{.TableNameCamelFCL}}Dao) deleteCache(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error {
	if d.cache != nil {
		return d.cache.Del(ctx, {{.ColumnNameCamelFCL}})
	}
	return nil
}

// Create a record, insert the record and the {{.ColumnNameCamelFCL}} value is written back to the table
func (d *{{.TableNameCamelFCL}}Dao) Create(ctx context.Context, table *model.{{.TableNameCamel}}) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteBy{{.ColumnNameCamel}} delete a record by {{.ColumnNameCamelFCL}}
func (d *{{.TableNameCamelFCL}}Dao) DeleteBy{{.ColumnNameCamel}}(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) error {
	err := d.db.WithContext(ctx).Where("{{.ColumnName}} = ?", {{.ColumnNameCamelFCL}}).Delete(&model.{{.TableNameCamel}}{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, {{.ColumnNameCamelFCL}})

	return nil
}

// UpdateBy{{.ColumnNameCamel}} update a record by {{.ColumnNameCamelFCL}}
func (d *{{.TableNameCamelFCL}}Dao) UpdateBy{{.ColumnNameCamel}}(ctx context.Context, table *model.{{.TableNameCamel}}) error {
	err := d.updateDataBy{{.ColumnNameCamel}}(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.{{.ColumnNameCamel}})

	return err
}

func (d *{{.TableNameCamelFCL}}Dao) updateDataBy{{.ColumnNameCamel}}(ctx context.Context, db *gorm.DB, table *model.{{.TableNameCamel}}) error {
	{{if .IsStringType}}if table.{{.ColumnNameCamel}} == "" {
		return errors.New("{{.ColumnNameCamelFCL}} cannot be empty")
	}
{{else}}	if table.{{.ColumnNameCamel}} < 1 {
		return errors.New("{{.ColumnNameCamelFCL}} cannot be 0")
	}
{{end}}

	update := map[string]interface{}{}
	// todo generate the update fields code to here
	// delete the templates code start
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Password != "" {
		update["password"] = table.Password
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Phone != "" {
		update["phone"] = table.Phone
	}
	if table.Avatar != "" {
		update["avatar"] = table.Avatar
	}
	if table.Age > 0 {
		update["age"] = table.Age
	}
	if table.Gender > 0 {
		update["gender"] = table.Gender
	}
	if table.LoginAt > 0 {
		update["login_at"] = table.LoginAt
	}
	// delete the templates code end

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetBy{{.ColumnNameCamel}} get a record by {{.ColumnNameCamelFCL}}
func (d *{{.TableNameCamelFCL}}Dao) GetBy{{.ColumnNameCamel}}(ctx context.Context, {{.ColumnNameCamelFCL}} {{.GoType}}) (*model.{{.TableNameCamel}}, error) {
	// no cache
	if d.cache == nil {
		record := &model.{{.TableNameCamel}}{}
		err := d.db.WithContext(ctx).Where("{{.ColumnName}} = ?", {{.ColumnNameCamelFCL}}).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, {{.ColumnNameCamelFCL}})
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same {{.ColumnNameCamelFCL}}, prevent high concurrent simultaneous access to database
		{{if .IsStringType}}val, err, _ := d.sfg.Do({{.ColumnNameCamelFCL}}, func() (interface{}, error) {
{{else}}		val, err, _ := d.sfg.Do(utils.{{.GoTypeFCU}}ToStr({{.ColumnNameCamelFCL}}), func() (interface{}, error) {
{{end}}
			table := &model.{{.TableNameCamel}}{}
			err := d.db.WithContext(ctx).Where("{{.ColumnName}} = ?", {{.ColumnNameCamelFCL}}).First(table).Error
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, {{.ColumnNameCamelFCL}}); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", {{.ColumnNameCamelFCL}}))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, {{.ColumnNameCamelFCL}}, table, cache.{{.TableNameCamel}}ExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("{{.ColumnNameCamelFCL}}", {{.ColumnNameCamelFCL}}))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.{{.TableNameCamel}})
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get paging records by column information,
// Note: query performance degrades when table rows are very large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	limit: lines per page
//	sort: sort fields, default is {{.ColumnNameCamelFCL}} backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Limit: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *{{.TableNameCamelFCL}}Dao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.{{.TableNameCamel}}, int64, error) {
	if params.Sort == "" {
		params.Sort = "-{{.ColumnName}}"
	}
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.{{.TableNameCamel}}{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.{{.TableNameCamel}}{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *{{.TableNameCamelFCL}}Dao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.{{.TableNameCamel}}) ({{.GoType}}, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.{{.ColumnNameCamel}}, err
}

// DeleteByTx delete a record by {{.ColumnNameCamelFCL}} in the database using the provided transaction
func (d *{{.TableNameCamelFCL}}Dao) DeleteByTx(ctx context.Context, tx *gorm.DB, {{.ColumnNameCamelFCL}} {{.GoType}}) error {
	err := tx.WithContext(ctx).Where("{{.ColumnName}} = ?", {{.ColumnNameCamelFCL}}).Delete(&model.{{.TableNameCamel}}{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, {{.ColumnNameCamelFCL}})

	return nil
}

// UpdateByTx update a record by {{.ColumnNameCamelFCL}} in the database using the provided transaction
func (d *{{.TableNameCamelFCL}}Dao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.{{.TableNameCamel}}) error {
	err := d.updateDataBy{{.ColumnNameCamel}}(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.{{.ColumnNameCamel}})

	return err
}
