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
	var targetVersion string

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade sponge version",
		Long: `upgrade sponge version.

Examples:
  # upgrade to latest version
  sponge upgrade
  # upgrade to specified version
  sponge upgrade --version=v1.5.6
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("upgrade sponge ......")
			if targetVersion == "" {
				targetVersion = latestVersion
			}
			ver, err := runUpgrade(targetVersion)
			if err != nil {
				return err
			}
			fmt.Printf("upgraded version to %s successfully.\n", ver)
			return nil
		},
	}

	cmd.Flags().StringVarP(&targetVersion, "version", "v", latestVersion, "upgrade sponge version")
	return cmd
}

func runUpgrade(targetVersion string) (string, error) {
	err := runUpgradeCommand(targetVersion)
	if err != nil {
		fmt.Println(lackSymbol + "upgrade sponge binary.")
		return "", err
	}
	fmt.Println(isntalledSymbol + "upgraded sponge binary.")
	ver, err := copyToTempDir(targetVersion)
	if err != nil {
		fmt.Println(lackSymbol + "upgrade template code.")
		return "", err
	}
	fmt.Println(isntalledSymbol + "upgraded template code.")
	err = updateSpongeInternalPlugin(ver)
	if err != nil {
		fmt.Println(lackSymbol + "upgrade protoc plugins.")
		return "", err
	}
	fmt.Println(isntalledSymbol + "upgraded protoc plugins.")
	return ver, nil
}

func runUpgradeCommand(targetVersion string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3) //nolint
	result := gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/sponge@"+targetVersion)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// copy the template files to a temporary directory
func copyToTempDir(targetVersion string) (string, error) {
	result, err := gobash.Exec("go", "env", "GOPATH")
	if err != nil {
		return "", fmt.Errorf("execute command failed, %v", err)
	}
	gopath := strings.ReplaceAll(string(result), "\n", "")
	if gopath == "" {
		return "", fmt.Errorf("$GOPATH is empty, you need set $GOPATH in your $PATH")
	}

	spongeDirName := ""
	if targetVersion == latestVersion {
		// find the new version of the sponge code directory
		arg := fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi", gopath)
		result, err = gobash.Exec("ls", adaptPathDelimiter(arg))
		if err != nil {
			return "", fmt.Errorf("execute command failed, %v", err)
		}

		spongeDirName = getLatestVersion(string(result))
		if spongeDirName == "" {
			return "", fmt.Errorf("not found sponge directory in '$GOPATH/pkg/mod/github.com/zhufuyi'")
		}
	} else {
		spongeDirName = "sponge@" + targetVersion
	}

	srcDir := adaptPathDelimiter(fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi/%s", gopath, spongeDirName))
	destDir := adaptPathDelimiter(GetSpongeDir() + "/")
	targetDir := adaptPathDelimiter(destDir + ".sponge")

	err = executeCommand("rm", "-rf", targetDir)
	if err != nil {
		return "", err
	}
	err = executeCommand("cp", "-rf", srcDir, targetDir)
	if err != nil {
		return "", err
	}
	err = executeCommand("chmod", "-R", "744", targetDir)
	if err != nil {
		return "", err
	}
	_ = executeCommand("rm", "-rf", targetDir+"/cmd/sponge")

	versionNum := strings.Replace(spongeDirName, "sponge@", "", 1)
	err = os.WriteFile(versionFile, []byte(versionNum), 0644)
	if err != nil {
		return "", err
	}

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

func updateSpongeInternalPlugin(targetVersion string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute) //nolint
	result := gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@"+targetVersion)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return result.Err
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Minute) //nolint
	result = gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@"+targetVersion)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return result.Err
	}

	return nil
}
