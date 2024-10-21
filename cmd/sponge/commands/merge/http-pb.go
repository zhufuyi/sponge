package merge

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// GinHandlerCode merge the gin handler code
func GinHandlerCode() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "http-pb",
		Short: "Merge the generated http related code into the template file",
		Long:  "Merge the generated http related code into the template file.",
		Example: color.HiBlackString(`  # Merge go template file in local server directory
  sponge merge http-pb

  # Merge go template file in specified directory
  sponge merge http-pb --dir=/path/to/server/directory`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir = adaptDir(dir)
			mergeHTTPECode(dir)
			mergeGinRouters(dir)
			mergeHTTPHandlerTmpl(dir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "input directory")

	return cmd
}
