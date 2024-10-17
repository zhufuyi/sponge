// Package ggorm is a library wrapped on top of gorm.io/gorm, with added features such as link tracing, paging queries, etc.
package ggorm

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	// DBDriverMysql mysql driver
	DBDriverMysql = "mysql"
	// DBDriverPostgresql postgresql driver
	DBDriverPostgresql = "postgresql"
	// DBDriverTidb tidb driver
	DBDriverTidb = "tidb"
	// DBDriverSqlite sqlite driver
	DBDriverSqlite = "sqlite"
)

// InitMysql init mysql or tidb
func InitMysql(dsn string, opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	o.apply(opts...)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(o.maxIdleConns)       // set the maximum number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(o.maxOpenConns)       // set the maximum number of open database connections
	sqlDB.SetConnMaxLifetime(o.connMaxLifetime) // set the maximum time a connection can be reused

	db, err := gorm.Open(mysqlDriver.New(mysqlDriver.Config{Conn: sqlDB}), gormConfig(o))
	if err != nil {
		return nil, err
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4") // automatic appending of table suffixes when creating tables

	// register trace plugin
	if o.enableTrace {
		err = db.Use(otelgorm.NewPlugin())
		if err != nil {
			return nil, fmt.Errorf("using gorm opentelemetry, err: %v", err)
		}
	}

	// register read-write separation plugin
	if len(o.slavesDsn) > 0 {
		err = db.Use(rwSeparationPlugin(o))
		if err != nil {
			return nil, err
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

// InitPostgresql init postgresql
func InitPostgresql(dsn string, schemaName string, opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	o.apply(opts...)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig(o))
	if err != nil {
		return nil, err
	}

	// Set schema search path if a schema is provided
	if schemaName != "" {
		err = db.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error
		if err != nil {
			return nil, fmt.Errorf("setting search path to schema %s, err: %v", schemaName, err)
		}
	}

	// register trace plugin
	if o.enableTrace {
		err = db.Use(otelgorm.NewPlugin())
		if err != nil {
			return nil, fmt.Errorf("using gorm opentelemetry, err: %v", err)
		}
	}

	// register read-write separation plugin
	if len(o.slavesDsn) > 0 {
		err = db.Use(rwSeparationPlugin(o))
		if err != nil {
			return nil, err
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

// InitTidb init tidb
func InitTidb(dsn string, opts ...Option) (*gorm.DB, error) {
	return InitMysql(dsn, opts...)
}

// InitSqlite init sqlite
func InitSqlite(dbFile string, opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	o.apply(opts...)

	dsn := fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental", dbFile)
	db, err := gorm.Open(sqlite.Open(dsn), gormConfig(o))
	if err != nil {
		return nil, err
	}
	db.Set("gorm:auto_increment", true)

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

// CloseDB close gorm db
func CloseDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	checkInUse(sqlDB, time.Second*5)

	return sqlDB.Close()
}

func checkInUse(sqlDB *sql.DB, duration time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), duration) //nolint
	for {
		select {
		case <-time.After(time.Millisecond * 250):
			if v := sqlDB.Stats().InUse; v == 0 {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// CloseSQLDB close sql db
func CloseSQLDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
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
			config.Logger = NewCustomGormLogger(o)
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

func rwSeparationPlugin(o *options) gorm.Plugin {
	slaves := []gorm.Dialector{}
	for _, dsn := range o.slavesDsn {
		slaves = append(slaves, mysqlDriver.New(mysqlDriver.Config{
			DSN: dsn,
		}))
	}

	masters := []gorm.Dialector{}
	for _, dsn := range o.mastersDsn {
		masters = append(masters, mysqlDriver.New(mysqlDriver.Config{
			DSN: dsn,
		}))
	}

	return dbresolver.Register(dbresolver.Config{
		Sources:  masters,
		Replicas: slaves,
		Policy:   dbresolver.RandomPolicy{},
	})
}
