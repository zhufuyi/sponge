package generate

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"
)

// HandlerCommand generate handler code
func HandlerCommand() *cobra.Command {
	var (
		moduleName string // module name for go.mod
		outPath    string // output directory
		dbTables   string // table names

		sqlArgs = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}

		serverName     string // server name
		suitedMonoRepo bool   // whether the generated code is suitable for mono-repo
	)

	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Generate handler CRUD code based on sql",
		Long:  "Generate handler CRUD code based on sql.",
		Example: color.HiBlackString(`  # Generate handler code.
  sponge web handler --module-name=yourModuleName --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # Generate handler code with multiple table names.
  sponge web handler --module-name=yourModuleName --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2

  # Generate handler code with extended api.
  sponge web handler --module-name=yourModuleName --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --extended-api=true

  # Generate handler code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge web handler --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir

  # If you want the generated code to suited to mono-repo, you need to set the parameter --suited-mono-repo=true --server-name=yourServerName`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mdName, srvName, smr := getNamesFromOutDir(outPath)
			if mdName != "" {
				moduleName = mdName
				serverName = srvName
				suitedMonoRepo = smr
			} else if moduleName == "" {
				return errors.New(`required flag(s) "module-name" not set, use "sponge web handler -h" for help`)
			}
			if suitedMonoRepo {
				if serverName == "" {
					return errors.New(`required flag(s) "server-name" not set, use "sponge web handler -h" for help`)
				}
				serverName = convertServerName(serverName)
				outPath = changeOutPath(outPath, serverName)
			}

			if sqlArgs.DBDriver == DBDriverMongodb {
				sqlArgs.IsEmbed = false
			}

			tableNames := strings.Split(dbTables, ",")
			for _, tableName := range tableNames {
				if tableName == "" {
					continue
				}

				sqlArgs.DBTable = tableName
				codes, err := sql2code.Generate(&sqlArgs)
				if err != nil {
					return err
				}

				g := &handlerGenerator{
					moduleName:     moduleName,
					dbDriver:       sqlArgs.DBDriver,
					codes:          codes,
					outPath:        outPath,
					isEmbed:        sqlArgs.IsEmbed,
					isExtendedAPI:  sqlArgs.IsExtendedAPI,
					serverName:     serverName,
					suitedMonoRepo: suitedMonoRepo,
				}
				outPath, err = g.generateCode()
				if err != nil {
					return err
				}
			}

			fmt.Printf(`
using help:
  1. move the folder "internal" to your project code folder.
  2. open a terminal and execute the command: make docs
  3. compile and run service: make run
  4. visit http://localhost:8080/swagger/index.html in your browser, and test the CRUD api interface.

`)
			fmt.Printf("generate \"handler\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	//_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	cmd.Flags().StringVarP(&sqlArgs.DBDriver, "db-driver", "k", "mysql", "database driver, support mysql, mongodb, postgresql, tidb, sqlite")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "database content address, e.g. user:password@(host:port)/database. Note: if db-driver=sqlite, db-dsn must be a local sqlite db file, e.g. --db-dsn=/tmp/sponge_sqlite.db") //nolint
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", false, "whether to embed gorm.model struct")
	cmd.Flags().BoolVarP(&sqlArgs.IsExtendedAPI, "extended-api", "a", false, "whether to generate extended crud api, additional includes: DeleteByIDs, GetByCondition, ListByIDs, ListByLatestID")
	cmd.Flags().BoolVarP(&suitedMonoRepo, "suited-mono-repo", "l", false, "whether the generated code is suitable for mono-repo")
	cmd.Flags().IntVarP(&sqlArgs.JSONNamedType, "json-name-type", "j", 1, "json tags name type, 0:snake case, 1:camel case")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./handler_<time>, "+flagTip("module-name"))

	return cmd
}

type handlerGenerator struct {
	moduleName     string
	dbDriver       string
	codes          map[string]string
	outPath        string
	serverName     string
	isEmbed        bool
	isExtendedAPI  bool
	suitedMonoRepo bool

	fields        []replacer.Field
	isCommonStyle bool
}

func (g *handlerGenerator) generateCode() (string, error) {
	subTplName := codeNameHandler
	r, _ := replacer.New(SpongeDir)
	if r == nil {
		return "", errors.New("replacer is nil")
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
		"internal/ecode": {
			"userExample_http.go",
		},
		"internal/handler": {
			"userExample.go", "userExample_test.go",
		},
		"internal/model": {
			"userExample.go",
		},
		"internal/routers": {
			"userExample.go",
		},
		"internal/types": {
			"userExample_types.go",
		},
	}

	info := g.codes[parser.CodeTypeCrudInfo]
	crudInfo, _ := unmarshalCrudInfo(info)
	if crudInfo.CheckCommonType() {
		g.isCommonStyle = true
		selectFiles = map[string][]string{
			"internal/cache": {
				"userExample.go.tpl",
			},
			"internal/dao": {
				"userExample.go.tpl",
			},
			"internal/ecode": {
				"userExample_http.go.tpl",
			},
			"internal/handler": {
				"userExample.go.tpl",
			},
			"internal/model": {
				"userExample.go",
			},
			"internal/routers": {
				"userExample.go.tpl",
			},
			"internal/types": {
				"userExample_types.go.tpl",
			},
		}
		var fields []replacer.Field
		if g.isExtendedAPI {
			selectFiles["internal/dao"] = []string{"userExample.go.exp.tpl"}
			selectFiles["internal/ecode"] = []string{"userExample_http.go.exp.tpl"}
			selectFiles["internal/handler"] = []string{"userExample.go.exp.tpl"}
			selectFiles["internal/routers"] = []string{"userExample.go.exp.tpl"}
			selectFiles["internal/types"] = []string{"userExample_types.go.exp.tpl"}
			fields = commonHandlerExtendedFields(r)
		} else {
			fields = commonHandlerFields(r)
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
				replaceFiles, fields = handlerExtendedAPI(r, codeNameHandler)
			}
			g.fields = append(g.fields, fields...)
		}

	case DBDriverMongodb:
		if g.isExtendedAPI {
			var fields []replacer.Field
			replaceFiles, fields = handlerMongoDBExtendedAPI(r, codeNameHandler)
			g.fields = append(g.fields, fields...)
		} else {
			replaceFiles = map[string][]string{
				"internal/cache": {
					"userExample.go.mgo",
				},
				"internal/dao": {
					"userExample.go.mgo",
				},
				"internal/handler": {
					"userExample.go.mgo",
				},
				"internal/types": {
					"userExample_types.go.mgo",
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

func (g *handlerGenerator) addFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field
	fields = append(fields, g.fields...)
	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoMgoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, typesFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, typesMgoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, handlerTestFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the model/userExample.go file
			Old: modelFileMark,
			New: g.codes[parser.CodeTypeModel],
		},
		{ // replace the contents of the dao/userExample.go file
			Old: daoFileMark,
			New: g.codes[parser.CodeTypeDAO],
		},
		{ // replace the contents of the handler/userExample.go file
			Old: handlerFileMark,
			New: adjustmentOfIDType(g.codes[parser.CodeTypeHandler], g.dbDriver, g.isCommonStyle),
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
			Old: "userExampleNO       = 1",
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(99)+1),
		},
		{
			Old: g.moduleName + pkgPathSuffix,
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old: "userExample_types.go.mgo",
			New: "userExample_types.go",
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

func handlerExtendedAPI(r replacer.Replacer, codeName string) (map[string][]string, []replacer.Field) {
	replaceFiles := map[string][]string{
		"internal/dao": {
			"userExample.go.exp", "userExample_test.go.exp",
		},
		"internal/ecode": {
			"systemCode_http.go", "userExample_http.go.exp",
		},
		"internal/handler": {
			"userExample.go.exp", "userExample_test.go.exp",
		},
		"internal/routers": {
			"routers.go", "userExample.go.exp",
		},
		"internal/types": {
			"swagger_types.go", "userExample_types.go.exp",
		},
	}
	if codeName == codeNameHandler {
		replaceFiles["internal/ecode"] = []string{"userExample_http.go.exp"}
		replaceFiles["internal/routers"] = []string{"userExample.go.exp"}
		replaceFiles["internal/types"] = []string{"userExample_types.go.exp"}
	}

	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoFile+expSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile+expSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, typesFile+expSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, handlerTestFile+expSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample_types.go.exp",
			New: "userExample_types.go",
		},
		{
			Old: "userExample_http.go.exp",
			New: "userExample_http.go",
		},
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

func handlerMongoDBExtendedAPI(r replacer.Replacer, codeName string) (map[string][]string, []replacer.Field) {
	replaceFiles := map[string][]string{
		"internal/cache": {
			"userExample.go.mgo",
		},
		"internal/dao": {
			"userExample.go.mgo.exp",
		},
		"internal/ecode": {
			"systemCode_http.go", "userExample_http.go.exp",
		},
		"internal/handler": {
			"userExample.go.mgo.exp",
		},
		"internal/model": {
			"init.go.mgo", "userExample.go",
		},
		"internal/routers": {
			"routers.go", "userExample.go.exp",
		},
		"internal/types": {
			"swagger_types.go", "userExample_types.go.mgo.exp",
		},
	}
	if codeName == codeNameHandler {
		replaceFiles["internal/ecode"] = []string{"userExample_http.go.exp"}
		replaceFiles["internal/model"] = []string{"userExample.go"}
		replaceFiles["internal/routers"] = []string{"userExample.go.exp"}
		replaceFiles["internal/types"] = []string{"userExample_types.go.mgo.exp"}
	}

	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoMgoFile+expSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, typesMgoFile+expSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample_http.go.exp",
			New: "userExample_http.go",
		},
		{
			Old: "userExample_types.go.mgo.exp",
			New: "userExample_types.go",
		},
		{
			Old: "userExample.go.mgo.exp",
			New: "userExample.go",
		},
		{
			Old: "userExample.go.exp",
			New: "userExample.go",
		},
	}...)

	return replaceFiles, fields
}

func commonHandlerFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoFile+tplSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile+tplSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, typesFile+tplSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample_http.go.tpl",
			New: "userExample_http.go",
		},
		{
			Old: "userExample_types.go.tpl",
			New: "userExample_types.go",
		},
		{
			Old: "userExample.go.tpl",
			New: "userExample.go",
		},
	}...)

	return fields
}

func commonHandlerExtendedFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, daoFile+expSuffix+tplSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile+expSuffix+tplSuffix, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, typesFile+expSuffix+tplSuffix, startMark, endMark)...)

	fields = append(fields, []replacer.Field{
		{
			Old: "userExample_http.go.exp.tpl",
			New: "userExample_http.go",
		},
		{
			Old: "userExample_types.go.exp.tpl",
			New: "userExample_types.go",
		},
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
