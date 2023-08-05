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
	"protoc-gen-go":          "google.golang.org/protobuf/cmd/protoc-gen-go@latest",
	"protoc-gen-go-grpc":     "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
	"protoc-gen-validate":    "github.com/envoyproxy/protoc-gen-validate@latest",
	"protoc-gen-gotag":       "github.com/srikrsna/protoc-gen-gotag@latest",
	"protoc-gen-go-gin":      "github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest",
	"protoc-gen-go-rpc-tmpl": "github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@latest",
	"protoc-gen-openapiv2":   "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
	"protoc-gen-doc":         "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest",
	"golangci-lint":          "github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	"swag":                   "github.com/swaggo/swag/cmd/swag@v1.8.12",
	"go-callvis":             "github.com/ofabry/go-callvis@latest",
}

const (
	isntalledSymbol = "✔ "
	lackSymbol      = "❌ "
	warnSymbol      = "⚠ "
)

// ToolsCommand tools management
func ToolsCommand() *cobra.Command {
	var installFlag bool

	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Managing sponge dependency tools",
		Long: `managing sponge dependency tools.

Examples:
  # show all dependencies tools.
  sponge tools

  # install all dependencies tools.
  sponge tools --install
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			installedNames, lackNames := checkInstallTools()
			if installFlag {
				installTools(lackNames)
			} else {
				showDependencyTools(installedNames, lackNames)
			}

			return nil
		},
	}
	cmd.Flags().BoolVarP(&installFlag, "install", "i", false, "install dependent tools")

	return cmd
}

func checkInstallTools() ([]string, []string) {
	var installedNames, lackNames = []string{}, []string{}
	for _, name := range toolNames {
		_, err := gobash.Exec("which", name)
		if err != nil {
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
		content += "\nInstalling dependency tools using the command: sponge tools --install\n"
	} else {
		content += "\nAll dependent tools installed.\n"
	}

	fmt.Println(content)
}

func installTools(lackNames []string) {
	if len(lackNames) == 0 {
		fmt.Printf("\n    All dependent tools installed.\n\n")
		return
	}
	fmt.Printf("install a total of %d dependent tools, need to wait a little time.\n\n", len(lackNames))

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
			pkgAddr, ok := installToolCommands[name]
			if !ok {
				return
			}
			pkgAddr = adaptInternalCommand(name, pkgAddr)
			result := gobash.Run(ctx, "go", "install", pkgAddr)
			for v := range result.StdOut {
				_ = v
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
	fmt.Println()
}

func adaptInternalCommand(name string, pkgAddr string) string {
	if name == "protoc-gen-go-gin" || name == "protoc-gen-go-rpc-tmpl" {
		if version != "v0.0.0" {
			return strings.ReplaceAll(pkgAddr, "@latest", "@"+version)
		}
	}

	return pkgAddr
}
