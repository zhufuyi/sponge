package commands

import (
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands/generate"
)

// GenWebCommand generate web server code
func GenWebCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "web",
		Short:         "Generate model, cache, dao, handler, http code",
		Long:          "Generate model, cache, dao, handler, http code.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		generate.ModelCommand("web"),
		generate.DaoCommand("web"),
		generate.CacheCommand("web"),
		generate.HandlerCommand(),
		generate.HTTPCommand(),
		generate.HTTPPbCommand(),
		generate.ConvertSwagJSONCommand("web"),
		generate.HandlerPbCommand(),
	)

	return cmd
}
