package generate

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/sponge/pkg/replacer"

	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"
)

// RPCPbCommand generate rpc server codes bash on protobuf file
func RPCPbCommand() *cobra.Command {
	var (
		moduleName   string // go.mod文件的module名称
		serverName   string // 服务名称
		projectName  string // 项目名称
		repoAddr     string // 镜像仓库地址
		outPath      string // 输出目录
		protobufFile string // proto file文件
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "rpc-pb",
		Short: "Generate rpc server codes based on protobuf file",
		Long: `generate rpc server codes based on protobuf file.

Examples:
  # generate rpc server codes.
  sponge micro rpc-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./demo.proto

  # generate rpc server codes and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge micro rpc-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./demo.proto --out=./yourServerDir

  # generate rpc server codes and specify the docker image repository address.
  sponge micro rpc-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --repo-addr=192.168.3.37:9443/user-name --protobuf-file=./demo.proto
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenRPCPbCommand(moduleName, serverName, projectName, protobufFile, repoAddr, outPath)
		},
	}

	cmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "module-name is the name of the module in the 'go.mod' file")
	_ = cmd.MarkFlagRequired("module-name")
	cmd.Flags().StringVarP(&serverName, "server-name", "s", "", "server name")
	_ = cmd.MarkFlagRequired("server-name")
	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "project name")
	_ = cmd.MarkFlagRequired("project-name")
	cmd.Flags().StringVarP(&protobufFile, "protobuf-file", "f", "", "proto file")
	_ = cmd.MarkFlagRequired("protobuf-file")
	cmd.Flags().StringVarP(&repoAddr, "repo-addr", "r", "", "docker image repository address, excluding http and repository names")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./serverName_rpc-pb_<time>")

	return cmd
}

func runGenRPCPbCommand(moduleName string, serverName string, projectName string, protobufFile string, repoAddr string, outPath string) error {
	protobufFiles, isImportTypes, err := parseProtobufFiles(protobufFile)
	if err != nil {
		return err
	}

	subTplName := "rpc-pb"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{ // 只处理的子目录
		"api/types", "cmd/serverNameExample_grpcPbExample",
		"sponge/build", "sponge/configs", "sponge/deployments", "sponge/scripts", "sponge/third_party",
		"internal/config", "internal/ecode", "internal/server", "internal/service",
	}
	subFiles := []string{ // 只处理子文件
		"sponge/.gitignore", "sponge/.golangci.yml", "sponge/go.mod", "sponge/go.sum",
		"sponge/Jenkinsfile", "sponge/Makefile", "sponge/README.md",
	}
	ignoreDirs := []string{} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{ // 指定子目录下忽略处理的文件
		"types.pb.validate.go", "types.pb.go", // api/types
		"userExample_rpc.go", "systemCode_http.go", "userExample_http.go", // internal/ecode
		"http.go", "http_option.go", "http_test.go", // internal/server
		"userExample.go", "userExample_client_test.go", "userExample_logic.go", "userExample_logic_test.go", "userExample_test.go", // internal/service
	}

	if !isImportTypes {
		ignoreFiles = append(ignoreFiles, "types.proto")
	}

	r.SetSubDirsAndFiles(subDirs, subFiles...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreSubFiles(ignoreFiles...)
	fields := addRPCPbFields(moduleName, serverName, projectName, repoAddr, r)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, serverName+"_"+subTplName)
	if err = r.SaveFiles(); err != nil {
		return err
	}

	_ = saveProtobufFiles(moduleName, serverName, r.GetOutputDir(), protobufFiles)
	_ = saveGenInfo(moduleName, serverName, r.GetOutputDir())

	fmt.Printf("generate %s's rpc gateway server codes successfully, out = %s\n\n", serverName, r.GetOutputDir())

	return nil
}

func addRPCPbFields(moduleName string, serverName string, projectName string, repoAddr string, r replacer.Replacer) []replacer.Field {
	var fields []replacer.Field

	repoHost, _ := parseImageRepoAddr(repoAddr)

	fields = append(fields, deleteFieldsMark(r, dockerFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerFileBuild, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, dockerComposeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sDeploymentFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, k8sServiceFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, imageBuildFile, wellOnlyGrpcStartMark, wellOnlyGrpcEndMark)...)
	fields = append(fields, deleteFieldsMark(r, makeFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, gitIgnoreFile, wellStartMark, wellEndMark)...)
	fields = append(fields, deleteFieldsMark(r, protoShellFile, wellStartMark, wellEndMark)...)
	fields = append(fields, replaceFileContentMark(r, readmeFile, "## "+serverName)...)
	fields = append(fields, []replacer.Field{
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
		{ // 替换proto.sh文件内容
			Old: protoShellFileMark,
			New: protoShellServiceTmplCode,
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
			Old: "_grpcPbExample",
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
	}...)

	return fields
}
