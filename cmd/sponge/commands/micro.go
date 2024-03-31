package commands

import (
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"
)

// GenMicroCommand generate micro service code
func GenMicroCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "micro",
		Short:         "Generate protobuf, model, cache, dao, service, grpc, grpc-gw, grpc+http, grpc-cli code",
		Long:          "generate protobuf, model, cache, dao, service, grpc, grpc-gw, grpc+http, grpc-cli code.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		generate.ProtobufCommand(),
		generate.ModelCommand("micro"),
		generate.DaoCommand("micro"),
		generate.CacheCommand("micro"),
		generate.ServiceCommand(),
		generate.RPCCommand(),
		generate.RPCGwPbCommand(),
		generate.RPCPbCommand(),
		generate.GRPCConnectionCommand(),
		generate.ConvertSwagJSONCommand("micro"),
		generate.GRPCAndHTTPPbCommand(),
		generate.ServiceAndHandlerCRUDCommand(),
	)

	return cmd
}
