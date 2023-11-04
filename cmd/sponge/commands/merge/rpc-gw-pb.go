package merge

import (
	"github.com/spf13/cobra"
)

// GinServiceCode merge the gin service code
func GinServiceCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc-gw-pb",
		Short: "Merge the generated grpc gateway related code into the template file",
		Long: `merge the generated grpc gateway related code into the template file.

Examples:
  sponge merge rpc-gw-pb
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mergeGRPCECode()
			mergeGinRouters()
			mergeGRPCServiceClientTmpl()
			return nil
		},
	}

	return cmd
}
