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

// HTTPCommand generate http code
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
		Short: "Generate http server code",
		Long: `generate http server code.

Examples:
  # generate http code and embed 'gorm.model' struct.
  sponge http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user

  # generate http code, structure fields correspond to the column names of the table.
  sponge http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --embed=false

  # generate http code and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user --out=./yourServerDir

  # generate http code and specify the docker image repository address.
  sponge http --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --repo-addr=192.168.3.37:9443/user-name --db-dsn=root:123456@(192.168.3.37:3306)/test --db-table=user
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
	subDirs := []string{} // 只处理的子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{"cmd/sponge", "sponge/.github", "sponge/.git", "sponge/api", "sponge/pkg",
		"sponge/assets", "sponge/test", "sponge/third_party", "internal/service"} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"swagger.json", "swagger.yaml", "proto.html", "protoc.sh",
		"proto-doc.sh", "grpc_health_probe.sh", "grpc.go", "grpc_option.go", "grpc_test.go", "LICENSE", "doc.go",
		"grpc_userExample.go", "grpc_systemCode.go", "grpc_systemCode_test.go", "codecov.yml"} // 指定子目录下忽略处理的文件

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addHTTPFields(moduleName, serverName, projectName, repoAddr, r, dbDSN, codes)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, serverName+"_"+subTplName)
	if err := r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate %s's http code successfully, out = %s\n\n", serverName, r.GetOutputDir())
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
	fields = append(fields, deleteFieldsMark(r, mainFile, startMark, endMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFileBuild, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerComposeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sDeploymentFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sServiceFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, imageBuildFile, wellOnlyGrpcStartMark, wellOnlyGrpcEndMark)...)
	fields = append(fields, deleteFieldsMark(r, makeFile, wellStartMark, wellEndMark)...)
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
		{ // 替换main.go文件内容
			Old: mainFileMark,
			New: mainFileHTTPCode,
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
			New: fmt.Sprintf("userExampleNO = %d", rand.Intn(1000)),
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
			Old: "tmp.go.mod",
			New: "go.mod",
		},
		{
			Old: "tmp.gitignore",
			New: ".gitignore",
		},
		{
			Old: "tmp.golangci.yml",
			New: ".golangci.yml",
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
