package merge

import (
	"github.com/spf13/cobra"
)

// GinHandlerCode merge the gin handler code
func GinHandlerCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http-pb",
		Short: "Merge the generated http related code into the template file",
		Long: `merge the generated http related code into the template file.

Examples:
  sponge merge http-pb
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mergeHTTPECode()
			mergeGinRouters()
			mergeHTTPHandlerTmpl()
			return nil
		},
	}

	return cmd
}
