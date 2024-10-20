package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const latestVersion = "latest"

// InitCommand initial sponge
func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize sponge",
		Long: color.HiBlackString(`initialize sponge.

Examples:
  # run init, download code and install plugins.
  sponge init
`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("initializing sponge, please wait a moment ......")

			targetVersion := latestVersion
			// download sponge template code
			_, err := runUpgrade(targetVersion)
			if err != nil {
				return err
			}

			// installing dependency plugins
			_, lackNames := checkInstallPlugins()
			installPlugins(lackNames)

			return nil
		},
	}

	return cmd
}
