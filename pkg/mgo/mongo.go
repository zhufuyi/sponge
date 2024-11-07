// Package mgo is a library wrapped on go.mongodb.org/mongo-driver/mongo, with added features paging queries, etc.
package mgo

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DB = mongo.Database

var ErrNoDocuments = mongo.ErrNoDocuments

const (
	// DBDriverName mongodb driver
	DBDriverName = "mongodb"
)

// Init connecting to mongo
func Init(dsn string, opts ...*options.ClientOptions) (*mongo.Database, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	dbName := strings.TrimLeft(u.Path, "/")
	if dbName == "" {
		return nil, errors.New("database name is empty")
	}

	var uri string
	if u.RawQuery == "" {
		uri = strings.TrimRight(dsn, u.Path)
	} else {
		tmp := strings.TrimRight(dsn, u.RawQuery)
		uri = strings.TrimRight(tmp, dbName+"?") + "?" + u.RawQuery
	}

	return Init2(uri, dbName, opts...)
}

// Init2 connecting to mongo using uri
func Init2(uri string, dbName string, opts ...*options.ClientOptions) (*mongo.Database, error) {
	ctx := context.Background()
	mongoOpts := []*options.ClientOptions{
		options.Client().ApplyURI(uri),
	}
	mongoOpts = append(mongoOpts, opts...)
	client, err := mongo.Connect(ctx, mongoOpts...)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return db, nil
}

// Close mongodb
func Close(db *mongo.Database) error {
	return db.Client().Disconnect(context.Background())
}

// WithOption set option for mongodb
func WithOption() *options.ClientOptions {
	return options.Client()
}

type customLogger struct {
	zapLogger *zap.Logger
}

func (l *customLogger) Info(_ int, msg string, kvs ...interface{}) {
	l.zapLogger.Info(msg, zap.String("msg", msg), zap.Any("kv", kvs))
}

func (l *customLogger) Error(err error, msg string, kvs ...interface{}) {
	l.zapLogger.Warn(msg, zap.Error(err), zap.String("msg", msg), zap.Any("kv", kvs))
}

// NewCustomLogger create a custom logger for mongodb, debug level is used by default.
// example: WithOption().SetLoggerOptions(NewCustomLogger(logger.Get(), true))
func NewCustomLogger(l *zap.Logger, isDebugLevel bool) *options.LoggerOptions {
	if l == nil {
		l, _ = zap.NewProduction()
	}
	sink := &customLogger{zapLogger: l}

	level := options.LogLevelInfo
	if isDebugLevel {
		level = options.LogLevelDebug
	}

	// Create a client with our logger options.
	return options.
		Logger().
		SetSink(sink).
		SetMaxDocumentLength(300).
		SetComponentLevel(options.LogComponentCommand, level)
}
