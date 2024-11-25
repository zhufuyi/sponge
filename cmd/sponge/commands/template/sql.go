package template

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"
)

var (
	printSQLOnce    sync.Once
	printSQLContent *strings.Builder
)

// SQLCommand generate code based on sql and custom template
func SQLCommand() *cobra.Command {
	var (
		tplDir     = ""   // template directory
		fieldsFile = ""   // fields defined in json
		dbTables   string // table names
		outPath    string // output directory
		onlyPrint  bool   // only print template code and all fields

		sqlArgs = sql2code.Args{
			JSONTag:          true,
			GormType:         true,
			IsCustomTemplate: true,
		}
	)
	printSQLOnce = sync.Once{}
	printSQLContent = new(strings.Builder)

	cmd := &cobra.Command{
		Use:   "sql",
		Short: "Generate code based on sql and custom template",
		Long:  "Generate code based on sql and custom template.",
		Example: color.HiBlackString(`  # Generate code.
  sponge template sql --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --tpl-dir=yourTemplateDir

  # Generate code and specify fields defined in json file.
  sponge template sql --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --tpl-dir=yourTemplateDir --fields=yourDefineFields.json

  # Print template code and all fields, do not generate code.
  sponge template sql --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --tpl-dir=yourTemplateDir --fields=yourDefineFields.json --only-print

  # Generate code with multiple table names, each table generates code independently based on the template files.
  sponge template sql --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=t1,t2 --tpl-dir=yourTemplateDir

  # Generate code and specify output directory. Note: code generation will be canceled when the latest generated file already exists.
  sponge template sql --db-driver=mysql --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --tpl-dir=yourTemplateDir --out=./yourDir`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if files, err := gofile.ListFiles(tplDir); err != nil {
				return err
			} else if len(files) == 0 {
				return fmt.Errorf("no template files found in directory '%s'", tplDir)
			}

			var m map[string]interface{}
			if fieldsFile != "" {
				var err error
				m, err = parseFields(fieldsFile)
				if err != nil {
					return err
				}
			}

			tableNames := strings.Split(dbTables, ",")
			l := len(tableNames)
			for i, tableName := range tableNames {
				if tableName == "" {
					continue
				}

				sqlArgs.DBTable = tableName
				codes, err := sql2code.Generate(&sqlArgs)
				if err != nil {
					return err
				}

				tableInfo, err := parser.UnMarshalTableInfo(codes[parser.CodeTypeTableInfo])
				if err != nil {
					return err
				}
				fields, err := mergeFields(tableInfo, m)
				if err != nil {
					return err
				}

				g := sqlGenerator{
					tplDir:    tplDir,
					fields:    fields,
					onlyPrint: onlyPrint,
					outPath:   outPath,
				}
				outPath, err = g.generateCode()
				if err != nil {
					return err
				}

				if i != l-1 {
					printSQLContent.WriteString("\n    " +
						"------------------------------------------------------------------\n\n\n")
				}
			}

			if onlyPrint {
				fmt.Printf("%s", printSQLContent.String())
			} else {
				fmt.Printf("generate custom code successfully, out = %s\n", outPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&sqlArgs.DBDriver, "db-driver", "k", "", "database driver, support mysql, mongodb, postgresql, sqlite")
	_ = cmd.MarkFlagRequired("db-driver")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "database content address, e.g. user:password@(host:port)/database. Note: if db-driver=sqlite, db-dsn must be a local sqlite db file, e.g. --db-dsn=/tmp/sponge_sqlite.db") //nolint
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&dbTables, "db-table", "t", "", "table name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().StringVarP(&sqlArgs.TablePrefix, "table-prefix", "p", "", "table name prefix, e.g. t_")

	cmd.Flags().StringVarP(&tplDir, "tpl-dir", "i", "", "directory where your template code is located")
	_ = cmd.MarkFlagRequired("tpl-dir")
	cmd.Flags().StringVarP(&fieldsFile, "fields", "f", "", "fields defined in json file")
	cmd.Flags().BoolVarP(&onlyPrint, "only-print", "n", false, "only print template code and all fields, do not generate code")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./sql_to_template_<time>")

	return cmd
}

type sqlGenerator struct {
	tplDir    string
	fields    map[string]interface{}
	onlyPrint bool
	outPath   string
}

func (g *sqlGenerator) generateCode() (string, error) {
	subTplName := "sql_to_template"
	r, _ := replacer.New(g.tplDir)
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	files := r.GetFiles()
	if len(files) == 0 {
		return "", errors.New("no template files found")
	}

	if g.onlyPrint {
		printSQLOnce.Do(func() {
			listTemplateFiles(printSQLContent, files)
			printSQLContent.WriteString("\n\nAll fields name and value:\n")
		})
		listFields(printSQLContent, g.fields)
		return "", nil
	}

	_ = r.SetOutputDir(g.outPath, subTplName)
	if err := r.SaveTemplateFiles(g.fields, gofile.GetSuffixDir(g.tplDir)); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}
