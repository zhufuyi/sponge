package generate

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// ModelCommand generate model code
func ModelCommand() *cobra.Command {
	var (
		outPath string // 输出目录
		sqlArgs = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "model",
		Short: "Generate model code",
		Long: `generate model code.

Examples:
  # generate model code and embed 'gorm.Model' struct.
  sponge model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate model code, structure fields correspond to the column names of the table.
  sponge model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate model code and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			return runGenModelCommand(codes, outPath)
		},
	}

	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./model_<time>")

	return cmd
}

func runGenModelCommand(codes map[string]string, outPath string) error {
	subTplName := "model"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{"internal/model"}              // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{}                           // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"init.go", "init_test.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addModelFields(r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, out = %s\n\n", subTplName, r.GetOutputDir())
	return nil
}

func addModelFields(r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
