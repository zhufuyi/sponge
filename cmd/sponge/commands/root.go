package commands

import (
	"os"

	"github.com/zhufuyi/sponge/cmd/sponge/commands/generate"

	"github.com/spf13/cobra"
)

var (
	version     = "v0.0.0"
	versionFile = os.TempDir() + "/sponge/.github/version"
)

// NewRootCMD command entry
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "sponge",
		Long:          "sponge is a microservices framework for quickly creating http or grpc code.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       getVersion(),
	}

	cmd.AddCommand(
		InitCommand(),
		UpdateCommand(),
		ToolsCommand(),
		NewWebCommand(),
		MicroCommand(),
		generate.ConfigCommand(),
	)

	return cmd
}

func getVersion() string {
	data, _ := os.ReadFile(versionFile)
	v := string(data)
	if v != "" {
		return v
	}
	return version
}
