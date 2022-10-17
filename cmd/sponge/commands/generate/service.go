package generate

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/spf13/cobra"
)

// ServiceCommand generate service code
func ServiceCommand() *cobra.Command {
	var (
		moduleName string // go.mod文件的module名称
		serverName string // 服务名称
		outPath    string // 输出目录
		sqlArgs    = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	cmd := &cobra.Command{
		Use:   "service",
		Short: "Generate grpc service code",
		Long: `generate grpc service code.

Examples:
  # generate service code and embed 'gorm.model' struct.
  sponge service --module-name=yourModuleName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate service code, structure fields correspond to the column names of the table.
  sponge service --module-name=yourModuleName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate service code and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge service --module-name=yourModuleName --server-name=yourServerName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			return runGenServiceCommand(moduleName, serverName, codes, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "p", "", "module-name is the name of the module in the 'go.mod' file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./service_<time>")

	return cmd
}

func runGenServiceCommand(moduleName string, serverName string, codes map[string]string, outPath string) error {
	subTplName := "service"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	if serverName == "" {
		serverName = moduleName
	}

	// 设置模板信息
	subDirs := []string{"internal/model", "internal/cache", "internal/dao",
		"internal/ecode", "internal/service", "api/serverNameExample"} // 只处理的指定子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"init.go", "init_test.go", "http_systemCode.go", "http_userExample.go",
		"grpc_systemCode.go", "grpc_systemCode_test.go", "service.go", "service_test.go", "userExample.pb.go",
		"userExample.pb.validate.go", "userExample_grpc.pb.go"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addServiceFields(moduleName, serverName, r, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate '%s' code successfully, out = %s\n\n", subTplName, r.GetOutputDir())
	return nil
}

func addServiceFields(moduleName string, serverName string, r replacer.Replacer, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, protoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, serviceClientFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, serviceTestFile, startMark, endMark)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
		},
		{ // 替换v1/userExample.proto文件内容
			Old: protoFileMark,
			New: codes[parser.CodeTypeProto],
		},
		{ // 替换service/userExample_client_test.go文件内容
			Old: serviceFileMark,
			New: adjustmentOfIDType(codes[parser.CodeTypeService]),
		},
		{
			Old: selfPackageName + "/" + r.GetSourcePath(),
			New: moduleName,
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: moduleName,
		},
		// 替换目录名称
		{
			Old: strings.Join([]string{"api", "serverNameExample", "v1"}, gofile.GetPathDelimiter()),
			New: strings.Join([]string{"api", serverName, "v1"}, gofile.GetPathDelimiter()),
		},
		{
			Old: "api/serverNameExample/v1",
			New: fmt.Sprintf("api/%s/v1", serverName),
		},
		{
			Old: "api.serverNameExample.v1",
			New: fmt.Sprintf("api.%s.v1", strings.ReplaceAll(serverName, "-", "_")), // proto package 不能存在"-"号
		},
		{
			Old: "userExampleNO = 1",
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(1000)),
		},
		{
			Old: moduleName + "/pkg",
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
