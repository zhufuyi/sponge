package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/replacer"
	"github.com/go-dev-frame/sponge/pkg/sql2code"
	"github.com/go-dev-frame/sponge/pkg/sql2code/parser"
)

// ModelCommand generate model code
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
		Short: "Generate model code based on sql",
		Long:  "Generate model code based on sql.",
		Example: color.HiBlackString(fmt.Sprintf(`  # Generate model code.
  sponge %s model --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # Generate model code with multiple table names.
  sponge %s model --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2

  # Generate model code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge %s model --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir`,
			parentName, parentName, parentName)),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			tableNames := strings.Split(dbTables, ",")
			for _, tableName := range tableNames {
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

				g := &modelGenerator{
					codes:   codes,
					outPath: outPath,
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
			fmt.Printf("generate \"model\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&sqlArgs.DBDriver, "db-driver", "k", "mysql", "database driver, support mysql, mongodb, postgresql, sqlite")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "database content address, e.g. user:password@(host:port)/database. Note: if db-driver=sqlite, db-dsn must be a local sqlite db file, e.g. --db-dsn=/tmp/sponge_sqlite.db") //nolint
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", false, "whether to embed gorm.model struct")
	cmd.Flags().IntVarP(&sqlArgs.JSONNamedType, "json-name-type", "j", 1, "json tags name type, 0:snake case, 1:camel case")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./model_<time>")

	return cmd
}

type modelGenerator struct {
	codes   map[string]string
	outPath string
}

func (g *modelGenerator) generateCode() (string, error) {
	subTplName := codeNameModel
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// specify the subdirectory and files
	subDirs := []string{}
	subFiles := []string{"internal/model/userExample.go"}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	fields := g.addFields(r)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(g.outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func (g *modelGenerator) addFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the model/userExample.go file
			Old: modelFileMark,
			New: g.codes[parser.CodeTypeModel],
		},
		{
			Old:             "UserExample",
			New:             g.codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
