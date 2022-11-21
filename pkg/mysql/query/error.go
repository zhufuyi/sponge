package query

import "gorm.io/gorm"

var (
	// ErrNotFound record
	ErrNotFound = gorm.ErrRecordNotFound
)
