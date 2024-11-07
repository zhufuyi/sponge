package dbclose

import (
	"database/sql"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestCloseDB(t *testing.T) {
	sqlDB := new(sql.DB)
	checkInUse(sqlDB, time.Millisecond*100)
	checkInUse(sqlDB, time.Millisecond*600)
	db := new(gorm.DB)
	defer func() { recover() }()
	_ = Close(db)
}
