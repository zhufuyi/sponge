// Package patch is command set for patching service code.
package patch

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

var copyCount = 0

// CopyProtoCommand copy proto file from the grpc service directory
func CopyProtoCommand() *cobra.Command {
	var (
		serverDir     string // server dir
		versionFolder string // proto file version folder
		outPath       string // output directory
		protoFile     string // proto file names
		targetModule  string // target module
	)

	cmd := &cobra.Command{
		Use:   "copy-proto",
		Short: "Copy proto file from the grpc service directory",
		Long: `Copy proto file from the grpc service, if the proto file exists, it will be forced to overwrite it,
don't worry about losing the proto file after overwriting it, before copying proto it will be backed up to 
the directory /tmp/sponge_copy_backup_proto_files.`,
		Example: color.HiBlackString(`  # Copy all proto files from a grpc service directory
  sponge patch copy-proto --server-dir=../grpc-server

  # Copy all proto files from multiple grpc services directory
  sponge patch copy-proto --server-dir=../grpc-server1,../rpc-server2

  # Copy the specified proto files in the grpc service directory
  sponge patch copy-proto --server-dir=../grpc-server --proto-file=name1.proto,name2.proto`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !gofile.IsExists(outPath) {
				_ = os.MkdirAll(outPath, 0766)
			}

			if targetModule == "" {
				moduleName, _, err := getModuleAndServerName(".")
				if err != nil {
					return err
				}
				targetModule = moduleName
			}

			var selectProtoFiles []string
			if protoFile != "" {
				selectProtoFiles = strings.Split(protoFile, ",")
			}

			serverDirs := strings.Split(serverDir, ",")
			for _, srcDir := range serverDirs {
				srcModuleName, srcServerName, err := getModuleAndServerName(srcDir)
				if err != nil {
					return err
				}
				pc := &protoCopier{
					moduleName:       targetModule,
					outPath:          outPath,
					srcModuleName:    srcModuleName,
					srcServerName:    srcServerName,
					srcDir:           srcDir,
					srcVersionFolder: versionFolder,
					selectProtoFiles: selectProtoFiles,
					copiedFiles:      make(map[string]struct{}),
				}
				err = pc.copyProtoFiles()
				if err != nil {
					return err
				}
			}

			if copyCount == 0 {
				fmt.Printf("\nno proto files to copy, server-dir = %v\n", serverDirs)
			} else {
				fmt.Printf("\ncopy proto files successfully, out = %s\n", outPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&serverDir, "server-dir", "s", "", "server directory, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("server-dir")
	cmd.Flags().StringVarP(&protoFile, "proto-file", "p", "", "proto files, multiple names separated by commas")
	cmd.Flags().StringVarP(&targetModule, "target-module", "t", "", "target module name, same module name as the target project's go.mod")
	cmd.Flags().StringVarP(&versionFolder, "version-folder", "v", "v1", "proto file version folder")
	cmd.Flags().StringVarP(&outPath, "out", "o", "api", "output directory, if the proto file already exists, it will be overwritten directly")
	return cmd
}

func getModuleAndServerName(dir string) (moduleName string, serverName string, err error) {
	if dir == "" {
		return "", "", errors.New("param \"server-dir\" is empty")
	}
	data, err := os.ReadFile(dir + "/docs/gen.info")
	if err != nil {
		return "", "", err
	}

	ms := strings.Split(string(data), ",")
	if len(ms) >= 2 {
		return ms[0], ms[1], nil
	}

	return "", "", errors.New("not found server name in docs/gen.info")
}

type protoCopier struct {
	moduleName string
	outPath    string

	srcModuleName    string
	srcServerName    string
	srcDir           string
	srcVersionFolder string

	selectProtoFiles []string
	copiedFiles      map[string]struct{}
}

func (c *protoCopier) copyProtoFiles() error {
	srcProtoFolder := c.srcDir + "/api/" + c.srcServerName + "/" + c.srcVersionFolder
	targetProtoFolder := c.outPath + "/" + c.srcServerName + "/" + c.srcVersionFolder

	protoFiles, err := gofile.ListFiles(srcProtoFolder, gofile.WithSuffix(".proto"))
	if err != nil {
		return err
	}

	// match proto files
	if len(c.selectProtoFiles) > 0 {
		var matchProtoFiles []string
		for _, filePath := range protoFiles {
			for _, file := range c.selectProtoFiles {
				if gofile.GetFilename(filePath) == file {
					matchProtoFiles = append(matchProtoFiles, filePath)
					break
				}
			}
		}
		if len(matchProtoFiles) == 0 {
			return nil
		}
		protoFiles = matchProtoFiles
	}

	if len(protoFiles) > 0 {
		err = backupProtoFiles(c.outPath)
		if err != nil {
			return err
		}
		fmt.Println()
	}

	for _, pf := range protoFiles {
		data, err := os.ReadFile(pf)
		if err != nil {
			return err
		}
		err = c.copyDependencyProtoFile(data)
		if err != nil {
			return err
		}
		targetProtoFile := targetProtoFolder + "/" + gofile.GetFilename(pf)
		err = c.copyProtoFile(pf, targetProtoFile, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *protoCopier) copyDependencyProtoFile(data []byte) error {
	regStr := `import(.*?)"api/([\w\W]*?.proto)`
	reg := regexp.MustCompile(regStr)
	match := reg.FindAllStringSubmatch(string(data), -1)
	if len(match) == 0 {
		return nil
	}

	var pfs []string
	for _, v := range match {
		if len(v) == 3 {
			pfs = append(pfs, v[2])
		}
	}

	for _, pf := range pfs {
		srcProtoFile := c.srcDir + "/api/" + pf
		targetProtoFile := c.outPath + "/" + pf
		pData, err := os.ReadFile(srcProtoFile)
		if err != nil {
			return err
		}
		err = c.copyProtoFile(srcProtoFile, targetProtoFile, true)
		if err != nil {
			return err
		}

		if copyCount > 1000 {
			return errors.New("import dependencies circle or too many files")
		}
		err = c.copyDependencyProtoFile(pData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *protoCopier) copyProtoFile(srcProtoFile string, targetProtoFile string, isDependency bool) error {
	if c.isCopied(targetProtoFile) {
		return nil
	}

	targetProtoDir := gofile.GetFileDir(targetProtoFile)
	_ = os.MkdirAll(targetProtoDir, 0766)

	// replace go_package
	pbContent, err := os.ReadFile(srcProtoFile)
	if err != nil {
		return fmt.Errorf("read file %s error, %v", srcProtoFile, err)
	}
	pbContent = c.replacePackage(pbContent, isDependency)

	tmpFile := os.TempDir() + gofile.GetPathDelimiter() + gofile.GetFilename(srcProtoFile)
	err = os.WriteFile(tmpFile, pbContent, 0666)
	if err != nil {
		return err
	}

	fmt.Printf("    %s  -->  %s\n", cutPath(srcProtoFile), targetProtoFile)
	_, err = gobash.Exec("mv", "-f", tmpFile, targetProtoFile)
	if err != nil {
		return err
	}
	copyCount++
	return nil
}

func (c *protoCopier) replacePackage(data []byte, isDependency bool) []byte {
	if bytes.Contains(data, []byte("\r\n")) {
		data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	}

	regStr2 := `go_package [\w\W]*?;\n`
	reg2 := regexp.MustCompile(regStr2)
	goPackageName := reg2.Find(data)

	if len(goPackageName) > 0 {
		if isDependency {
			newGoPackageName := bytes.Replace(goPackageName, []byte(c.srcModuleName), []byte(c.moduleName), 1)
			data = bytes.Replace(data, goPackageName, newGoPackageName, 1)
		} else {
			newGoPackage := fmt.Sprintf("go_package = \"%s/api/%s/%s;%s\";\n", c.moduleName, c.srcServerName, c.srcVersionFolder, c.srcVersionFolder)
			data = bytes.Replace(data, goPackageName, []byte(newGoPackage), 1)
		}
	}

	return data
}

func (c *protoCopier) isCopied(targetProtoFile string) bool {
	if _, ok := c.copiedFiles[targetProtoFile]; !ok {
		c.copiedFiles[targetProtoFile] = struct{}{}
		return false
	}
	return true
}

func backupProtoFiles(outPath string) error {
	prefixPath, err := filepath.Abs(outPath)
	if err != nil {
		return err
	}

	pfs, _ := gofile.ListFiles(outPath, gofile.WithSuffix(".proto"))
	backupDir := os.TempDir() + gofile.GetPathDelimiter() + "sponge_copy_backup_proto_files" +
		gofile.GetPathDelimiter() + time.Now().Format("20060102T150405")
	for _, srcProtoFile := range pfs {
		suffixPath := strings.ReplaceAll(srcProtoFile, prefixPath, "")
		targetProtoFile := backupDir + suffixPath
		targetProtoDir := gofile.GetFileDir(targetProtoFile)

		_ = os.MkdirAll(targetProtoDir, 0744)
		_, err = gobash.Exec("cp", srcProtoFile, targetProtoFile)
		if err != nil {
			return err
		}
	}
	return nil
}
