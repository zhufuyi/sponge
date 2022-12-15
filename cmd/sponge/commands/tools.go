package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/zhufuyi/sponge/pkg/gobash"

	"github.com/spf13/cobra"
)

var toolNames = []string{
	"go",
	"protoc",
	"protoc-gen-go",
	"protoc-gen-go-grpc",
	"protoc-gen-validate",
	"protoc-gen-gotag",
	"protoc-gen-go-gin",
	"protoc-gen-go-rpc-tmpl",
	"protoc-gen-openapiv2",
	"protoc-gen-doc",
	"swag",
	"golangci-lint",
	"go-callvis",
}

var installToolCommands = map[string]string{
	"go":                     "go: please install manually yourself, download url is https://go.dev/dl/ or https://golang.google.cn/dl/",
	"protoc":                 "protoc: please install manually yourself, download url is https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.3",
	"protoc-gen-go":          "go install google.golang.org/protobuf/cmd/protoc-gen-go@latest",
	"protoc-gen-go-grpc":     "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
	"protoc-gen-validate":    "go install github.com/envoyproxy/protoc-gen-validate@latest",
	"protoc-gen-gotag":       "go install github.com/srikrsna/protoc-gen-gotag@latest",
	"protoc-gen-go-gin":      "go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest",
	"protoc-gen-go-rpc-tmpl": "go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@latest",
	"protoc-gen-openapiv2":   "go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
	"protoc-gen-doc":         "go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest",
	"golangci-lint":          "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	"swag":                   "go install github.com/swaggo/swag/cmd/swag@latest",
	"go-callvis":             "go install github.com/ofabry/go-callvis@latest",
}

const (
	isntalledSymbol = "✔ "
	lackSymbol      = "❌ "
	warnSymbol      = "⚠ "
)

// ToolsCommand tools management
func ToolsCommand() *cobra.Command {
	var executor string
	var enableCNGoProxy bool
	var installFlag bool

	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Managing sponge dependency tools",
		Long: `managing sponge dependency tools.

Examples:
  # for linux, show all dependencies tools.
  sponge tools

  # for windows, show all dependencies tools.
  sponge tools --executor="D:\Program Files\cmder\vendor\git-for-windows\bin\bash.exe"

  # use goproxy https://goproxy.cn
  sponge tools -g

  # install all tools.
  sponge tools --install
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if executor != "" {
				gobash.SetExecutorPath(executor)
			}

			installedNames, lackNames := checkInstallTools()
			if installFlag {
				installTools(lackNames, enableCNGoProxy)
			} else {
				showDependencyTools(installedNames, lackNames)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&executor, "executor", "e", "", "for windows systems, you need to specify the bash executor path.")
	cmd.Flags().BoolVarP(&enableCNGoProxy, "enable-cn-goproxy", "g", false, "is $GOPROXY turn on 'https://goproxy.cn'")
	cmd.Flags().BoolVarP(&installFlag, "install", "i", false, "install dependent tools")

	return cmd
}

func checkInstallTools() ([]string, []string) {
	var installedNames, lackNames = []string{}, []string{}
	for _, name := range toolNames {
		command := "which " + name
		_, err := gobash.Exec(command)
		if err != nil {
			checkExit("tools", err)
			lackNames = append(lackNames, name)
			continue
		}
		installedNames = append(installedNames, name)
	}

	data, _ := os.ReadFile(versionFile)
	v := string(data)
	if v != "" {
		version = v
	}

	return installedNames, lackNames
}

func showDependencyTools(installedNames []string, lackNames []string) {
	var content string

	if len(installedNames) > 0 {
		content = "Installed dependency tools:\n"
		for _, name := range installedNames {
			content += "    " + isntalledSymbol + " " + name + "\n"
		}
	}

	if len(lackNames) > 0 {
		content += "\nUninstalled dependency tools:\n"
		for _, name := range lackNames {
			content += "    " + lackSymbol + " " + name + "\n"
		}
		content += "\nInstalling dependency tools using the command: sponge tools --install"
	} else {
		content += "\nAll dependent tools installed."
	}

	fmt.Println(content)
}

func installTools(lackNames []string, enableCNGoProxy bool) {
	if len(lackNames) == 0 {
		fmt.Printf("\n    All dependent tools installed.\n\n")
		return
	}
	fmt.Printf("install dependent tools ......\n\n")

	var wg = &sync.WaitGroup{}
	var manuallyNames []string
	for _, name := range lackNames {
		if name == "go" || name == "protoc" {
			manuallyNames = append(manuallyNames, name)
			continue
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			ctx, _ := context.WithTimeout(context.Background(), time.Minute*3) //nolint
			command, ok := installToolCommands[name]
			if !ok {
				return
			}
			command = adaptInternalCommand(name, command)
			fmt.Println(command)
			if enableCNGoProxy {
				command = "GOPROXY=https://goproxy.cn,direct && " + command
			}
			result := gobash.Run(ctx, command)
			for range result.StdOut {
			}
			if result.Err != nil {
				fmt.Printf("%s %s, %v\n", lackSymbol, name, result.Err)
			} else {
				fmt.Printf("%s %s\n", isntalledSymbol, name)
			}
		}(name)
	}

	wg.Wait()

	for _, name := range manuallyNames {
		fmt.Println(warnSymbol + " " + installToolCommands[name])
	}
}

func checkExit(name string, err error) {
	str := err.Error()
	if strings.Contains(str, "exec: ") && strings.Contains(str, "file does not exist") {
		fmt.Printf(err.Error()+", must specify the executor location, example:\n"+`    sponge %s --executor="your bash file path"
`, name)
		os.Exit(1)
	}
}

func adaptInternalCommand(name string, command string) string {
	if name == "protoc-gen-go-gin" || name == "protoc-gen-go-rpc-tmpl" {
		if version != "v0.0.0" {
			return strings.ReplaceAll(command, "@latest", "@"+version)
		}
	}

	return command
}
