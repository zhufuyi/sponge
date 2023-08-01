package merge

import (
	"github.com/spf13/cobra"
)

// GRPCServiceCode merge the grpc service code
func GRPCServiceCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grpc-service",
		Short: "Merge the grpc service code",
		Long: `merge the grpc service code.

Examples:
  # merge grpc service code
  sponge merge grpc-service
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	return cmd
}
