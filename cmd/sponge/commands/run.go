package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/zhufuyi/sponge/cmd/sponge/server"

	"github.com/spf13/cobra"
)

// NewRunCommand sponge run commands
func NewRunCommand() *cobra.Command {
	var port = 24631

	cmd := &cobra.Command{
		Use:           "run",
		Short:         "Start a web service for sponge",
		Long:          "start a web service for sponge.",
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://localhost:%d", port)
			fmt.Printf("sponge service started, visit %s in your browser.\n\n", url)
			go func() {
				_ = open(url)
			}()
			server.RunHTTPServer(fmt.Sprintf(":%d", port))
			return nil
		},
	}

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
