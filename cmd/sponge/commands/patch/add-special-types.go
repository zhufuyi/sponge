package patch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

var (
	protoFileType    = ".proto"
	specialTypesFile = "special_types.go"
)

// AddSpecialTypesCommand add common special types that proto files depend on
func AddSpecialTypesCommand() *cobra.Command {
	var (
		dir        string
		goFileName = "special_types.go"
	)

	cmd := &cobra.Command{
		Use:   "add-special-types",
		Short: "Add common special types that proto files depend on",
		Long: `add common special types that proto files depend on

Examples:
  # add common special types to api directory
  sponge patch add-special-types

  # add common special types to the specified directory
  sponge patch add-special-types --dir=./api/serverName/v1

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dir == "" {
				dir = "api"
			}

			dirFiles, err := getNeedHandleDirs(dir)
			if err != nil {
				return err
			}

			count := 0
			for d, files := range dirFiles {
				isNeed, err := addSpecialTypes(d, files)
				if err != nil {
					return err
				}
				if isNeed {
					fmt.Printf("add %s to %s successful.\n", goFileName, d)
					count++
				}
			}
			if count > 0 {
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "input specified directory")

	return cmd
}

func getNeedHandleDirs(dir string) (map[string][]string, error) {
	dirs, err := gofile.ListDirs(dir)
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		dirs = append(dirs, dir)
	}

	dirs = gofile.FilterDirs(dirs, gofile.WithSuffix(protoFileType))
	if len(dirs) == 0 {
		return nil, nil
	}

	dirFiles := map[string][]string{}
	for _, d := range dirs {
		files, _ := gofile.ListFiles(d, gofile.WithContain(specialTypesFile))
		if len(files) > 0 {
			continue
		}
		files, _ = gofile.ListFiles(d, gofile.WithSuffix(protoFileType))
		if len(files) > 0 {
			dirFiles[d] = files
		}
	}

	return dirFiles, nil
}

func addSpecialTypes(dir string, files []string) (bool, error) {
	isNeedAdd := false
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return false, err
		}
		content := string(data)
		for _, mark := range importProtoMarks {
			if strings.Contains(content, mark) {
				isNeedAdd = true
				break
			}
		}
	}

	if !isNeedAdd {
		return isNeedAdd, nil
	}

	packageName := filepath.Base(dir)
	if packageName != "v1" {
		specialTypesFileData = strings.ReplaceAll(specialTypesFileData, "package v1", fmt.Sprintf("package %s", packageName))
	}
	file := fmt.Sprintf("%s%s%s", dir, gofile.GetPathDelimiter(), specialTypesFile)
	err := os.WriteFile(file, []byte(specialTypesFileData), 0666)

	return isNeedAdd, err
}

var (
	importProtoMarks = []string{
		`"google/protobuf/empty.proto"`,
		`"google/protobuf/any.proto"`,
		`"google/protobuf/duration.proto"`,
		`"google/protobuf/timestamp.proto"`,
		`"google/protobuf/struct.proto"`,
		`"google/protobuf/wrappers.proto"`,
	}

	specialTypesFileData = `package v1

import (
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// import "google/protobuf/empty.proto"
type (
	Empty = emptypb.Empty
)

// import "google/protobuf/any.proto"
type (
	Any = anypb.Any
)

// import "google/protobuf/duration.proto"
type (
	Duration = durationpb.Duration
)

// import "google/protobuf/timestamp.proto"
type (
	Timestamp = timestamppb.Timestamp
)

// import "google/protobuf/struct.proto"
type (
	Struct = structpb.Struct
)

// import "google/protobuf/wrappers.proto"
type (
	BoolValue   = wrapperspb.BoolValue
	StringValue = wrapperspb.StringValue
	Int32Value  = wrapperspb.Int32Value
	Int64Value  = wrapperspb.Int64Value
	UInt32Value = wrapperspb.UInt32Value
	UInt64Value = wrapperspb.UInt64Value
	FloatValue  = wrapperspb.FloatValue
	DoubleValue = wrapperspb.DoubleValue
	BytesValue  = wrapperspb.BytesValue
)

// -----------------------------------------------------------

// If there are more special types, you can add them here
`
)
