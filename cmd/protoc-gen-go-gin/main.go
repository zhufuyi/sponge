package main

import (
	"flag"
	"fmt"
	"strings"

	"protoc-gen-go-gin/internal/generate/handlerFile"
	"protoc-gen-go-gin/internal/generate/routerFile"
	"protoc-gen-go-gin/internal/generate/serviceFile"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var example bool
	flag.BoolVar(&example, "example", false, "usage example")
	flag.Parse()
	if example {
		fmt.Println(`
# generate *_router.pb.go file
protoc --proto_path=. --proto_path=./third_party --go-gin_out=. --go-gin_opt=paths=source_relative *.proto

# generate *_router.pb.go and *_handler.go files, Note: You need to move *_handler.go to the internal/handler directory
protoc --proto_path=. --proto_path=./third_party --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugins=handler *.proto

# generate *_router.pb.go and *_service.go files, Note: You need to move *_service.go to the internal/service directory
protoc --proto_path=. --proto_path=./third_party --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugins=service *.proto
`)
		return
	}

	var flags flag.FlagSet
	var plugins = flags.String("plugins", "", "list of plugins to enable (supported values: handler,service)")

	options := protogen.Options{
		ParamFunc: flags.Set,
	}

	options.Run(func(gen *protogen.Plugin) error {
		handlerFlag, serviceFlag := false, false
		pluginNames := strings.Split(*plugins, ",")
		for _, plugin := range pluginNames {
			switch strings.ReplaceAll(plugin, " ", "") {
			case "handler":
				handlerFlag = true
			case "service":
				serviceFlag = true
			case "":
			default:
				return fmt.Errorf("protoc-gen-go-gin: unknown plugin %q", plugin)
			}
		}

		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			routerFile.GenerateFile(gen, f)

			if handlerFlag {
				handlerFile.GenerateFile(gen, f)
			}
			if serviceFlag {
				serviceFile.GenerateFile(gen, f)
			}
		}
		return nil
	})
}
