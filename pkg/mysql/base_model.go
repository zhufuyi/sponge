package mysql

import (
	"reflect"
	"time"

	"github.com/huandu/xstrings"
	"gorm.io/gorm"
)

//type Model = gorm.Model

// Model 表内嵌字段，嵌入到其他表结构体时添加 `gorm:"embedded"`
type Model struct {
	ID        uint64         `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// KV map类型
type KV = map[string]interface{}

// GetTableName 获取表名
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
