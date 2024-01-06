// Package commands are subcommands of the sponge command.
package commands

import (
	"fmt"
	"os"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

var (
	version     = "v0.0.0"
	versionFile = GetSpongeDir() + "/.sponge/.github/version"
)

// NewRootCMD command entry
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sponge",
		Long: `Sponge is a powerful golang productivity tool that integrates automatic code generation, 
web and microservice framework, basic development framework.
repo: https://github.com/zhufuyi/sponge
docs: https://go-sponge.com`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       getVersion(),
	}

	cmd.AddCommand(
		InitCommand(),
		UpgradeCommand(),
		PluginsCommand(),
		GenWebCommand(),
		GenMicroCommand(),
		generate.ConfigCommand(),
		OpenUICommand(),
		MergeCommand(),
		PatchCommand(),
	)

	return cmd
}

func getVersion() string {
	data, _ := os.ReadFile(versionFile)
	v := string(data)
	if v != "" {
		return v
	}
	return "unknown, execute command \"sponge init\" to get version"
}

// GetSpongeDir get sponge home directory
func GetSpongeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("can't get home directory'")
		return ""
	}

	return dir
}
