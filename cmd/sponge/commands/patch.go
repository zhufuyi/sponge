package commands

import (
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands/patch"
)

// PatchCommand patch server code
func PatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "patch",
		Short:         "Patch the generated code",
		Long:          `Patch the generated code.`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		patch.DeleteJSONOmitemptyCommand(),
		patch.GenerateDBInitCommand(),
		patch.GenTypesPbCommand(),
		patch.CopyProtoCommand(),
		patch.CopyThirdPartyProtoCommand(),
		patch.CopyGOModCommand(),
		patch.ModifyDuplicateNumCommand(),
		patch.ModifyDuplicateErrCodeCommand(),
		patch.AdaptMonoRepoCommand(),
		patch.ModifyProtoPackageCommand(),
	)

	return cmd
}
