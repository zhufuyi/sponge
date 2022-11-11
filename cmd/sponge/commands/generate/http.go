package generate

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/zhufuyi/sponge/pkg/replacer"
	"github.com/zhufuyi/sponge/pkg/sql2code"
	"github.com/zhufuyi/sponge/pkg/sql2code/parser"

	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"
)

// HTTPCommand generate http server codes
func HTTPCommand() *cobra.Command {
	var (
		moduleName  string // go.mod文件的module名称
		serverName  string // 服务名称
		projectName string // 项目名称
		repoAddr    string // 镜像仓库地址
		outPath     string // 输出目录
		sqlArgs     = sql2code.Args{
			Package:  "model",
			JSONTag:  true,
			GormType: true,
		}
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Generate http server codes based on mysql",
		Long: `generate http server codes based on mysql.

Examples:
  # generate http server codes and embed 'gorm.model' struct.
  sponge web http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate http server codes, structure fields correspond to the column names of the table.
  sponge web http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate http server codes and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge web http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir

  # generate http server codes and specify the docker image repository address.
  sponge web http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --repo-addr=192.168.3.37:9443/user-name --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := sql2code.Generate(&sqlArgs)
			if err != nil {
				return err
			}

			return runGenHTTPCommand(moduleName, serverName, projectName, repoAddr, sqlArgs.DBDsn, codes, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the 'go.mod' file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&sqlArgs.DBDsn, "db-dsn", "d", "", "db content addr, e.g. user:password@(host:port)/database")
	_ = cmd.MarkFlagRequired("db-dsn")
	cmd.Flags().StringVarP(&sqlArgs.DBTable, "db-table", "t", "", "table name")
	_ = cmd.MarkFlagRequired("db-table")
	cmd.Flags().BoolVarP(&sqlArgs.IsEmbed, "embed", "e", true, "whether to embed 'gorm.Model' struct")
	cmd.Flags().StringVarP(&repoAddr, "repo-addr", "r", "", "docker image repository address, excluding http and repository names")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./serverName_http_<time>")

	return cmd
}

func runGenHTTPCommand(moduleName string, serverName string, projectName string, repoAddr string,
	dbDSN string, codes map[string]string, outPath string) error {
	subTplName := "http"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{ // 指定处理的子目录
		"cmd/serverNameExample_httpExample", "sponge/build", "sponge/configs", "sponge/deployments",
		"sponge/docs", "sponge/scripts", "sponge/internal",
	}
	subFiles := []string{ // 指定处理的子文件
		"sponge/.gitignore", "sponge/.golangci.yml", "sponge/go.mod", "sponge/go.sum",
		"sponge/Jenkinsfile", "sponge/Makefile", "sponge/README.md",
	}
	ignoreDirs := []string{ // 指定子目录下忽略处理的目录
		"internal/service", "internal/rpcclient",
	}
	ignoreFiles := []string{ // 指定子目录下忽略处理的文件
		"swagger.json", "swagger.yaml", "apis.swagger.json", "apis.html", "apis.go", // sponge/docs
		"userExample_rpc.go", "systemCode_rpc.go", // internal/ecode
		"routers_pbExample.go", "routers_pbExample_test.go", "userExample_service.pb.go", // internal/routers
		"grpc.go", "grpc_option.go", "grpc_test.go", // internal/server
	}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addHTTPFields(moduleName, serverName, projectName, repoAddr, r, dbDSN, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, serverName+"_"+subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	_ = saveGenInfo(moduleName, serverName, r.GetOutputDir())

	fmt.Printf("generate %s's http server codes successfully, out = %s\n\n", serverName, r.GetOutputDir())
	return nil
}

func addHTTPFields(moduleName string, serverName string, projectName string, repoAddr string,
	r replacer.Replacer, dbDSN string, codes map[string]string) []replacer.Field {
	var fields []replacer.Field

	repoHost, _ := parseImageRepoAddr(repoAddr)

	fields = append(fields, deleteFieldsMark(r, modelFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, daoTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, handlerFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, handlerTestFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, httpFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFileBuild, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerComposeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sDeploymentFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sServiceFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, imageBuildFile, wellOnlyGrpcStartMark, wellOnlyGrpcEndMark)...)
	fields = append(fields, deleteFieldsMark(r, makeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, gitIgnoreFile, wellStartMark, wellEndMark)...)
	fields = append(fields, replaceFileContentMark(r, readmeFile, "## "+serverName)...)
	fields = append(fields, []replacer.Field{
		{ // 替换model/userExample.go文件内容
			Old: modelFileMark,
			New: codes[parser.CodeTypeModel],
		},
		{ // 替换dao/userExample.go文件内容
			Old: daoFileMark,
			New: codes[parser.CodeTypeDAO],
		},
		{ // 替换handler/userExample.go文件内容
			Old: handlerFileMark,
			New: adjustmentOfIDType(codes[parser.CodeTypeHandler]),
		},
		{ // 替换Dockerfile文件内容
			Old: dockerFileMark,
			New: dockerFileHTTPCode,
		},
		{ // 替换Dockerfile_build文件内容
			Old: dockerFileBuildMark,
			New: dockerFileBuildHTTPCode,
		},
		{ // 替换docker-compose.yml文件内容
			Old: dockerComposeFileMark,
			New: dockerComposeFileHTTPCode,
		},
		{ // 替换*-deployment.yml文件内容
			Old: k8sDeploymentFileMark,
			New: k8sDeploymentFileHTTPCode,
		},
		{ // 替换*-svc.yml文件内容
			Old: k8sServiceFileMark,
			New: k8sServiceFileHTTPCode,
		},
		// 替换github.com/zhufuyi/sponge/templates/sponge
		{
			Old: selfPackageName + "/" + r.GetSourcePath(),
			New: moduleName,
		},
		{
			Old: "github.com/zhufuyi/sponge",
			New: moduleName,
		},
		{
			Old: moduleName + "/pkg",
			New: "github.com/zhufuyi/sponge/pkg",
		},
		{
			Old: "sponge api docs",
			New: serverName + " api docs",
		},
		{
			Old: "userExampleNO = 1",
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(100)),
		},
		{
			Old: "serverNameExample",
			New: serverName,
		},
		// docker镜像和k8s部署脚本替换
		{
			Old: "server-name-example",
			New: xstrings.ToKebabCase(serverName),
		},
		{
			Old: "projectNameExample",
			New: projectName,
		},
		// docker镜像和k8s部署脚本替换
		{
			Old: "project-name-example",
			New: xstrings.ToKebabCase(projectName),
		},
		{
			Old: "repo-addr-example",
			New: repoAddr,
		},
		{
			Old: "image-repo-host",
			New: repoHost,
		},
		{
			Old: "_httpExample",
			New: "",
		},
		{
			Old: "_mixExample",
			New: "",
		},
		{
			Old: "_pbExample",
			New: "",
		},
		{
			Old: "root:123456@(192.168.3.37:3306)/account",
			New: dbDSN,
		},
		{
			Old:             "UserExample",
			New:             codes[parser.TableName],
			IsCaseSensitive: true,
		},
	}...)

	return fields
}
