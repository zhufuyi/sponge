// Package sql2code provides for generating code for different purposes according to sql,
// support generating json, gorm model, update parameter, request parameter code,
// sql can be obtained from parameter, file, db three ways, priority from high to low.
package sql2code

import (
	"errors"
	"fmt"
	"os"

	"github.com/zhufuyi/sponge/pkg/sql2code/parser"
)

// Args generate code arguments
type Args struct {
	SQL string // DDL sql

	DDLFile string // DDL file

	DBDsn   string // connecting to mysql's dsn
	DBTable string

	Package        string // specify the package name (only valid for model types)
	GormType       bool   // whether to display the gorm type name (only valid for model type codes)
	JSONTag        bool   // does it include a json tag
	JSONNamedType  int    // json field naming type, 0: snake case such as my_field_name, 1: camel sase, such as myFieldName
	IsEmbed        bool   // is gorm.Model embedded
	CodeType       string // specify the different types of code to be generated, namely model (default), json, dao, handler, proto
	ForceTableName bool
	Charset        string
	Collation      string
	TablePrefix    string
	ColumnPrefix   string
	NoNullType     bool
	NullStyle      string
}

func (a *Args) checkValid() error {
	if a.SQL == "" && a.DDLFile == "" && (a.DBDsn == "" && a.DBTable == "") {
		return errors.New("you must specify sql or ddl file")
	}
	return nil
}

func getSQL(args *Args) (string, error) {
	if args.SQL != "" {
		return args.SQL, nil
	}

	sql := ""
	if args.DDLFile != "" {
		b, err := os.ReadFile(args.DDLFile)
		if err != nil {
			return sql, fmt.Errorf("read %s failed, %s", args.DDLFile, err)
		}
		return string(b), nil
	} else if args.DBDsn != "" {
		if args.DBTable == "" {
			return sql, errors.New("miss mysql table")
		}
		sqlStr, err := parser.GetTableInfo(args.DBDsn, args.DBTable)
		if err != nil {
			return sql, err
		}
		return sqlStr, nil
	}

	return sql, errors.New("no SQL input(-sql|-f|-db-dsn)")
}

func getOptions(args *Args) []parser.Option {
	var opts []parser.Option

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
	if err := args.checkValid(); err != nil {
		return nil, err
	}

	sql, err := getSQL(args)
	if err != nil {
		return nil, err
	}

	opt := getOptions(args)

	return parser.ParseSQL(sql, opt...)
}
