package commands

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gobash"
)

// GenGraphCommand generate graph command
func GenGraphCommand() *cobra.Command {
	var (
		isAll      bool
		projectDir string
		serverDir  []string
	)

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Draw a business architecture diagram based on the project created by sponge",
		Long:  "Draw a business architecture diagram based on the project created by sponge.",
		Example: color.HiBlackString(`  # If there are multiple servers in a project, simply specify the project directory path to generate a diagram between servers 
  sponge graph --project-dir=/path/to/project

  # You can also specify multiple services to generate a business framework diagram
  sponge graph --server-dir=/path/to/server1 --server-dir=/path/to/server2

  # Includes database related servers
  sponge graph --project-dir=/path/to/project --all`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectDir == "" && len(serverDir) == 0 {
				return errors.New("no project directory or server directory specified\n\n" + cmd.Example)
			}

			_, err := gobash.Exec("spograph", "-h")
			if err != nil {
				fmt.Printf("not found spograph command, please install it by running the following command: %s\n",
					color.HiCyanString("go install github.com/zhufuyi/spograph@latest"))
				return nil
			}

			var params []string
			if projectDir != "" {
				params = append(params, "--project-dir="+projectDir)
			}
			for _, dir := range serverDir {
				params = append(params, "--server-dir="+dir)
			}
			if isAll {
				params = append(params, "--all")
			}
			result, err := gobash.Exec("spograph", params...)
			if err != nil {
				return err
			}
			fmt.Printf("%s", string(result))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&isAll, "all", "a", false, "does it include services such as databases")
	cmd.Flags().StringVarP(&projectDir, "project-dir", "p", "", "project directory")
	cmd.Flags().StringSliceVarP(&serverDir, "server-dir", "s", []string{}, "server directory, multiple parameters can be set")

	return cmd
}
