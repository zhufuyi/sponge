package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// ModelCommand generate model codes
func ModelCommand(parentName string) *cobra.Command {
	var (
		outPath  string // output directory
		dbTables string // table names

		sqlArgs = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "model",
		Short: "Generate model codes based on mysql table",
		Long: fmt.Sprintf(`generate model codes based on mysql table.

Examples:
  # generate model codes and embed 'gorm.Model' struct.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate model codes with multiple table names.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2

  # generate model codes, structure fields correspond to the column names of the table.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate model codes and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge %s model --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`, parentName, parentName, parentName, parentName),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
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

				dir, err := runGenModelCommand(codes, outPath)
				if err != nil {
					return err
				}
				if outPath == "" {
					outPath = dir
				}
			}

			fmt.Printf("generate 'model' codes successfully, out = %s\n\n", outPath)
			fmt.Printf(`Instructions for use:
	Generate structures corresponding to gorm based on mysql tables.

`)
			return nil
		},
	}

	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./model_<time>")

	return cmd
}

func runGenModelCommand(codes map[string]string, outPath string) (string, error) {
	subTplName := "model"
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{"internal/model"} // only the specified subdirectory is processed, if empty or no subdirectory is specified, it means all files
	ignoreDirs := []string{}              // specify the directory in the subdirectory where processing is ignored
	ignoreFiles := []string{              // specify the files in the subdirectory to be ignored for processing
		"init.go", "init_test.go",
	}

	r.SetSubDirsAndFiles(subDirs)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addModelFields(r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func addModelFields(r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the model/userExample.go file
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
