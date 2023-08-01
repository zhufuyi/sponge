package commands

import (
	"github.com/zhufuyi/sponge/cmd/sponge/commands/merge"

	"github.com/spf13/cobra"
)

// MergeCommand merge the generated code
func MergeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "merge",
		Short:         "Merge the generated code, including gin-handler, gin-service and grpc-service",
		Long:          "merge the generated code, including gin-handler, gin-service and grpc-service.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		merge.GinHandlerCode(),
		merge.GinServiceCode(),
		merge.GRPCServiceCode(),
	)

	return cmd
}
