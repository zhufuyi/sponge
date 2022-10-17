package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"

	"github.com/spf13/cobra"
)

// UpdateCommand update sponge binaries
func UpdateCommand() *cobra.Command {
	var executor string
	var enableCNGoProxy bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update sponge to the latest version",
		Long: `update sponge to the latest version.

Examples:
  # for linux
  sponge update
  # for windows, need to specify the bash file
  sponge update --executor="D:\Program Files\cmder\vendor\git-for-windows\bin\bash.exe"
  # use https://goproxy.cn goproxy
  sponge update -g
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if executor != "" {
				gobash.SetExecutorPath(executor)
			}
			err := runUpdateCommand(enableCNGoProxy)
			if err != nil {
				return err
			}
			err = CopyToTempDir()
			if err != nil {
				return err
			}
			fmt.Println("update sponge successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&executor, "executor", "e", "", "for windows systems, you need to specify the bash executor path.")
	cmd.Flags().BoolVarP(&enableCNGoProxy, "enable-cn-goproxy", "g", false, "is $GOPROXY turn on 'https://goproxy.cn'")

	return cmd
}

func runUpdateCommand(enableCNGoProxy bool) error {
	fmt.Println("update sponge ......")
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*5) //nolint
	command := "go install github.com/zhufuyi/sponge/cmd/sponge@latest"
	if enableCNGoProxy {
		command = "GOPROXY=https://goproxy.cn,direct && " + command
	}
	result := gobash.Run(ctx, command)
	for v := range result.StdOut {
		fmt.Printf("%s", v)
	}
	if result.Err != nil {
		return fmt.Errorf("exec command failed, %v", result.Err)
	}

	return nil
}

// CopyToTempDir 复制模板文件到临时目录下
func CopyToTempDir() error {
	result, err := gobash.Exec("go env GOPATH")
	if err != nil {
		return fmt.Errorf("Exec() error %v", err)
	}
	gopath := strings.ReplaceAll(string(result), "\n", "")
	if gopath == "" {
		return fmt.Errorf("$GOPATH is empty, you need set $GOPATH in your $PATH")
	}

	// 找出新版本sponge代码文件夹
	command := "ls $(go env GOPATH)/pkg/mod/github.com/zhufuyi | grep sponge@ | sort -r | head -1"
	result, err = gobash.Exec(command)
	if err != nil {
		return fmt.Errorf("Exec() error %v", err)
	}
	latestSpongeDirName := strings.ReplaceAll(string(result), "\n", "")
	if latestSpongeDirName == "" {
		return fmt.Errorf("not found 'sponge' directory in '$GOPATH/pkg/mod/github.com/zhufuyi'")
	}
	srcDir := fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi/%s", gopath, latestSpongeDirName)
	destDir := os.TempDir() + "/sponge"

	// 复制到临时目录
	_ = os.RemoveAll(adaptPathDelimiter(destDir))
	command = fmt.Sprintf(`cp -rf %s %s`, adaptPathDelimiter(srcDir), adaptPathDelimiter(destDir))
	_, err = gobash.Exec(command)
	return err
}

func adaptPathDelimiter(filePath string) string {
	if gofile.IsWindows() {
		filePath = strings.ReplaceAll(filePath, "\\", "/")
	}
	return filePath
}
