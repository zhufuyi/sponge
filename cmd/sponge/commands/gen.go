package commands

import (
	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

// GenCommand generate dependency code
func GenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gen",
		Short:         "Generate dependency code, e.g. mysql and redis initialization code, types.proto",
		Long:          `generate dependency code.`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(
		generate.MysqlAndRedisCommand(),
		generate.TypesPbCommand(),
	)
	return cmd
}
