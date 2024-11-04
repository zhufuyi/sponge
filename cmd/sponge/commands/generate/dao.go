package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"
)

// DaoCommand generate dao code
func DaoCommand(parentName string) *cobra.Command {
	var (
		moduleName      string // go.mod module name
		outPath         string // output directory
		dbTables        string // table names
		isIncludeInitDB bool

		sqlArgs = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}

		serverName     string // server name
		suitedMonoRepo bool   // whether the generated code is suitable for mono-repo
	)

	cmd := &cobra.Command{
		Use:   "dao",
		Short: "Generate dao code based on sql",
		Long:  "Generate dao code based on sql.",
		Example: color.HiBlackString(fmt.Sprintf(`  # Generate dao code.
  sponge %s dao --module-name=yourModuleName --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # Generate dao code with multiple table names.
  sponge %s dao --module-name=yourModuleName --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2

  # Generate dao code with extened api.
  sponge %s dao --module-name=yourModuleName --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --extended-api=true

  # Generate dao code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge %s dao --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir

  # If you want the generated code to suited to mono-repo, you need to set the parameter --suited-mono-repo=true --server-name=yourServerName`,
			parentName, parentName, parentName, parentName)),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mdName, srvName, smr := getNamesFromOutDir(outPath)
			if mdName != "" {
				moduleName = mdName
				serverName = srvName
				suitedMonoRepo = smr
			} else if moduleName == "" {
				return fmt.Errorf(`required flag(s) "module-name" not set, use "sponge %s dao -h" for help`, parentName)
			}
			if suitedMonoRepo {
				if serverName == "" {
					return fmt.Errorf(`required flag(s) "server-name" not set, use "sponge %s dao -h" for help`, parentName)
				}
				serverName = convertServerName(serverName)
				outPath = changeOutPath(outPath, serverName)
			}

			tableNames := strings.Split(dbTables, ",")
			for count, tableName := range tableNames {
				if tableName == "" {
					continue
				}

				if sqlArgs.DBDriver == DBDriverMongodb {
					sqlArgs.IsEmbed = false
				}
				sqlArgs.DBTable = tableName
				codes, err := sql2code.Generate(&sqlArgs)
				if err != nil {
					return err
				}

				// control to generate the initialization db code only once
				if count == 0 && isIncludeInitDB {
					isIncludeInitDB = true
				} else {
					isIncludeInitDB = false
				}

				g := &daoGenerator{
					moduleName:      moduleName,
					dbDriver:        sqlArgs.DBDriver,
					isIncludeInitDB: isIncludeInitDB,
					codes:           codes,
					outPath:         outPath,
					serverName:      serverName,
					isEmbed:         sqlArgs.IsEmbed,
					isExtendedAPI:   sqlArgs.IsExtendedAPI,
					suitedMonoRepo:  suitedMonoRepo,
				}
				outPath, err = g.generateCode()
				if err != nil {
					return err
				}
			}

			fmt.Printf(`
using help:
  move the folder "internal" to your project code folder.

`)
			fmt.Printf("generate \"dao\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	//_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDriver, "db-driver", "k", "mysql", "database driver, support mysql, mongodb, postgresql, tidb, sqlite")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "database content address, e.g. user:password@(host:port)/database. Note: if db-driver=sqlite, db-dsn must be a local sqlite db file, e.g. --db-dsn=/tmp/sponge_sqlite.db") //nolint
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", false, "whether to embed gorm.model struct")
	cmd.Flags().BoolVarP(&sqlArgs.IsExtendedAPI, "extended-api", "a", false, "whether to generate extended crud api, additional includes: DeleteByIDs, GetByCondition, ListByIDs, ListByLatestID")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	cmd.Flags().BoolVarP(&suitedMonoRepo, "suited-mono-repo", "l", false, "whether the generated code is suitable for mono-repo")
	cmd.Flags().IntVarP(&sqlArgs.JSONNamedType, "json-name-type", "j", 1, "json tags name type, 0:snake case, 1:camel case")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./dao_<time>, "+flagTip("module-name"))
	cmd.Flags().BoolVarP(&isIncludeInitDB, "include-init-db", "i", false, "if true, includes mysql and redis initialization code")

	return cmd
}

type daoGenerator struct {
	moduleName      string
	dbDriver        string
	isIncludeInitDB bool
	codes           map[string]string
	outPath         string
	isEmbed         bool
	isExtendedAPI   bool
	serverName      string
	suitedMonoRepo  bool

	fields []replacer.Field
}

func (g *daoGenerator) generateCode() (string, error) {
	subTplName := codeNameDao
	r, _ := replacer.New(SpongeDir)
	if r == nil {
		return "", errors.New("r is nil")
	}

	// specify the subdirectory and files
	subDirs := []string{}
	subFiles := []string{}

	selectFiles := map[string][]string{
		"internal/cache": {
			"userExample.go", "userExample_test.go",
		},
		"internal/dao": {
			"userExample.go", "userExample_test.go",
		},
		"internal/model": {
			"userExample.go",
		},
	}

	info := g.codes[parser.CodeTypeCrudInfo]
	crudInfo, _ := unmarshalCrudInfo(info)
	if crudInfo.CheckCommonType() {
		selectFiles = map[string][]string{
			"internal/cache": {
				"userExample.go.tpl",
			},
			"internal/dao": {
				"userExample.go.tpl",
			},
			"internal/model": {
				"userExample.go",
			},
		}
		var fields []replacer.Field
		if g.isExtendedAPI {
			selectFiles["internal/dao"] = []string{"userExample.go.exp.tpl"}
			fields = commonDaoExtendedFields(r)
		} else {
			fields = commonDaoFields(r)
		}
		contentFields, err := replaceFilesContent(r, getTemplateFiles(selectFiles), crudInfo)
		if err != nil {
			return "", err
		}
		g.fields = append(g.fields, contentFields...)
		g.fields = append(g.fields, fields...)
	}

	replaceFiles := make(map[string][]string)
	switch strings.ToLower(g.dbDriver) {
	case DBDriverMysql, DBDriverPostgresql, DBDriverTidb, DBDriverSqlite:
		g.fields = append(g.fields, getExpectedSQLForDeletionField(g.isEmbed)...)
		if g.isExtendedAPI {
			var fields []replacer.Field
			if !crudInfo.CheckCommonType() {
				replaceFiles, fields = daoExtendedAPI(r)
			}
			g.fields = append(g.fields, fields...)
		}

	case DBDriverMongodb:
		if g.isExtendedAPI {
			var fields []replacer.Field
			replaceFiles, fields = daoMongoDBExtendedAPI(r)
			g.fields = append(g.fields, fields...)
		} else {
			replaceFiles = map[string][]string{
				"internal/cache": {
					"userExample.go.mgo",
				},
				"internal/dao": {
					"userExample.go.mgo",
				},
			}
		}

	default:
		return "", dbDriverErr(g.dbDriver)
	}

	subFiles = append(subFiles, getSubFiles(selectFiles, replaceFiles)...)

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	_ = r.SetOutputDir(g.outPath, subTplName)
	fields := g.addFields(r)
	r.SetReplacementFields(fields)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

// set fields
func (g *daoGenerator) addFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field
	fields = append(fields, g.fields...)
	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoMgoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the model/userExample.go file
			Old: modelFileMark,
			New: g.codes[parser.CodeTypeModel],
		},
		{
			Old: daoFileMark,
			New: g.codes[parser.CodeTypeDAO],
		},
		{
			Old: selfPackageName + "/" + r.GetSourcePath(),
			New: g.moduleName,
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: g.moduleName,
		},
		{
			Old: g.moduleName + pkgPathSuffix,
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old: "init.go.mgo",
			New: "init.go",
		},
		{
			Old: "userExample.go.mgo",
			New: "userExample.go",
		},
		{
			Old:             "UserExample",
			New:             g.codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	if g.suitedMonoRepo {
		fs := SubServerCodeFields(g.moduleName, g.serverName)
		fields = append(fields, fs...)
	}

	return fields
}

func daoExtendedAPI(r replacer.Replacer) (map[string][]string, []replacer.Field) {
	replaceFiles := map[string][]string{
		"internal/dao": {
			"userExample.go.exp", "userExample_test.go.exp",
		},
	}
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoFile+expSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile+expSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample.go.exp",
			New: "userExample.go",
		},
		{
			Old: "userExample_test.go.exp",
			New: "userExample_test.go",
		},
	}...)

	return replaceFiles, fields
}

func daoMongoDBExtendedAPI(r replacer.Replacer) (map[string][]string, []replacer.Field) {
	replaceFiles := map[string][]string{
		"internal/cache": {
			"userExample.go.mgo",
		},
		"internal/dao": {
			"userExample.go.mgo.exp",
		},
	}

	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoMgoFile+expSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample.go.mgo.exp",
			New: "userExample.go",
		},
		{
			Old: "userExample.go.mgo",
			New: "userExample.go",
		},
	}...)

	return replaceFiles, fields
}

func commonDaoFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoFile+tplSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile+tplSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample.go.tpl",
			New: "userExample.go",
		},
	}...)

	return fields
}

func commonDaoExtendedFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoFile+expSuffix+tplSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile+expSuffix+tplSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample.go.tpl",
			New: "userExample.go",
		},
		{
			Old: "userExample.go.exp.tpl",
			New: "userExample.go",
		},
	}...)

	return fields
}
