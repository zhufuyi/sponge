package commands

import (
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// Version 命令版本号
const Version = "0.0.0"

// NewRootCMD 命令入口
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "sponge",
		Long:          "sponge management tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       Version,
	}

	cmd.AddCommand(
		generate.ModelCommand(),
		generate.DaoCommand(),
		generate.HandlerCommand(),
		generate.HTTPCommand(),
		generate.ProtoCommand(),
		generate.ServiceCommand(),
		generate.GRPCCommand(),
		generate.ConfigCommand(),
		UpdateCommand(),
	)
	return cmd
}
