package query

import "gorm.io/gorm"

var (
	// ErrNotFound 空记录
	ErrNotFound = gorm.ErrRecordNotFound
)
