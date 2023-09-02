// Package patch is command set for patching server code.
package patch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"

	"github.com/spf13/cobra"
)

var copyCount = 0

// CopyProtoCommand copy proto file from the rpc server directory
func CopyProtoCommand() *cobra.Command {
	var (
		serverDir     string // server dir
		versionFolder string // proto file version folder
		outPath       string // output directory
	)

	cmd := &cobra.Command{
		Use:   "copy-proto",
		Short: "Copy proto file from the rpc server directory",
		Long: `copy proto file from the rpc server, if the proto file exists, it will be forced to overwrite it,
don't worry about losing the proto file after overwriting it, before copying proto it will be backed up to 
the directory /tmp/sponge_copy_backup_proto_files.

Examples:
  # copy proto file from a rpc server directory
  sponge patch copy-proto --server-dir=./rpc-server

  # copy proto file from multiple rpc servers directory
  sponge patch copy-proto --server-dir=./rpc-server1,./rpc-server2
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !gofile.IsExists(outPath) {
				_ = os.MkdirAll(outPath, 0766)
			}

			serverDirs := strings.Split(serverDir, ",")
			for _, dir := range serverDirs {
				sn, err := getServerName(dir)
				if err != nil {
					return err
				}
				err = copyProtoFiles(dir, sn, versionFolder, outPath)
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
	cmd.Flags().StringVarP(&versionFolder, "version-folder", "v", "v1", "proto file version folder")
	cmd.Flags().StringVarP(&outPath, "out", "o", "./api", "output directory, if the proto file already exists, it will be overwritten directly")
	return cmd
}

func getServerName(dir string) (string, error) {
	if dir == "" {
		return "", errors.New("param \"server-dir\" is empty")
	}
	data, err := os.ReadFile(dir + "/docs/gen.info")
	if err != nil {
		return "", err
	}

	ms := strings.Split(string(data), ",")
	if len(ms) != 2 {
		return "", errors.New("not found server name in docs/gen.info")
	}

	return ms[1], nil
}

func copyProtoFiles(dir string, serverName string, versionFolder string, outPath string) error {
	srcProtoFolder := dir + "/api/" + serverName + "/" + versionFolder
	targetProtoFolder := outPath + "/" + serverName + "/" + versionFolder

	protoFiles, err := gofile.ListFiles(srcProtoFolder, gofile.WithSuffix(".proto"))
	if err != nil {
		return err
	}

	if len(protoFiles) > 0 {
		err = backupProtoFiles(outPath)
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
		err = copyDependencyProtoFile(data, dir, outPath)
		if err != nil {
			return err
		}
		targetProtoFile := targetProtoFolder + "/" + gofile.GetFilename(pf)
		err = copyProtoFile(pf, targetProtoFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyDependencyProtoFile(data []byte, dir string, outPath string) error {
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
		srcProtoFile := dir + "/api/" + pf
		targetProtoFile := outPath + "/" + pf
		pData, err := os.ReadFile(srcProtoFile)
		if err != nil {
			return err
		}
		err = copyProtoFile(srcProtoFile, targetProtoFile)
		if err != nil {
			return err
		}

		if copyCount > 1000 {
			return errors.New("import dependencies circle or too many files")
		}
		err = copyDependencyProtoFile(pData, dir, outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyProtoFile(srcProtoFile string, targetProtoFile string) error {
	targetProtoDir := gofile.GetFileDir(targetProtoFile)
	_ = os.MkdirAll(targetProtoDir, 0766)

	fmt.Printf("copy  \"%s\"  -->  \"%s\"\n", srcProtoFile, targetProtoFile)
	_, err := gobash.Exec("cp", "-f", srcProtoFile, targetProtoFile)
	if err != nil {
		return err
	}
	copyCount++
	return nil
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
