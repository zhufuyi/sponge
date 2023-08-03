package merge

import (
	"github.com/spf13/cobra"
)

// GRPCServiceCode merge the grpc service code
func GRPCServiceCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc-pb",
		Short: "Merge the generated grpc related code into the template file",
		Long: `merge the generated grpc related code into the template file.

Examples:
  sponge merge rpc-pb
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mergeGRPCECode()
			mergeGRPCServiceTmpl()
			return nil
		},
	}

	return cmd
}
