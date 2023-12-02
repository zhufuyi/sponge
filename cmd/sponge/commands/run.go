package commands

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/zhufuyi/sponge/cmd/sponge/server"

	"github.com/spf13/cobra"
)

// OpenUICommand open the sponge UI interface
func OpenUICommand() *cobra.Command {
	var (
		port       int
		spongeAddr string
		isLog      bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Open the sponge UI interface",
		Long: `open the sponge UI interface.

Examples:
  # no log for running.
  sponge run

  # log for running.
  sponge run -l
`,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			if spongeAddr == "" {
				spongeAddr = "http://localhost:24631"
			} else {
				if err := checkSpongeAddr(spongeAddr); err != nil {
					return err
				}
			}
			fmt.Printf(`sponge command ui service started successfully, verson is %s, listening port is %d,
visit %s in your browser.
`, getVersion(), port, spongeAddr)
			go func() {
				_ = open(spongeAddr)
			}()
			server.RunHTTPServer(spongeAddr, port, isLog)
			return nil
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 24631, "port on which the sponge service listens")
	cmd.Flags().StringVarP(&spongeAddr, "addr", "a", "", "address of the front-end page requesting the sponge service, e.g. http://192.168.1.10:24631 or https://go-sponge.com/ui")
	cmd.Flags().BoolVarP(&isLog, "log", "l", false, "enable service logging")
	return cmd
}

func open(visitURL string) error {
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

	args = append(args, visitURL)
	return exec.Command(cmd, args...).Start()
}

func checkSpongeAddr(spongeAddr string) error {
	paramErr := errors.New("the addr parameter is invalid,  e.g. sponge run --addr=http://192.168.1.10:24631")
	u, err := url.Parse(spongeAddr)
	if err != nil {
		return paramErr
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return paramErr
	}

	return nil
}
