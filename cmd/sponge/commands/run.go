package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/zhufuyi/sponge/cmd/sponge/server"

	"github.com/spf13/cobra"
)

var servicePort = 24631

// NewRunCommand sponge run commands
func NewRunCommand() *cobra.Command {
	var isLog bool
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Open the Sponge UI interface",
		Long: `open the Sponge UI interface.

Examples:
  # no log for running.
  sponge run

  # log for running.
  sponge run -l
`,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://localhost:%d", servicePort)
			fmt.Printf("sponge command ui service started successfully, visit %s in your browser.\n\n", url)
			go func() {
				_ = open(url)
			}()
			server.RunHTTPServer(fmt.Sprintf(":%d", servicePort), isLog)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&isLog, "log", "l", false, "enable service logging")
	return cmd
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
