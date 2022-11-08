package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/generate/handler"
	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/generate/router"
	"github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin/internal/generate/service"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const exampleTip = `
# generate *_router.pb.go file
protoc --proto_path=. --proto_path=./third_party --go-gin_out=. --go-gin_opt=paths=source_relative *.proto

# generate *_router.pb.go and handler *_logic.go file
protoc --proto_path=. --proto_path=./third_party --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=handler \
  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/handler *.proto

# generate *_router.pb.go and service *_logic.go
protoc --proto_path=. --proto_path=./third_party --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/service *.proto
`

func main() {
	var example bool
	flag.BoolVar(&example, "example", false, "usage example")
	flag.Parse()
	if example {
		fmt.Printf("%s", exampleTip)
		return
	}

	var flags flag.FlagSet

	var plugin, moduleName, serverName, out string
	flags.StringVar(&plugin, "plugin", "", "list of plugin to enable (supported values: handler or service)")
	flags.StringVar(&moduleName, "moduleName", "", "import module name")
	flags.StringVar(&serverName, "serverName", "", "import server name")
	flags.StringVar(&out, "out", "", "plugin generation code output folder")

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(gen *protogen.Plugin) error {
		handlerFlag, serviceFlag := false, false
		pluginName := strings.ReplaceAll(plugin, " ", "")
		switch pluginName {
		case "handler":
			handlerFlag = true
		case "service":
			serviceFlag = true
		case "":
		default:
			return fmt.Errorf("protoc-gen-go-gin: unknown plugin %q", plugin)
		}

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			router.GenerateFile(gen, f)

			if handlerFlag {
				filename, gf := handler.GenerateFile(gen, f)
				if gf == nil {
					continue
				}
				content, err := gf.Content()
				if err != nil {
					continue // skip error, process the next protobuf file
				}
				err = saveFile(moduleName, serverName, out, filename, content, false)
				if err != nil {
					continue // skip error, process the next protobuf file
				}
				gf.Skip()
			} else if serviceFlag {
				filename, routerContent, gf := service.GenerateFile(gen, f)
				if gf == nil {
					continue
				}
				content, err := gf.Content()
				if err != nil {
					continue // skip error, process the next protobuf file
				}
				err = saveFile(moduleName, serverName, out, filename, content, false)
				if err != nil {
					continue // skip error, process the next protobuf file
				}
				err = saveFile(moduleName, serverName, "internal/routers", filename, routerContent, true)
				if err != nil {
					continue // skip error, process the next protobuf file
				}
				gf.Skip()
			}
		}
		return nil
	})
}

func saveFile(moduleName string, serverName string, out string, filename string, content []byte, isNeedCovered bool) error {
	if moduleName == "" {
		panic("--go-gin_opt option error, 'moduleName' cannot be empty\n" +
			"    usage example: --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/service")
	}
	if serverName == "" {
		panic("--go-gin_opt option error, 'serverName' cannot be empty\n" +
			"    usage example: --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/service")
	}
	if out == "" {
		panic("--go-gin_opt option error, 'out' cannot be empty\n" +
			"    usage example: --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/service")
	}

	_ = os.MkdirAll(out, 0666)
	_, name := filepath.Split(filename)
	if isNeedCovered {
		name = strings.ReplaceAll(name, "_logic.go", ".go")
	}
	file := out + "/" + name
	if !isNeedCovered && isExists(file) {
		file += ".gen." + time.Now().Format("150405")
	}

	content = bytes.ReplaceAll(content, []byte("moduleNameExample"), []byte(moduleName))
	content = bytes.ReplaceAll(content, []byte("serverNameExample"), []byte(serverName))
	return os.WriteFile(file, content, 0666)
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
