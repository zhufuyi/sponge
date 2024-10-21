package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/replacer"
)

// GRPCConnectionCommand generate grpc connection code
func GRPCConnectionCommand() *cobra.Command {
	var (
		moduleName      string // module name for go.mod
		outPath         string // output directory
		grpcServerNames string // grpc service names

		serverName     string // server name
		suitedMonoRepo bool   // whether the generated code is suitable for mono-repo
	)

	cmd := &cobra.Command{
		Use:   "rpc-conn",
		Short: "Generate grpc connection code",
		Long:  "Generate grpc connection code.",
		Example: color.HiBlackString(`  # Generate grpc connection code
  sponge micro rpc-conn --module-name=yourModuleName --rpc-server-name=yourGrpcName

  # Generate grpc connection code with multiple names.
  sponge micro rpc-conn --module-name=yourModuleName --rpc-server-name=name1,name2

  # Generate grpc connection code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge micro rpc-conn --rpc-server-name=user --out=./yourServerDir

  # If you want the generated code to suited to mono-repo, you need to set the parameter --suited-mono-repo=true --server-name=yourServerName`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mdName, srvName, smr := getNamesFromOutDir(outPath)
			if mdName != "" {
				moduleName = mdName
				serverName = srvName
				suitedMonoRepo = smr
			} else if moduleName == "" {
				return errors.New(`required flag(s) "module-name" not set, use "sponge micro rpc-conn -h" for help`)
			}
			if suitedMonoRepo {
				if serverName == "" {
					return errors.New(`required flag(s) "server-name" not set, use "sponge micro rpc-conn -h" for help`)
				}
				serverName = convertServerName(serverName)
				outPath = changeOutPath(outPath, serverName)
			}

			grpcNames := strings.Split(grpcServerNames, ",")
			for _, grpcName := range grpcNames {
				if grpcName == "" {
					continue
				}

				var err error
				var g = &grpcConnectionGenerator{
					moduleName: moduleName,
					grpcName:   grpcName,
					outPath:    outPath,

					serverName:     serverName,
					suitedMonoRepo: suitedMonoRepo,
				}
				outPath, err = g.generateCode()
				if err != nil {
					return err
				}
			}

			fmt.Printf(`
using help:
  move the folder "internal" to your project code folder.

`)
			fmt.Printf("generate \"rpc-conn\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	cmd.Flags().StringVarP(&grpcServerNames, "rpc-server-name", "r", "", "rpc service name, multiple names separated by commas")
	_ = cmd.MarkFlagRequired("rpc-server-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	cmd.Flags().BoolVarP(&suitedMonoRepo, "suited-mono-repo", "l", false, "whether the generated code is suitable for mono-repo")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./rpc-conn_<time>, "+flagTip("module-name"))

	return cmd
}

type grpcConnectionGenerator struct {
	moduleName string
	grpcName   string
	outPath    string

	serverName     string
	suitedMonoRepo bool
}

func (g *grpcConnectionGenerator) generateCode() (string, error) {
	subTplName := codeNameGRPCConn
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// specify the subdirectory and files
	subDirs := []string{}
	subFiles := []string{"internal/rpcclient/serverNameExample.go"}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	_ = r.SetOutputDir(g.outPath, subTplName)
	fields := g.addFields()
	r.SetReplacementFields(fields)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func (g *grpcConnectionGenerator) addFields() []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, []replacer.Field{
		{
			Old: "github.com/zhufuyi/sponge/configs",
			New: g.moduleName + "/configs",
		},
		{
			Old: "github.com/zhufuyi/sponge/internal/config",
			New: g.moduleName + "/internal/config",
		},
		{
			Old:             "serverNameExample",
			New:             g.grpcName,
			IsCaseSensitive: true,
		},
	}...)

	if g.suitedMonoRepo {
		fs := SubServerCodeFields(g.moduleName, g.serverName)
		fields = append(fields, fs...)
	}

	return fields
}
