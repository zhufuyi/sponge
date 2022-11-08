package commands

import (
	"fmt"

	"github.com/zhufuyi/sponge/pkg/gobash"

	"github.com/spf13/cobra"
)

// InitCommand initial sponge
func InitCommand() *cobra.Command {
	var executor string
	var enableCNGoProxy bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize sponge",
		Long: `initialize sponge.

Examples:
  # for linux.
  sponge init

  # for windows.
  sponge init --executor="D:\Program Files\cmder\vendor\git-for-windows\bin\bash.exe"

  # use goproxy https://goproxy.cn
  sponge init -g
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if executor != "" {
				gobash.SetExecutorPath(executor)
			}
			fmt.Println("initialize sponge codes ......")
			// 下载sponge模板代码
			err := runUpdateCommand(enableCNGoProxy)
			if err != nil {
				return err
			}
			_, err = copyToTempDir()
			if err != nil {
				return err
			}

			// 安装依赖插件
			_, lackNames := checkInstallTools()
			installTools(lackNames, enableCNGoProxy)

			return nil
		},
	}

	cmd.Flags().StringVarP(&executor, "executor", "e", "", "for windows systems, you need to specify the bash executor path.")
	cmd.Flags().BoolVarP(&enableCNGoProxy, "enable-cn-goproxy", "g", false, "is $GOPROXY turn on 'https://goproxy.cn'")

	return cmd
}
