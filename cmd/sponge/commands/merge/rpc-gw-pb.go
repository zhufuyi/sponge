package merge

import (
	"github.com/spf13/cobra"
)

// GinServiceCode merge the gin service code
func GinServiceCode() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "rpc-gw-pb",
		Short: "Merge the generated grpc gateway related code into the template file",
		Long: `merge the generated grpc gateway related code into the template file.

Examples:
  sponge merge rpc-gw-pb
  sponge merge rpc-gw-pb --dir=yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir = adaptDir(dir)
			mergeGRPCECode(dir)
			mergeGinRouters(dir)
			mergeGRPCServiceClientTmpl(dir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "input directory")

	return cmd
}
