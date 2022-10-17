package sql2code

import (
	"errors"
	"fmt"
	"os"

	"github.com/zhufuyi/sponge/pkg/sql2code/parser"
)

// Args 参数
type Args struct {
	SQL string // DDL sql

	DDLFile string // 读取文件的DDL sql

	DBDsn   string // 从db获取表的DDL sql
	DBTable string

	Package        string // 生成字段的包名(只有model类型有效)
	GormType       bool   // 是否显示gorm type名称(只有model类型代码有效)
	JSONTag        bool   // 是否包括json tag
	JSONNamedType  int    // json命名类型，0:和列名一致，其他值表示驼峰
	IsEmbed        bool   // 是否嵌入gorm.Model
	CodeType       string // 指定生成代码用途，支持4中类型，分别是 model(默认), json, dao, handler
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

// GenerateOne 根据sql生成gorm代码，sql可以从参数、文件、db三种方式获取，优先从高到低
func GenerateOne(args *Args) (string, error) {
	codes, err := Generate(args)
	if err != nil {
		return "", err
	}

	if args.CodeType == "" {
		args.CodeType = parser.CodeTypeModel // 默认为model code
	}
	out, ok := codes[args.CodeType]
	if !ok {
		return "", fmt.Errorf("unknown code type %s", args.CodeType)
	}

	return out, nil
}

// Generate 生成model, json, dao, handler不同用途代码
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
