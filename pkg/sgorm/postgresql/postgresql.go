// Package postgresql provides a gorm driver for postgresql.
package postgresql

import (
	"fmt"
	"log"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/zhufuyi/sponge/pkg/sgorm/dbclose"
	"github.com/zhufuyi/sponge/pkg/sgorm/glog"
)

// Init postgresql
func Init(dsn string, opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	o.apply(opts...)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig(o))
	if err != nil {
		return nil, err
	}

	// register trace plugin
	if o.enableTrace {
		err = db.Use(otelgorm.NewPlugin())
		if err != nil {
			return nil, fmt.Errorf("using gorm opentelemetry, err: %v", err)
		}
	}

	// register plugins
	for _, plugin := range o.plugins {
		err = db.Use(plugin)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// gorm setting
func gormConfig(o *options) *gorm.Config {
	config := &gorm.Config{
		// disable foreign key constraints, not recommended for production environments
		DisableForeignKeyConstraintWhenMigrating: o.disableForeignKey,
		// removing the plural of an epithet
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	// print SQL
	if o.isLog {
		if o.gLog == nil {
			config.Logger = logger.Default.LogMode(o.logLevel)
		} else {
			config.Logger = glog.NewCustomGormLogger(o.gLog, o.requestIDKey, o.logLevel)
		}
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// print only slow queries
	if o.slowThreshold > 0 {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // use the standard output asWriter
			logger.Config{
				SlowThreshold: o.slowThreshold,
				Colorful:      true,
				LogLevel:      logger.Warn, // set the logging level, only above the specified level will output the slow query log
			},
		)
	}

	return config
}

// Close close gorm db
func Close(db *gorm.DB) error {
	return dbclose.Close(db)
}
