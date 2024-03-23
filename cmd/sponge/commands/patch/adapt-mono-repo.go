package patch

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

// AdaptMonoRepoCommand Adapt to mono-repo command
func AdaptMonoRepoCommand() *cobra.Command {
	var (
		dir        string
		moduleName string // module name for go.mod
		serverName string // server name
	)

	cmd := &cobra.Command{
		Use:   "adapt-mono-repo",
		Short: "Adapt to mono-repo in api directory code",
		Long: `adapt to mono-repo in api directory code

Examples:
  # adapt to mono-repo code
  sponge patch adapt-mono-repo

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mdName, srvName, _ := getNamesFromOutDir(dir)
			if mdName != "" {
				moduleName = mdName
			} else if moduleName == "" {
				return errors.New(`can't get info from docs/gen.info`)
			}
			if srvName != "" {
				serverName = srvName
			} else if serverName == "" {
				return errors.New(`can't get info from docs/gen.info`)
			}

			files, err := gofile.ListFiles(dir, gofile.WithSuffix(".go"))
			if err != nil {
				return err
			}

			var oldStr = fmt.Sprintf("\"%s/api", moduleName)
			var newStr = fmt.Sprintf("\"%s/api", moduleName+"/"+serverName)
			for _, file := range files {
				data, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				if bytes.Contains(data, []byte(oldStr)) {
					data = bytes.ReplaceAll(data, []byte(oldStr), []byte(newStr))
					err = os.WriteFile(file, data, 0766)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "input directory")

	return cmd
}
