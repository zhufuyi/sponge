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

	"github.com/spf13/cobra"
)

// UpdateCommand update sponge binaries
func UpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update sponge to the latest version",
		Long: `update sponge to the latest version.

Examples:
  # run update
  sponge update
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("update sponge ......")
			err := runUpdateCommand()
			if err != nil {
				return err
			}
			ver, err := copyToTempDir()
			if err != nil {
				return err
			}
			fmt.Printf("update sponge version to %s successfully.\n", ver)
			return nil
		},
	}

	return cmd
}

func runUpdateCommand() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3) //nolint
	result := gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/sponge@latest")
	for range result.StdOut {
	}
	if result.Err != nil {
		return fmt.Errorf("exec command failed, %v", result.Err)
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Minute) //nolint
	result = gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest")
	for range result.StdOut {
	}
	if result.Err != nil {
		return fmt.Errorf("exec command failed, %v", result.Err)
	}

	ctx, _ = context.WithTimeout(context.Background(), time.Minute) //nolint
	result = gobash.Run(ctx, "go", "install", "github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@latest")
	for range result.StdOut {
	}
	if result.Err != nil {
		return fmt.Errorf("exec command failed, %v", result.Err)
	}

	return nil
}

// copy the template files to a temporary directory
func copyToTempDir() (string, error) {
	result, err := gobash.Exec("go", "env", "GOPATH")
	if err != nil {
		return "", fmt.Errorf("cxec command failed, %v", err)
	}
	gopath := strings.ReplaceAll(string(result), "\n", "")
	if gopath == "" {
		return "", fmt.Errorf("$GOPATH is empty, you need set $GOPATH in your $PATH")
	}

	// find the new version of the sponge code directory
	arg := fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi", gopath)
	result, err = gobash.Exec("ls", adaptPathDelimiter(arg))
	if err != nil {
		return "", fmt.Errorf("cxec command failed, %v", err)
	}

	latestSpongeDirName := getLatestVersion(string(result))
	if latestSpongeDirName == "" {
		return "", fmt.Errorf("not found 'sponge' directory in '$GOPATH/pkg/mod/github.com/zhufuyi'")
	}

	srcDir := adaptPathDelimiter(fmt.Sprintf("%s/pkg/mod/github.com/zhufuyi/%s", gopath, latestSpongeDirName))
	destDir := adaptPathDelimiter(os.TempDir() + "/sponge")
	destDirBk := destDir + ".bk"

	// copy to temporary directory
	_ = os.Rename(destDir, destDirBk)
	_, err = gobash.Exec("cp", "-rf", srcDir, destDir)
	if err != nil {
		_ = os.Rename(destDirBk, destDir)
		return "", fmt.Errorf("cxec command failed, %v", err)
	}
	_ = os.RemoveAll(destDirBk)

	versionNum := strings.Replace(latestSpongeDirName, "sponge@", "", 1)
	_ = os.WriteFile(versionFile, []byte(versionNum), 0666)
	return versionNum, nil
}

func adaptPathDelimiter(filePath string) string {
	if gofile.IsWindows() {
		filePath = strings.ReplaceAll(filePath, "/", "\\")
	}
	return filePath
}

func getLatestVersion(s string) string {
	dirs := strings.Split(s, "\n")
	var allVersions []string
	for _, dirName := range dirs {
		if strings.Contains(dirName, "sponge@") {
			allVersions = append(allVersions, dirName)
		}
	}
	if allVersions == nil {
		return ""
	}
	sort.Strings(allVersions)
	return allVersions[len(allVersions)-1]
}
