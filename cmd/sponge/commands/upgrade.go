package commands

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/gobash"
	"github.com/go-dev-frame/sponge/pkg/gofile"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

// UpgradeCommand upgrade sponge binaries
func UpgradeCommand() *cobra.Command {
	var targetVersion string

	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade sponge version",
		Long:  "Upgrade sponge version.",
		Example: color.HiBlackString(`  # Upgrade to latest version
  sponge upgrade

  # Upgrade to specified version
  sponge upgrade --version=v1.5.6`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
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
	runningTip := "upgrading sponge binary "
	finishTip := "upgrade sponge binary done " + installedSymbol
	failTip := "upgrade sponge binary failed " + lackSymbol
	p := utils.NewWaitPrinter(time.Millisecond * 500)
	p.LoopPrint(runningTip)
	err := runUpgradeCommand(targetVersion)
	if err != nil {
		p.StopPrint(failTip)
		return "", err
	}
	p.StopPrint(finishTip)

	runningTip = "upgrading template code "
	finishTip = "upgrade template code done " + installedSymbol
	failTip = "upgrade template code failed " + lackSymbol
	p = utils.NewWaitPrinter(time.Millisecond * 500)
	p.LoopPrint(runningTip)
	ver, err := copyToTempDir(targetVersion)
	if err != nil {
		p.StopPrint(failTip)
		return "", err
	}
	p.StopPrint(finishTip)

	runningTip = "upgrading the built-in plugins of sponge "
	finishTip = "upgrade the built-in plugins of sponge done " + installedSymbol
	failTip = "upgrade the built-in plugins of sponge failed " + lackSymbol
	p = utils.NewWaitPrinter(time.Millisecond * 500)
	p.LoopPrint(runningTip)
	err = updateSpongeInternalPlugin(ver)
	if err != nil {
		p.StopPrint(failTip)
		return "", err
	}
	p.StopPrint(finishTip)
	return ver, nil
}

func runUpgradeCommand(targetVersion string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3) //nolint
	spongeVersion := "github.com/go-dev-frame/sponge/cmd/sponge@" + targetVersion
	if targetVersion != latestVersion && targetVersion < "v1.11.2" {
		spongeVersion = strings.ReplaceAll(spongeVersion, "go-dev-frame", "zhufuyi")
	}
	result := gobash.Run(ctx, "go", "install", spongeVersion)
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
		arg := fmt.Sprintf("%s/pkg/mod/github.com/go-dev-frame", gopath)
		result, err = gobash.Exec("ls", adaptPathDelimiter(arg))
		if err != nil {
			return "", fmt.Errorf("execute command failed, %v", err)
		}

		spongeDirName = getLatestVersion(string(result))
		if spongeDirName == "" {
			return "", fmt.Errorf("not found sponge directory in '$GOPATH/pkg/mod/github.com/go-dev-frame'")
		}
	} else {
		spongeDirName = "sponge@" + targetVersion
	}

	srcDir := adaptPathDelimiter(fmt.Sprintf("%s/pkg/mod/github.com/go-dev-frame/%s", gopath, spongeDirName))
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
	_ = executeCommand("rm", "-rf", targetDir+"/cmd/protoc-gen-go-gin")
	_ = executeCommand("rm", "-rf", targetDir+"/cmd/protoc-gen-go-rpc-tmpl")
	_ = executeCommand("rm", "-rf", targetDir+"/cmd/protoc-gen-json-field")
	_ = executeCommand("rm", "-rf", targetDir+"/pkg")
	_ = executeCommand("rm", "-rf", targetDir+"/test")
	_ = executeCommand("rm", "-rf", targetDir+"/assets")

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
	genGinVersion := "github.com/go-dev-frame/sponge/cmd/protoc-gen-go-gin@" + targetVersion
	if targetVersion < "v1.11.2" {
		genGinVersion = strings.ReplaceAll(genGinVersion, "go-dev-frame", "zhufuyi")
	}
	result := gobash.Run(ctx, "go", "install", genGinVersion)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return result.Err
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Minute) //nolint
	genRPCVersion := "github.com/go-dev-frame/sponge/cmd/protoc-gen-go-rpc-tmpl@" + targetVersion
	if targetVersion < "v1.11.2" {
		genRPCVersion = strings.ReplaceAll(genRPCVersion, "go-dev-frame", "zhufuyi")
	}
	result = gobash.Run(ctx, "go", "install", genRPCVersion)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		return result.Err
	}

	// v1.x.x version does not support protoc-gen-json-field
	if !strings.HasPrefix(targetVersion, "v1") {
		ctx, _ = context.WithTimeout(context.Background(), time.Minute) //nolint
		genJSONVersion := "github.com/go-dev-frame/sponge/cmd/protoc-gen-json-field@" + targetVersion
		if targetVersion < "v1.11.2" {
			genJSONVersion = strings.ReplaceAll(genJSONVersion, "go-dev-frame", "zhufuyi")
		}
		result = gobash.Run(ctx, "go", "install", genJSONVersion)
		for v := range result.StdOut {
			_ = v
		}
		if result.Err != nil {
			return result.Err
		}
	}

	return nil
}
