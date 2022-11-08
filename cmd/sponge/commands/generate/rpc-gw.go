package generate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/zhufuyi/sponge/pkg/replacer"

	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"
)

// RPCGwCommand generate rpc gateway server codes
func RPCGwCommand() *cobra.Command {
	var (
		moduleName   string // go.mod文件的module名称
		serverName   string // 服务名称
		projectName  string // 项目名称
		repoAddr     string // 镜像仓库地址
		outPath      string // 输出目录
		protobufFile string // proto file文件，指定这个文件生成路由和service
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "rpc-gw",
		Short: "Generate rpc gateway server codes based on protobuf",
		Long: `generate rpc gateway server codes based on protobuf.

Examples:
  # generate rpc gateway server codes.
  sponge micro rpc-gw --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./userExample.proto

  # generate rpc gateway server codes and specify the output directory, Note: if the file already exists, code generation will be canceled.
  sponge micro rpc-gw --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./userExample.proto --out=./yourServerDir

  # generate rpc gateway server codes and specify the docker image repository address.
  sponge micro rpc-gw --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --repo-addr=192.168.3.37:9443/user-name --protobuf-file=./userExample.proto
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenRPCGwCommand(moduleName, serverName, projectName, protobufFile, repoAddr, outPath)
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
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output directory, default is ./serverName_rpc_<time>")

	return cmd
}

func runGenRPCGwCommand(moduleName string, serverName string, projectName string, protobufFile string, repoAddr string, outPath string) error {
	protoData, err := os.ReadFile(protobufFile)
	if err != nil {
		return err
	}

	err = getServiceName(protoData)
	if err != nil {
		return err
	}

	subTplName := "rpc-gw"
	r := Replacers[TplNameSponge]
	if r == nil {
		return errors.New("replacer is nil")
	}

	// 设置模板信息
	subDirs := []string{} // 只处理的子目录，如果为空或者没有指定的子目录，表示所有文件
	ignoreDirs := []string{"cmd/sponge", "cmd/protoc-gen-go-gin", "cmd/serverNameExample_mixExample",
		"cmd/serverNameExample_grpcExample", "cmd/serverNameExample_httpExample",
		"sponge/.github", "sponge/.git", "sponge/pkg", "sponge/assets", "sponge/test",
		"internal/model", "internal/cache", "internal/dao", "internal/ecode", "internal/service",
		"internal/handler", "internal/types", "api/serverNameExample",
	} // 指定子目录下忽略处理的目录
	ignoreFiles := []string{"grpc.go", "grpc_option.go", "grpc_test.go", "LICENSE",
		"grpc_userExample.go", "grpc_systemCode.go", "grpc_systemCode_test.go",
		"codecov.yml", "routers.go", "routers_test.go", "userExample_gwExample.go", "userExample.go",
		"routers_gwExample_test.go", "userExample_gwExample.go", "types.pb.validate.go", "types.pb.go",
		"swagger.json", "swagger.yaml", "apis.swagger.json", "proto.html", "docs.go", "doc.go",
	} // 指定子目录下忽略处理的文件

	if !bytes.Contains(protoData, []byte("api/types/types.proto")) {
		ignoreFiles = append(ignoreFiles, "types.proto")
	}

	r.SetSubDirs(subDirs...)
	r.SetIgnoreSubDirs(ignoreDirs...)
	r.SetIgnoreFiles(ignoreFiles...)
	fields := addRPCGwFields(moduleName, serverName, projectName, repoAddr, r)
	r.SetReplacementFields(fields)
	_ = r.SetOutputDir(outPath, serverName+"_"+subTplName)
	if err = r.SaveFiles(); err != nil {
		return err
	}

	fmt.Printf("generate %s's rpc gateway server codes successfully, out = %s\n\n", serverName, r.GetOutputDir())

	// 保存moduleName和serverName到指定文件，给外部使用
	genInfo := moduleName + "," + serverName
	file := r.GetOutputDir() + "/docs/gen.info"
	err = os.WriteFile(file, []byte(genInfo), 0666)
	if err != nil {
		fmt.Printf("save file %s error, %v\n", file, err)
	}

	// 复制protobuf文件
	_, name := filepath.Split(protobufFile)
	dir := r.GetOutputDir() + "/api/" + serverName + "/v1"
	_ = os.MkdirAll(dir, 0666)
	file = dir + "/" + name
	err = os.WriteFile(file, protoData, 0666)
	if err != nil {
		fmt.Printf("save file %s error, %v\n", file, err)
	}

	return nil
}

func addRPCGwFields(moduleName string, serverName string, projectName string, repoAddr string,
	r replacer.Replacer) []replacer.Field {
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
	fields = append(fields, deleteFieldsMark(r, protoShellFile, wellStartMark, wellEndMark)...)
	fields = append(fields, replaceFileContentMark(r, readmeFile, "## "+serverName)...)
	//fields = append(fields, replaceFileContentMark(r, protoFile, string(protoData))...)
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
			New: protoShellServiceCode,
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
			Old: "_httpExample",
			New: "",
		},
		{
			Old: "_mixExample",
			New: "",
		},
		{
			Old: "_gwExample",
			New: "",
		},
	}...)

	return fields
}

func getServiceName(data []byte) error {
	servicePattern := `\nservice (\w+)`
	re := regexp.MustCompile(servicePattern)
	matchArr := re.FindStringSubmatch(string(data))
	if len(matchArr) < 2 {
		return fmt.Errorf("not found service name in protobuf file, the protobuf file requires at least one service")
	}
	return nil
}
