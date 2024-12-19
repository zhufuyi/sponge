package commands

import (
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands/merge"
)

// MergeCommand merge the generated code
func MergeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge the generated code into the template file",
		Long: `Merge the generated code into the template file, you don't worry about it affecting
the logic code you have already written, in case of accidents, you can find the
pre-merge code in the directory /tmp/sponge_merge_backup_code`,
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
