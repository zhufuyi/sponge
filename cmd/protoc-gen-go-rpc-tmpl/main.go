// Package main is to generate *.go(tmpl), *_client_test.go, *_rpc.go files.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl/internal/generate/service"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	helpInfo = `
# generate *.go file
protoc --proto_path=. --proto_path=./third_party --go-rpc-tmpl_out=. --go-rpc-tmpl_opt=paths=source_relative \
  --go-rpc-tmpl_opt=moduleName=yourModuleName --go-rpc-tmpl_opt=serverName=yourServerName *.proto

Note:
    If you want to merge the code, after generating the code, execute the command "sponge merge rpc-pb",
    you don't worry about it affecting the logic code you have already written, in case of accidents,
    you can find the pre-merge code in the directory /tmp/sponge_merge_backup_code.
`

	optErrFormat = `--go-rpc-tmpl_opt error, '%s' cannot be empty.

Usage example: 
    protoc --proto_path=. --proto_path=./third_party \
      --go-rpc-tmpl_out=. --go-rpc-tmpl_opt=paths=source_relative \
      --go-rpc-tmpl_opt=moduleName=yourModuleName --go-rpc-tmpl_opt=serverName=yourServerName \
      *.proto
`
)

func main() {
	var h bool
	flag.BoolVar(&h, "h", false, "help information")
	flag.Parse()
	if h {
		fmt.Printf("%s", helpInfo)
		return
	}

	var flags flag.FlagSet

	var moduleName, serverName, tmplDir, ecodeOut string
	flags.StringVar(&moduleName, "moduleName", "", "module name")
	flags.StringVar(&serverName, "serverName", "", "server name")
	flags.StringVar(&tmplDir, "tmplDir", "internal/service", "rpc template file directory, the default value is internal/service")
	flags.StringVar(&ecodeOut, "ecodeOut", "internal/ecode", "rpc error code file directory, the default value is internal/ecode")

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			err := saveRPCTmplFiles(f, moduleName, serverName, tmplDir, ecodeOut)
			if err != nil {
				continue // skip error, process the next protobuf file
			}
		}
		return nil
	})
}

func saveRPCTmplFiles(f *protogen.File, moduleName string, serverName string, tmplOut string, ecodeOut string) error {
	filenamePrefix := f.GeneratedFilenamePrefix
	tmplFileContent, testTmplFileContent, ecodeFileContent := service.GenerateFiles(filenamePrefix, f)

	filePath := filenamePrefix + ".go"
	err := saveFile(moduleName, serverName, tmplOut, filePath, tmplFileContent, false)
	if err != nil {
		return err
	}

	filePath = filenamePrefix + "_client_test.go"
	err = saveFile(moduleName, serverName, tmplOut, filePath, testTmplFileContent, true)
	if err != nil {
		return err
	}

	filePath = filenamePrefix + "_rpc.go"
	err = saveFileSimple(ecodeOut, filePath, ecodeFileContent, false)
	if err != nil {
		return err
	}

	return nil
}

func saveFile(moduleName string, serverName string, out string, filePath string, content []byte, isNeedCovered bool) error {
	if len(content) == 0 {
		return nil
	}

	if moduleName == "" {
		panic(fmt.Sprintf(optErrFormat, "moduleName"))
	}
	if serverName == "" {
		panic(fmt.Sprintf(optErrFormat, "serverName"))
	}

	_ = os.MkdirAll(out, 0766)
	_, name := filepath.Split(filePath)
	file := out + "/" + name
	if !isNeedCovered && isExists(file) {
		file += ".gen" + time.Now().Format("20060102T150405")
	}

	content = bytes.ReplaceAll(content, []byte("moduleNameExample"), []byte(moduleName))
	content = bytes.ReplaceAll(content, []byte("serverNameExample"), []byte(serverName))
	content = bytes.ReplaceAll(content, firstLetterToUpper("serverNameExample"), firstLetterToUpper(serverName))
	return os.WriteFile(file, content, 0666)
}

func saveFileSimple(out string, filePath string, content []byte, isNeedCovered bool) error {
	if len(content) == 0 {
		return nil
	}

	_ = os.MkdirAll(out, 0766)
	_, name := filepath.Split(filePath)
	file := out + "/" + name
	if !isNeedCovered && isExists(file) {
		file += ".gen" + time.Now().Format("20060102T150405")
	}

	return os.WriteFile(file, content, 0666)
}

func isExists(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func firstLetterToUpper(s string) []byte {
	if s == "" {
		return []byte{}
	}

	return []byte(strings.ToUpper(s[:1]) + s[1:])
}
