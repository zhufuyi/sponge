package commands

import (
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// MicroCommand micro commands
func MicroCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "micro",
		Short:         "Generate proto, model, dao, service, rpc, rpc-gw codes",
		Long:          "generate proto, model, dao, service, rpc, rpc-gw codes.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		generate.ProtoBufCommand(),
		generate.ModelCommand("micro"),
		generate.DaoCommand("micro"),
		generate.ServiceCommand(),
		generate.RPCCommand(),
		generate.RPCGwPbCommand(),
		generate.RPCPbCommand(),
	)

	return cmd
}
