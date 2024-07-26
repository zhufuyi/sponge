package mysql

import (
	"reflect"
	"time"

	"github.com/huandu/xstrings"
	"gorm.io/gorm"
)

//type Model = gorm.Model

// Model embedded structs, add `gorm: "embedded"` when defining table structs
// Deprecated: moved to package pkg/ggorm Model
type Model struct {
	ID        uint64         `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// Model2 embedded structs, json tag named is snake case
// Deprecated: moved to package pkg/ggorm Model2
type Model2 struct {
	ID        uint64         `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// KV map type
// Deprecated: moved to package pkg/ggorm KV
type KV = map[string]interface{}

// GetTableName get table name
// Deprecated: moved to package pkg/ggorm GetTableName
func GetTableName(object interface{}) string {
	tableName := ""

	typeof := reflect.TypeOf(object)
	switch typeof.Kind() {
	case reflect.Ptr:
		tableName = typeof.Elem().Name()
	case reflect.Struct:
		tableName = typeof.Name()
	default:
		return tableName
	}

	return xstrings.ToSnakeCase(tableName)
}
