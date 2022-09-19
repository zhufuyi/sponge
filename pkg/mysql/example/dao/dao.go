package dao

import (
	"gorm.io/gorm"
)

// Dao 对象
type Dao struct {
	db *gorm.DB
}

// New 新建dao实例
func New(db *gorm.DB) *Dao {
	return &Dao{
		db: db,
	}
}
