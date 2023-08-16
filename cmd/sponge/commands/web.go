package commands

import (
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// NewWebCommand web commands
func NewWebCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "web",
		Short:         "Generate model, cache, dao, handler, http codes",
		Long:          "generate model, cache, dao, handler, http codes.",
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
