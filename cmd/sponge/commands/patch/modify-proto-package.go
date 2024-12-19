package patch

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/gofile"
)

// ModifyProtoPackageCommand modifies the package and go_package names of proto files.
func ModifyProtoPackageCommand() *cobra.Command {
	var (
		dir        string
		moduleName string
		serverDir  string
	)

	cmd := &cobra.Command{
		Use:   "modify-proto-package",
		Short: "Modifies the package and go_package names of proto files",
		Long:  "Modifies the package and go_package names of proto files.",
		Example: color.HiBlackString(`  # Modify the package and go_package names of all proto files in the api directory.
  sponge patch modify-proto-package --dir=api --module-name=foo

  # Modify the package and go_package names of all proto files in the api directory, get module name from docs/gen.
  sponge patch modify-proto-package --dir=api --server-dir=server`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverDir != "" {
				mdName, _, _ := getNamesFromOutDir(serverDir)
				if mdName != "" {
					moduleName = mdName
				}
			}

			if moduleName == "" {
				return errors.New("'module-name' is required")
			}

			protoFiles, err := gofile.ListFiles(dir, gofile.WithSuffix(".proto"), gofile.WithNoAbsolutePath())
			if err != nil {
				return err
			}
			if len(protoFiles) == 0 {
				fmt.Printf("no proto files found in the directory '%s'.\n", dir)
				return nil
			}

			var successFiles []string
			for _, file := range protoFiles {
				ss := splitProtoFilePath(gofile.GetDir(file))
				packageName, goPackageName := getPackageName(ss, moduleName)
				err = replaceProtoPackages(file, packageName, goPackageName)
				if err != nil {
					return err
				}
				successFiles = append(successFiles, file)
			}

			if len(successFiles) > 0 {
				fmt.Printf(`modified the package and go_package names of files:
    %s`, strings.Join(successFiles, "\n    "))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "input specified directory")
	_ = cmd.MarkFlagRequired("dir")
	cmd.Flags().StringVarP(&serverDir, "server-dir", "s", "", "server directory, get module name and server name from docs/gen.info")
	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "go module name")

	return cmd
}

func getPackageName(ss []string, moduleName string) (packageName string, goPackageName string) {
	l := len(ss)
	switch l {
	case 0:
		packageName = "v1"
		goPackageName = `"v1"`
		return packageName, goPackageName
	case 1:
		if ss[0] == "." {
			ss[0] = "v1"
		}
		packageName = ss[0]
		goPackageName = fmt.Sprintf(`"%s/%s;%s"`, moduleName, ss[0], ss[0])
		return packageName, goPackageName
	case 2:
		packageName = strings.Join(ss, ".")
		goPackageName = fmt.Sprintf(`"%s/%s;%s"`, moduleName, strings.Join(ss, "/"), ss[1])
		return packageName, goPackageName
	}
	packageName = strings.Join(ss[l-3:], ".")
	goPackageName = fmt.Sprintf(`"%s/%s;%s"`, moduleName, strings.Join(ss, "/"), ss[l-1])
	return packageName, goPackageName
}

func splitProtoFilePath(protoFilePath string) []string {
	ss := strings.Split(protoFilePath, gofile.GetPathDelimiter())
	if len(ss) > 0 {
		if ss[0] == ".." || ss[0] == "." {
			return ss[1:]
		}
	}
	return ss
}

func replaceProtoPackages(protoFilePath, packageName, goPackage string) error {
	data, err := os.ReadFile(protoFilePath)
	if err != nil {
		return err
	}

	if bytes.Contains(data, []byte("\r\n")) {
		data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	}

	regStr2 := `go_package [\w\W]*?;\n`
	reg2 := regexp.MustCompile(regStr2)
	srcGoPackageName := reg2.Find(data)
	newGoPackage := fmt.Sprintf("go_package = %s;\n", goPackage)
	if len(srcGoPackageName) > 0 {
		data = bytes.Replace(data, srcGoPackageName, []byte(newGoPackage), 1)
	} else {
		data = bytes.Replace(data, []byte("\n\n"), []byte("\n\n"+newGoPackage+"\n\n"), 1)
	}

	regStr := `\npackage [\w\W]*?;`
	reg := regexp.MustCompile(regStr)
	srcPackageName := reg.Find(data)
	newPackage := fmt.Sprintf("\npackage %s;", packageName)
	if len(srcPackageName) > 0 {
		data = bytes.Replace(data, srcPackageName, []byte(newPackage), 1)
	} else {
		data = bytes.Replace(data, []byte("\n\n"), []byte("\n\n"+newPackage+"\n\n"), 1)
	}

	return os.WriteFile(protoFilePath, data, 0666)
}
