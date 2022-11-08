package generate

import (
	"fmt"
	"os"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/jy2struct"

	"github.com/spf13/cobra"
)

// ConfigCommand covert yaml to struct command
func ConfigCommand() *cobra.Command {
	var (
		ysArgs = jy2struct.Args{
			Tags:      "json",
			SubStruct: true,
		}
		serverDir = ""
	)

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Generate go config codes from yaml file",
		Long: `generate go config codes from yaml file.

Examples:
  # generate config codes in server directory, the yaml configuration file must be in <yourServerDir>/configs directory.
  sponge config --server-dir=/yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := getYAMLFile(serverDir)
			if err != nil {
				return err
			}

			err = runGenConfigCommand(files, ysArgs)
			if err != nil {
				return err
			}
			fmt.Println("covert yaml to go struct successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&serverDir, "server-dir", "d", "", "server directory")
	_ = cmd.MarkFlagRequired("server-dir")

	return cmd
}

func runGenConfigCommand(files map[string]configType, ysArgs jy2struct.Args) error {
	for outputFile, config := range files {
		ysArgs.Format = "yaml"
		ysArgs.InputFile = config.configFile

		var startCode string
		if config.isConfigCenter {
			ysArgs.Name = "Center"
			startCode = configFileCcCode
		} else {
			ysArgs.Name = "Config"
			startCode = configFileCode
		}
		structCodes, err := jy2struct.Covert(&ysArgs)
		if err != nil {
			return err
		}
		err = saveFile(config.configFile, outputFile, startCode+structCodes)
		if err != nil {
			return err
		}
	}

	return nil
}

type configType struct {
	configFile     string
	isConfigCenter bool
}

// 从configs目录读取所有yaml文件目录，一个是.yml，另一个是cc.yml
func getYAMLFile(serverDir string) (map[string]configType, error) {
	// 生成目标文件:配置文件
	files := make(map[string]configType)
	configsDir := serverDir + gofile.GetPathDelimiter() + "configs"
	goConfigDir := serverDir + gofile.GetPathDelimiter() + "internal" + gofile.GetPathDelimiter() + "config"

	ymlFiles, err := gofile.ListFiles(configsDir, gofile.WithSuffix(".yml"))
	if err != nil {
		return nil, err
	}

	yamlFiles, err := gofile.ListFiles(configsDir, gofile.WithSuffix(".yaml"))
	if err != nil {
		return nil, err
	}

	if len(ymlFiles) == 0 && len(yamlFiles) == 0 {
		return nil, fmt.Errorf("not found config files in directory %s", configsDir)
	}

	if len(ymlFiles) != 0 && len(yamlFiles) != 0 {
		return nil, fmt.Errorf("please use 'yml' or 'yaml' suffixes for configuration files, do not mix them")
	}

	if len(ymlFiles) > 0 {
		for _, file := range ymlFiles {
			name := gofile.GetFilename(file)
			files[goConfigDir+gofile.GetPathDelimiter()+strings.ReplaceAll(name, ".yml", ".go")] = configType{
				configFile:     file,
				isConfigCenter: strings.Contains(name, "cc.yml"),
			}
		}
		return files, nil
	}

	if len(yamlFiles) > 0 {
		for _, file := range yamlFiles {
			name := gofile.GetFilename(file)
			files[goConfigDir+gofile.GetPathDelimiter()+strings.ReplaceAll(name, ".yaml", ".go")] = configType{
				configFile:     file,
				isConfigCenter: strings.Contains(name, "cc.yaml"),
			}
		}
	}

	return files, nil
}

func saveFile(inputFile string, outputFile string, code string) error {
	err := os.WriteFile(outputFile, []byte(code), 0666)
	if err != nil {
		return err
	}
	fmt.Printf("%s ----> %s\n", inputFile, outputFile)
	return nil
}
