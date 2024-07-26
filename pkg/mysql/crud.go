package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
)

// TableName get table name
// Deprecated: moved to package pkg/ggorm TableName
func TableName(table interface{}) string {
	return GetTableName(table)
}

// Create a new record
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm Create
func Create(ctx context.Context, db *gorm.DB, table interface{}) error {
	return db.WithContext(ctx).Create(table).Error
}

// Delete record
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm Delete
func Delete(ctx context.Context, db *gorm.DB, table interface{}, queryCondition interface{}, args ...interface{}) error {
	return db.WithContext(ctx).Where(queryCondition, args...).Delete(table).Error
}

// DeleteByID delete record by id
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm DeleteByID
func DeleteByID(ctx context.Context, db *gorm.DB, table interface{}, id interface{}) error {
	return db.WithContext(ctx).Where("id = ?", id).Delete(table).Error
}

// Update record
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm Update
func Update(ctx context.Context, db *gorm.DB, table interface{}, column string, value interface{}, queryCondition interface{}, args ...interface{}) error {
	return db.WithContext(ctx).Model(table).Where(queryCondition, args...).Update(column, value).Error
}

// Updates record
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm Updates
func Updates(ctx context.Context, db *gorm.DB, table interface{}, update KV, queryCondition interface{}, args ...interface{}) error {
	return db.WithContext(ctx).Model(table).Where(queryCondition, args...).Updates(update).Error
}

// Get one record
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm Get
func Get(ctx context.Context, db *gorm.DB, table interface{}, queryCondition interface{}, args ...interface{}) error {
	return db.WithContext(ctx).Where(queryCondition, args...).First(table).Error
}

// GetByID get record by id
// Deprecated: moved to package pkg/ggorm GetByID
func GetByID(ctx context.Context, db *gorm.DB, table interface{}, id interface{}) error {
	return db.WithContext(ctx).Where("id = ?", id).First(table).Error
}

// List multiple records, starting from page 0
// the param of 'tables' must be a slice, eg: []StructName
// Deprecated: moved to package pkg/ggorm List
func List(ctx context.Context, db *gorm.DB, tables interface{}, page *query.Page, queryCondition interface{}, args ...interface{}) error {
	return db.WithContext(ctx).Order(page.Sort()).Limit(page.Size()).Offset(page.Offset()).Where(queryCondition, args...).Find(tables).Error
}

// Count number of records
// the param of 'table' must be pointer, eg: &StructName
// Deprecated: moved to package pkg/ggorm Count
func Count(ctx context.Context, db *gorm.DB, table interface{}, queryCondition interface{}, args ...interface{}) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(table).Where(queryCondition, args...).Count(&count).Error
	return count, err
}
