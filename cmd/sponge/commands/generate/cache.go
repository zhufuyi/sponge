package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/replacer"
)

// CacheCommand generate cache code
func CacheCommand(parentName string) *cobra.Command {
	var (
		moduleName string // module name for go.mod
		outPath    string // output directory
		cacheName  string // cache name
		prefixKey  string // prefix key
		keyName    string // key name
		keyType    string // key type
		valueName  string // value name
		valueType  string // value type

		serverName     string // server name
		suitedMonoRepo bool   // whether the generated code is suitable for mono-repo
	)

	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Generate cache code",
		Long:  "Generate cache code.",
		Example: color.HiBlackString(fmt.Sprintf(`  # Generate kv cache code
  sponge %s cache --module-name=yourModuleName --cache-name=userToken --key-name=id --key-type=uint64 --value-name=token --value-type=string

  # Generate kv cache code and specify the server directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge %s cache --module-name=yourModuleName --cache-name=token --prefix-key=user:token --key-name=id --key-type=uint64 --value-name=token --value-type=string --out=./yourServerDir

  # If you want the generated code to suited to mono-repo, you need to set the parameter --suited-mono-repo=true --server-name=yourServerName`,
			parentName, parentName)),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mdName, srvName, smr := getNamesFromOutDir(outPath)
			if mdName != "" {
				moduleName = mdName
				serverName = srvName
				suitedMonoRepo = smr
			} else if moduleName == "" {
				return errors.New(`required flag(s) "module-name" not set, use "sponge micro cache -h" for help`)
			}
			if suitedMonoRepo {
				if serverName == "" {
					return fmt.Errorf(`required flag(s) "server-name" not set, use "sponge %s cache -h" for help`, parentName)
				}
				serverName = convertServerName(serverName)
				outPath = changeOutPath(outPath, serverName)
			}
			cacheName = strings.ReplaceAll(cacheName, ":", "")

			if prefixKey == "" || prefixKey == ":" {
				prefixKey = cacheName + ":"
			} else if prefixKey[len(prefixKey)-1] != ':' {
				prefixKey += ":"
			}

			var err error
			var g = &stringCacheGenerator{
				moduleName: moduleName,
				cacheName:  cacheName,
				prefixKey:  prefixKey,
				keyName:    keyName,
				keyType:    keyType,
				valueName:  valueName,
				valueType:  valueType,
				outPath:    outPath,

				serverName:     serverName,
				suitedMonoRepo: suitedMonoRepo,
			}
			outPath, err = g.generateCode()
			if err != nil {
				return err
			}

			fmt.Printf(`
using help:
  move the folder "internal" to your project code folder.

`)
			fmt.Printf("generate \"cache\" code successfully, out = %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the go.mod file")
	cmd.Flags().StringVarP(&cacheName, "cache-name", "c", "", "cache name, e.g. userToken")
	_ = cmd.MarkFlagRequired("cache-name")
	cmd.Flags().StringVarP(&prefixKey, "prefix-key", "p", "", "cache prefix key, e.g. user:token")
	cmd.Flags().StringVarP(&keyName, "key-name", "k", "", "key name, e.g. id")
	_ = cmd.MarkFlagRequired("key-name")
	cmd.Flags().StringVarP(&keyType, "key-type", "t", "", "key go type, e.g. uint64")
	_ = cmd.MarkFlagRequired("key-type")
	cmd.Flags().StringVarP(&valueName, "value-name", "v", "", "value name, e.g. token")
	_ = cmd.MarkFlagRequired("value-name")
	cmd.Flags().StringVarP(&valueType, "value-type", "w", "", "value go type, e.g. string")
	_ = cmd.MarkFlagRequired("value-type")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	cmd.Flags().BoolVarP(&suitedMonoRepo, "suited-mono-repo", "l", false, "whether the generated code is suitable for mono-repo")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./cache_<time>, "+flagTip("module-name"))

	return cmd
}

type stringCacheGenerator struct {
	moduleName string
	cacheName  string
	prefixKey  string
	keyName    string
	keyType    string
	valueName  string
	valueType  string
	outPath    string

	serverName     string
	suitedMonoRepo bool
}

func (g *stringCacheGenerator) generateCode() (string, error) {
	subTplName := codeNameCache
	r := Replacers[TplNameSponge]
	if r == nil {
		return "", errors.New("replacer is nil")
	}

	// specify the subdirectory and files
	subDirs := []string{}
	subFiles := []string{"internal/cache/cacheNameExample.go"}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	_ = r.SetOutputDir(g.outPath, subTplName)
	fields := g.addFields(r)
	r.SetReplacementFields(fields)
	if err := r.SaveFiles(); err != nil {
		return "", err
	}

	return r.GetOutputDir(), nil
}

func (g *stringCacheGenerator) addFields(r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field
	fields = append(fields, deleteFieldsMark(r, cacheFile, startMark, endMark)...)

	// match the case where the value type is a pointer
	if g.valueType[0] == '*' {
		fields = append(fields, []replacer.Field{
			{
				Old:             "var valueNameExample valueTypeExample",
				New:             fmt.Sprintf("%s := &%s{}", g.valueName, g.valueType[1:]),
				IsCaseSensitive: false,
			},
			{
				Old:             "&valueNameExample",
				New:             g.valueName,
				IsCaseSensitive: false,
			},
		}...)
	}

	fields = append(fields, []replacer.Field{
		{
			Old: "github.com/go-dev-frame/sponge/internal/model",
			New: g.moduleName + "/internal/model",
		},
		{
			Old:             "cacheNameExample",
			New:             g.cacheName,
			IsCaseSensitive: true,
		},
		{
			Old:             "prefixKeyExample:",
			New:             g.prefixKey,
			IsCaseSensitive: false,
		},
		{
			Old:             "keyNameExample",
			New:             g.keyName,
			IsCaseSensitive: false,
		},
		{
			Old:             "keyTypeExample",
			New:             g.keyType,
			IsCaseSensitive: false,
		},
		{
			Old:             "valueNameExample",
			New:             g.valueName,
			IsCaseSensitive: false,
		},
		{
			Old:             "valueTypeExample",
			New:             g.valueType,
			IsCaseSensitive: false,
		},
	}...)

	if g.suitedMonoRepo {
		fs := SubServerCodeFields(g.moduleName, g.serverName)
		fields = append(fields, fs...)
	}

	return fields
}
