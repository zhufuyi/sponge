package commands

import (
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/patch"
)

// PatchCommand patch server code
func PatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "patch",
		Short:         "Command set for patching service code",
		Long:          `command set for patching service code.`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		patch.DeleteJSONOmitemptyCommand(),
		patch.GenerateDBInitCommand(),
		patch.GenMysqlInitCommand(),
		patch.GenTypesPbCommand(),
		patch.CopyProtoCommand(),
		patch.ModifyDuplicateNumCommand(),
		patch.ModifyDuplicateErrCodeCommand(),
		patch.AdaptMonoRepoCommand(),
		patch.AddSpecialTypesCommand(),
		patch.ModifyProtoPackageCommand(),
	)

	return cmd
}
