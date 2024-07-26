package query

import "gorm.io/gorm"

var (
	// ErrNotFound record
	// Deprecated: moved to package pkg/gorm/query ErrorNotFound
	ErrNotFound = gorm.ErrRecordNotFound
)
