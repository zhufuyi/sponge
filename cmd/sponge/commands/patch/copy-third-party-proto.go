package patch

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

// CopyThirdPartyProtoCommand copy third-party proto files
func CopyThirdPartyProtoCommand() *cobra.Command {
	var (
		outPath    string // output directory
		isLogExist bool
	)

	cmd := &cobra.Command{
		Use:   "copy-third-party-proto",
		Short: "Copy third-party proto files",
		Long:  "Copy third-party proto files to local directory.",
		Example: color.HiBlackString(`  # Copy third-party proto files to current directory
  sponge patch copy-third-party-proto

  # Copy third-party proto files to yourServerDir
  sponge patch copy-third-party-proto --out=./yourServerDir`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := outPath + gofile.GetPathDelimiter() + "third_party"
			if gofile.IsExists(out) {
				if isLogExist {
					fmt.Printf("%s proto files already exists, skip copying.\n", out)
				}
				return nil
			}

			var err error
			out, err = runCopyThirdPartyProtoCommand(outPath)
			if err != nil {
				return err
			}
			fmt.Printf("copied third_party proto files to %s\n", out)

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "out", "o", ".", "output directory")
	cmd.Flags().BoolVarP(&isLogExist, "is-log-exist", "l", false, "is log file exist")

	return cmd
}

func runCopyThirdPartyProtoCommand(out string) (string, error) {
	r := generate.Replacers[generate.TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{"sponge/third_party"}

	r.SetSubDirsAndFiles(subDirs)
	_ = r.SetOutputDir(out)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}
