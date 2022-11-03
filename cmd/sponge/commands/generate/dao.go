package generate

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// DaoCommand generate dao codes
func DaoCommand() *cobra.Command {
	var (
		moduleName string // 服务名称，也就是包名称
		outPath    string // 输出目录
		sqlArgs    = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "dao",
		Short: "Generate dao codes",
		Long: `generate dao codes.

Examples:
  # generate dao codes and embed 'gorm.model' struct.
  sponge dao --module-name=yourModuleName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate dao codes, structure fields correspond to the column names of the table.
  sponge dao --module-name=yourModuleName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate dao codes and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge dao --module-name=yourModuleName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			return runGenDaoCommand(moduleName, codes, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the 'go.mod' file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./dao_<time>")

	return cmd
}

func runGenDaoCommand(moduleName string, codes map[string]string, outPath string) error {
	subTplName := "dao"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("r is nil")
	}

	// 设置模板信息
	subDirs := []string{"internal/model", "internal/cache", "internal/dao"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{}                                                // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"init.go", "init_test.go"}                      // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addDAOFields(moduleName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' codes successfully, out = %s\n\n", subTplName, r.GetOutputDir())
	return nil
}

func addDAOFields(moduleName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
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
