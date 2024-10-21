package patch

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

// CopyGOModCommand copy go mod files
func CopyGOModCommand() *cobra.Command {
	var (
		moduleName     string // module name for go.mod
		outPath        string // output directory
		isLogExist     bool
		isForceReplace bool
	)

	cmd := &cobra.Command{
		Use:   "copy-go-mod",
		Short: "Copy go mod files",
		Long:  "Copy go mod files to local directory.",
		Example: color.HiBlackString(`  # Copy go mod files to current directory
  sponge patch copy-go-mod --module-name=yourModuleName

  # Copy go mod files to yourServerDir, module name from out directory
  sponge patch copy-go-mod --out=./yourServerDir`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if moduleName == "" {
				mn, _, _ := getNamesFromOutDir(outPath)
				if mn == "" {
					return errors.New("module-name is required, please use --module-name to set it")
				}
				moduleName = mn
			}

			goModFile := outPath + gofile.GetPathDelimiter() + "go.mod"
			if gofile.IsExists(goModFile) {
				if !isForceReplace {
					if isLogExist {
						fmt.Printf("%s already exists, skip copying.\n", goModFile)
					}
					return nil
				}
				// delete the go.mod and go.sum file if it exists
				_ = os.RemoveAll(goModFile)
				_ = os.RemoveAll(strings.TrimSuffix(goModFile, ".mod") + ".sum")
			}

			out, err := runCopyGoModCommand(moduleName, outPath)
			if err != nil {
				return err
			}
			fmt.Printf("copied go.mod to %s\n", out)

			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	cmd.Flags().StringVarP(&outPath, "out", "o", ".", "output directory")
	cmd.Flags().BoolVarP(&isLogExist, "is-log-exist", "l", false, "whether to log file exist")
	cmd.Flags().BoolVarP(&isForceReplace, "is-force-replace", "f", false, "whether to force  replace the go.mod file")

	return cmd
}

func runCopyGoModCommand(moduleName string, out string) (string, error) {
	r := generate.Replacers[generate.TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subFiles := []string{"sponge/go.mod", "sponge/go.sum"}
	r.SetSubDirsAndFiles(nil, subFiles...)
	r.SetReplacementFields(generate.GetGoModFields(moduleName))
	_ = r.SetOutputDir(out)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}
