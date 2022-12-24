package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zhufuyi/sponge/cmd/sponge/server"
	"github.com/zhufuyi/sponge/pkg/utils"
)

// NewRunCommand sponge run commands
func NewRunCommand() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:           "run",
		Short:         "Start a web server for sponge",
		Long:          "start a web server for sponge.",
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			if port == 0 {
				port, _ = utils.GetAvailablePort()
			}

			fmt.Printf("sponge server start up, port=%d\n\n", port)
			server.RunHTTPServer(fmt.Sprintf(":%d", port))

			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 0, "web server port")

	return cmd
}
