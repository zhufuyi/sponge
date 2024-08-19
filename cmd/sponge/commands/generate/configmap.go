package generate

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/replacer"
)

// ConfigmapCommand generate k8s configmap command
func ConfigmapCommand() *cobra.Command {
	var (
		serverName  = ""
		projectName = ""
		configFile  = ""
		outPath     = ""
	)

	cmd := &cobra.Command{
		Use:   "configmap",
		Short: "Generate k8s configmap",
		Long: color.HiBlackString(`generate k8s configmap.

Examples:
  # generate k8s configmap
  sponge configmap --server-name=yourServerName --project-name=yourProjectName --config-file=yourConfigFile.yml

  # generate grpc connection code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge configmap --server-name=yourServerName --project-name=yourProjectName --config-file=yourConfigFile.yml --out=./yourServerDir
`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := convertYamlConfig(configFile)
			if err != nil {
				return err
			}
			g := copyConfigGenerator{
				serverName:  serverName,
				projectName: projectName,
				content:     content,
				outPath:     outPath,
			}
			outPath, err = g.generateCode()
			if err != nil {
				return err
			}
			fmt.Printf("\ngenerate \"configmap\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&configFile, "config-file", "f", "", "server config file")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./configmap_<time>")

	return cmd
}

type copyConfigGenerator struct {
	serverName  string
	projectName string
	content     string
	outPath     string
}

func (g *copyConfigGenerator) generateCode() (string, error) {
	subTplName := "configmap"
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// setting up template information
	subDirs := []string{ // only the specified subdirectory is processed, if empty or no subdirectory is specified, it means all files
		"deployments/kubernetes",
	}
	ignoreDirs := []string{} // specify the directory in the subdirectory where processing is ignored
	ignoreFiles := []string{ // specify the files in the subdirectory to be ignored for processing
		"projectNameExample-namespace.yml", "README.md", "serverNameExample-deployment.yml", "serverNameExample-svc.yml",
	}

	r.SetSubDirsAndFiles(subDirs)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := g.addFields()
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(g.outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func (g *copyConfigGenerator) addFields() []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, []replacer.Field{
		{
			Old: configmapFileMark,
			New: g.content,
		},
		{
			Old:             "serverNameExample",
			New:             g.serverName,
			IsCaseSensitive: true,
		},
		{
			Old: "server-name-example",
			New: xstrings.ToKebabCase(g.serverName), // snake_case to kebab_case
		},
		{
			Old: "project-name-example",
			New: g.projectName,
		},
	}...)

	return fields
}
