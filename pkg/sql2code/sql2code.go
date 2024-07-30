// Package sql2code provides for generating code for different purposes according to sql,
// support generating json, gorm model, update parameter, request parameter code,
// sql can be obtained from parameter, file, db three ways, priority from high to low.
package sql2code

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"
	"github.com/zhufuyi/sponge/pkg/utils"
)

// Args generate code arguments
type Args struct {
	SQL string // DDL sql

	DDLFile string // DDL file

	DBDriver   string            // db driver name, such as mysql, mongodb, postgresql, tidb, sqlite, default is mysql
	DBDsn      string            // connecting to mysql's dsn, if DBDriver is sqlite, DBDsn is local db file
	DBTable    string            // table name
	fieldTypes map[string]string // field name:type

	Package        string // specify the package name (only valid for model types)
	GormType       bool   // whether to display the gorm type name (only valid for model type codes)
	JSONTag        bool   // does it include a json tag
	JSONNamedType  int    // json field naming type, 0: snake case such as my_field_name, 1: camel sase, such as myFieldName
	IsEmbed        bool   // is gorm.Model embedded
	IsWebProto     bool   // proto file type, true: include router path and swagger info, false: normal proto file without router and swagger
	CodeType       string // specify the different types of code to be generated, namely model (default), json, dao, handler, proto
	ForceTableName bool
	Charset        string
	Collation      string
	TablePrefix    string
	ColumnPrefix   string
	NoNullType     bool
	NullStyle      string
	IsExtendedAPI  bool // true: generate extended api (9 api), false: generate basic api (5 api)
}

func (a *Args) checkValid() error {
	if a.SQL == "" && a.DDLFile == "" && (a.DBDsn == "" && a.DBTable == "") {
		return errors.New("you must specify sql or ddl file")
	}
	if a.DBTable != "" {
		tables := strings.Split(a.DBTable, ",")
		for _, name := range tables {
			if strings.HasSuffix(name, "_test") {
				return fmt.Errorf(`the table name (%s) suffix "_test" is not supported for code generation, please delete suffix "_test" or change it to another name. `, name)
			}
		}
	}

	if a.DBDriver == "" {
		a.DBDriver = parser.DBDriverMysql
	} else if a.DBDriver == parser.DBDriverSqlite {
		if !gofile.IsExists(a.DBDsn) {
			return fmt.Errorf("sqlite db file %s not found in local host", a.DBDsn)
		}
	}
	if a.fieldTypes == nil {
		a.fieldTypes = make(map[string]string)
	}
	return nil
}

func getSQL(args *Args) (string, map[string]string, error) {
	if args.SQL != "" {
		return args.SQL, nil, nil
	}

	sql := ""
	dbDriverName := strings.ToLower(args.DBDriver)
	if args.DDLFile != "" {
		if dbDriverName != parser.DBDriverMysql {
			return sql, nil, fmt.Errorf("not support driver %s for parsing the sql file, only mysql is supported", args.DBDriver)
		}
		b, err := os.ReadFile(args.DDLFile)
		if err != nil {
			return sql, nil, fmt.Errorf("read %s failed, %s", args.DDLFile, err)
		}
		return string(b), nil, nil
	} else if args.DBDsn != "" {
		if args.DBTable == "" {
			return sql, nil, errors.New("miss database table")
		}

		switch dbDriverName {
		case parser.DBDriverMysql, parser.DBDriverTidb:
			dsn := utils.AdaptiveMysqlDsn(args.DBDsn)
			sqlStr, err := parser.GetMysqlTableInfo(dsn, args.DBTable)
			return sqlStr, nil, err
		case parser.DBDriverPostgresql:
			dsn := utils.AdaptivePostgresqlDsn(args.DBDsn)
			fields, err := parser.GetPostgresqlTableInfo(dsn, args.DBTable)
			if err != nil {
				return "", nil, err
			}
			sqlStr, pgTypeMap := parser.ConvertToSQLByPgFields(args.DBTable, fields)
			return sqlStr, pgTypeMap, nil
		case parser.DBDriverSqlite:
			sqlStr, err := parser.GetSqliteTableInfo(args.DBDsn, args.DBTable)
			return sqlStr, nil, err
		case parser.DBDriverMongodb:
			dsn := utils.AdaptiveMongodbDsn(args.DBDsn)
			fields, err := parser.GetMongodbTableInfo(dsn, args.DBTable)
			if err != nil {
				return "", nil, err
			}
			sqlStr, mongoTypeMap := parser.ConvertToSQLByMgoFields(args.DBTable, fields)
			return sqlStr, mongoTypeMap, nil
		default:
			return "", nil, fmt.Errorf("getsql error, unsupported database driver: " + dbDriverName)
		}
	}

	return sql, nil, errors.New("no SQL input(-sql|-f|-db-dsn)")
}

func setOptions(args *Args) []parser.Option {
	var opts []parser.Option

	if args.DBDriver != "" {
		opts = append(opts, parser.WithDBDriver(args.DBDriver))
	}
	if args.fieldTypes != nil {
		opts = append(opts, parser.WithFieldTypes(args.fieldTypes))
	}

	if args.Charset != "" {
		opts = append(opts, parser.WithCharset(args.Charset))
	}
	if args.Collation != "" {
		opts = append(opts, parser.WithCollation(args.Collation))
	}
	if args.JSONTag {
		opts = append(opts, parser.WithJSONTag(args.JSONNamedType))
	}
	if args.TablePrefix != "" {
		opts = append(opts, parser.WithTablePrefix(args.TablePrefix))
	}
	if args.ColumnPrefix != "" {
		opts = append(opts, parser.WithColumnPrefix(args.ColumnPrefix))
	}
	if args.NoNullType {
		opts = append(opts, parser.WithNoNullType())
	}
	if args.IsEmbed {
		opts = append(opts, parser.WithEmbed())
	}
	if args.IsWebProto {
		opts = append(opts, parser.WithWebProto())
	}

	if args.NullStyle != "" {
		switch args.NullStyle {
		case "sql":
			opts = append(opts, parser.WithNullStyle(parser.NullInSql))
		case "ptr":
			opts = append(opts, parser.WithNullStyle(parser.NullInPointer))
		default:
			fmt.Printf("invalid null style: %s\n", args.NullStyle)
			return nil
		}
	} else {
		opts = append(opts, parser.WithNullStyle(parser.NullDisable))
	}
	if args.Package != "" {
		opts = append(opts, parser.WithPackage(args.Package))
	}
	if args.GormType {
		opts = append(opts, parser.WithGormType())
	}
	if args.ForceTableName {
		opts = append(opts, parser.WithForceTableName())
	}
	if args.IsExtendedAPI {
		opts = append(opts, parser.WithExtendedAPI())
	}

	return opts
}

// GenerateOne generate gorm code from sql, which can be obtained from parameters, files and db, with priority from highest to lowest
func GenerateOne(args *Args) (string, error) {
	codes, err := Generate(args)
	if err != nil {
		return "", err
	}

	if args.CodeType == "" {
		args.CodeType = parser.CodeTypeModel // default is model code
	}
	out, ok := codes[args.CodeType]
	if !ok {
		return "", fmt.Errorf("unknown code type %s", args.CodeType)
	}

	return out, nil
}

// Generate model, json, dao, handler, proto codes
func Generate(args *Args) (map[string]string, error) {
	args.FormatDsn()
	if err := args.checkValid(); err != nil {
		return nil, err
	}

	sql, fieldTypes, err := getSQL(args)
	if err != nil {
		return nil, err
	}
	if fieldTypes != nil {
		args.fieldTypes = fieldTypes
	}
	if sql == "" {
		return nil, fmt.Errorf("get sql from %s error, maybe the table %s doesn't exist", args.DBDriver, args.DBTable)
	}

	opt := setOptions(args)

	return parser.ParseSQL(sql, opt...)
}

func (a *Args) FormatDsn() {
	dbParams := strings.Split(a.DBDsn, ";")
	a.DBDsn = dbParams[0]
	newParams := dbParams[1:]
	for _, v := range newParams {
		ss := strings.SplitN(v, "=", 2)
		switch ss[0] {
		case "prefix":
			a.TablePrefix = ss[1]
		}
	}
}
