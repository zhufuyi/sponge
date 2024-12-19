package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/gofile"
	"github.com/go-dev-frame/sponge/pkg/replacer"
)

var (
	printCustomContent *strings.Builder
)

// FieldCommand generate code based on custom template and fields
func FieldCommand() *cobra.Command {
	var (
		tplDir     = ""   // template directory
		fieldsFile = ""   // json file
		onlyPrint  bool   // only print template code and all fields, do not generate code
		outPath    string // output directory
	)
	printCustomContent = new(strings.Builder)

	cmd := &cobra.Command{
		Use:   "field",
		Short: "Generate code based on custom template and fields",
		Long:  "Generate code based on custom template and fields.",
		Example: color.HiBlackString(`  # Generate code.
  sponge template field --tpl-dir=yourTemplateDir --fields=yourDefineFields.json

  # Print template code and all fields, do not generate code.
  sponge template field --tpl-dir=yourTemplateDir --fields=yourDefineFields.json --only-print

  # Generate code and specify output directory. Note: code generation will be canceled when the latest generated file already exists.
  sponge template field --tpl-dir=yourTemplateDir --fields=yourDefineFields.json --out=./yourDir`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if files, err := gofile.ListFiles(tplDir); err != nil {
				return err
			} else if len(files) == 0 {
				return fmt.Errorf("no template files found in directory '%s'", tplDir)
			}

			m, err := parseFields(fieldsFile)
			if err != nil {
				return err
			}
			if len(m) == 0 {
				return fmt.Errorf("no fields found in json file %s", fieldsFile)
			}

			g := customGenerator{
				tplDir:    tplDir,
				fields:    m,
				onlyPrint: onlyPrint,
				outPath:   outPath,
			}
			outPath, err = g.generateCode()
			if err != nil {
				return err
			}

			if onlyPrint {
				fmt.Printf("%s", printCustomContent.String())
			} else {
				fmt.Printf("generate custom code successfully, out = %s\n", outPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&tplDir, "tpl-dir", "i", "", "directory where your template code is located")
	_ = cmd.MarkFlagRequired("tpl-dir")
	cmd.Flags().StringVarP(&fieldsFile, "fields", "f", "", "fields defined in JSON file")
	_ = cmd.MarkFlagRequired("fields")
	cmd.Flags().BoolVarP(&onlyPrint, "only-print", "n", false, "only print template code and all fields, do not generate code")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./custom_<time>")

	return cmd
}

type customGenerator struct {
	tplDir    string
	fields    map[string]interface{}
	onlyPrint bool
	outPath   string
}

func (g *customGenerator) generateCode() (string, error) {
	subTplName := "custom"
	r, _ := replacer.New(g.tplDir)
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	files := r.GetFiles()
	if len(files) == 0 {
		return "", errors.New("no template files found")
	}

	if g.onlyPrint {
		listTemplateFiles(printCustomContent, files)
		printCustomContent.WriteString("\n\nAll fields name and value:\n")
		listFields(printCustomContent, g.fields)
		return "", nil
	}

	_ = r.SetOutputDir(g.outPath, subTplName)
	if err := r.SaveTemplateFiles(g.fields, gofile.GetSuffixDir(g.tplDir)); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}
