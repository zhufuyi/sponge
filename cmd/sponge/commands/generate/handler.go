package generate

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// HandlerCommand generate handler codes
func HandlerCommand() *cobra.Command {
	var (
		moduleName string // go.mod文件的module名称
		outPath    string // 输出目录
		sqlArgs    = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Generate handler codes based on mysql",
		Long: `generate handler codes based on mysql.

Examples:
  # generate handler codes and embed 'gorm.model' struct.
  sponge web handler --module-name=yourModuleName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate handler codes, structure fields correspond to the column names of the table.
  sponge web handler --module-name=yourModuleName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate handler codes and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge web handler --module-name=yourModuleName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			return runGenHandlerCommand(moduleName, codes, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the 'go.mod' file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")

	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./handler_<time>")

	return cmd
}

func runGenHandlerCommand(moduleName string, codes map[string]string, outPath string) error {
	subTplName := "handler"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{"internal/model", "internal/cache", "internal/dao",
		"internal/ecode", "internal/handler", "internal/routers", "internal/types"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"init.go", "init_test.go", "swagger_types.go", "http_systemCode.go",
		"grpc_systemCode.go", "grpc_userExample.go", "grpc_systemCode_test.go", "routers.go",
		"routers_test.go", "routers_gwExample.go", "routers_gwExample_test.go", "userExample_gwExample.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addHandlerFields(moduleName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' codes successfully, out = %s\n\n", subTplName, r.GetOutputDir())
	return nil
}

func addHandlerFields(moduleName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, handlerFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, handlerTestFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
		},
		{ // 替换handler/userExample.go文件内容
			Old: handlerFileMark,
			New: adjustmentOfIDType(codes[parser.CodeTypeHandler]),
		},
		{
			Old: selfPackageName + "/" + r.GetSourcePath(),
			New: moduleName,
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: moduleName,
		},
		{
			Old: "userExampleNO = 1",
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(1000)),
		},
		{
			Old: moduleName + "/pkg",
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
