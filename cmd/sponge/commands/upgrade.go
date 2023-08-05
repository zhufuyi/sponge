package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/spf13/cobra"
)

// UpgradeCommand upgrade sponge binaries
func UpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade sponge to the latest version",
		Long: `upgrade sponge to the latest version.

Examples:
  # upgrade version
  sponge upgrade
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("upgrade sponge ......")
			err := runUpgradeCommand()
			if err != nil {
				return err
			}
			ver, err := copyToTempDir()
			if err != nil {
				return err
			}
			updateSpongeInternalPlugin(ver)
			fmt.Printf("upgrade sponge version to %s successfully.\n", ver)
			return nil
		},
	}

	return cmd
}

func runUpgradeCommand() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3) //nolint
	result := gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/sponge@latest")
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return fmt.Errorf("execute command failed, %v", result.Err)
	}
	fmt.Println("already downloaded the latest version of sponge.")
	return nil
}

// copy the template files to a temporary directory
func copyToTempDir() (string, error) {
	result, err := gobash.Exec("go", "env", "GOPATH")
	if err != nil {
		return "", fmt.Errorf("execute command failed, %v", err)
	}
	gopath := strings.ReplaceAll(string(result), "\n", "")
	if gopath == "" {
		return "", fmt.Errorf("$GOPATH is empty, you need set $GOPATH in your $PATH")
	}

	// find the new version of the sponge code directory
	arg := fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi", gopath)
	result, err = gobash.Exec("ls", adaptPathDelimiter(arg))
	if err != nil {
		return "", fmt.Errorf("execute command failed, %v", err)
	}

	latestSpongeDirName := getLatestVersion(string(result))
	if latestSpongeDirName == "" {
		return "", fmt.Errorf("not found 'sponge' directory in '$GOPATH/pkg/mod/github.com/zhufuyi'")
	}

	srcDir := adaptPathDelimiter(fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi/%s", gopath, latestSpongeDirName))
	destDir := adaptPathDelimiter(GetSpongeDir() + "/")
	targetDir := adaptPathDelimiter(destDir + ".sponge")

	err = executeCommand("rm", "-rf", targetDir)
	if err != nil {
		return "", err
	}
	err = executeCommand("cp", "-rf", srcDir, destDir)
	if err != nil {
		return "", err
	}
	err = executeCommand("chmod", "-R", "744", destDir+latestSpongeDirName)
	if err != nil {
		return "", err
	}
	err = os.Rename(destDir+latestSpongeDirName, targetDir)
	if err != nil {
		return "", fmt.Errorf("rename '%s' error, %v", destDir, err)
	}

	versionNum := strings.Replace(latestSpongeDirName, "sponge@", "", 1)
	_ = os.WriteFile(versionFile, []byte(versionNum), 0666)

	fmt.Println("the code template has been updated to the latest version.")

	return versionNum, nil
}

func executeCommand(name string, args ...string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30) //nolint
	result := gobash.Run(ctx, name, args...)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return fmt.Errorf("execute command failed, %v", result.Err)
	}
	return nil
}

func adaptPathDelimiter(filePath string) string {
	if gofile.IsWindows() {
		filePath = strings.ReplaceAll(filePath, "/", "\\")
	}
	return filePath
}

func getLatestVersion(s string) string {
	var dirNames = make(map[int]string)
	var nums []int

	dirs := strings.Split(s, "\n")
	for _, dirName := range dirs {
		if strings.Contains(dirName, "sponge@") {
			tmp := strings.ReplaceAll(dirName, "sponge@", "")
			ss := strings.Split(tmp, ".")
			if len(ss) != 3 {
				continue
			}
			if strings.Contains(ss[2], "v0.0.0") {
				continue
			}
			num := utils.StrToInt(ss[0])*10000 + utils.StrToInt(ss[1])*100 + utils.StrToInt(ss[2])
			nums = append(nums, num)
			dirNames[num] = dirName
		}
	}
	if len(nums) == 0 {
		return ""
	}

	sort.Ints(nums)
	return dirNames[nums[len(nums)-1]]
}

func updateSpongeInternalPlugin(latestVersionNum string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute) //nolint
	result := gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@"+latestVersionNum)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		fmt.Printf("upgrade plugin 'protoc-gen-go-gin' failed, target version=%s, error=%v, "+
			"please use the command 'sponge tools --install' to retry.\n", latestVersionNum, result.Err)
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Minute) //nolint
	result = gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@"+latestVersionNum)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		fmt.Printf("upgrade plugin 'protoc-gen-go-rpc-tmpl' failed, target version=%s, error=%v, "+
			"please use the command 'sponge tools --install' to retry.\n", latestVersionNum, result.Err)
	}
}
