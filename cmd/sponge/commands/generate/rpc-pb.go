package generate

import (
	"errors"
	"fmt"

	"github.com/zhufuyi/sponge/pkg/replacer"

	"github.com/huandu/xstrings"
	"github.com/spf13/cobra"
)

// RPCPbCommand generate rpc service codes bash on protobuf file
func RPCPbCommand() *cobra.Command {
	var (
		moduleName   string // module name for go.mod
		serverName   string // server name
		projectName  string // project name for deployment name
		repoAddr     string // image repo address
		outPath      string // output directory
		protobufFile string // protobuf file, support * matching
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "rpc-pb",
		Short: "Generate rpc service codes based on protobuf file",
		Long: `generate rpc service codes based on protobuf file.

Examples:
  # generate rpc service codes.
  sponge micro rpc-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./demo.proto

  # generate rpc service codes and specify the output directory, Note: code generation will be canceled when the latest generated file already exists.
  sponge micro rpc-pb --module-name=yourModuleName --server-name=yourServerName --project-name=yourProjectName --protobuf-file=./demo.proto --out=./yourServerDir

  # generate rpc service codes and specify the docker image repository address.
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

	// setting up template information
	subDirs := []string{ // processing-only subdirectories
		"api/types", "cmd/serverNameExample_grpcPbExample",
		"sponge/build", "sponge/configs", "sponge/deployments", "sponge/scripts", "sponge/third_party",
		"internal/config", "internal/ecode", "internal/server", "internal/service",
	}
	subFiles := []string{ // processing of sub-documents only
		"sponge/.gitignore", "sponge/.golangci.yml", "sponge/go.mod", "sponge/go.sum",
		"sponge/Jenkinsfile", "sponge/Makefile", "sponge/README.md",
	}
	ignoreDirs := []string{} // specify the directory in the subdirectory where processing is ignored
	ignoreFiles := []string{ // specify the files in the subdirectory to be ignored for processing
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

	fmt.Printf("generate %s's rpc server codes successfully, out = %s\n", serverName, r.GetOutputDir())
	fmt.Printf(`help for use:
	1. open a terminal and execute the command to generate the *pb.go file and the service template code: make proto
	2. open file 'internal/service/xxx.go', replace panic("implement me") according to template code example.
	3. compiling and starting services: make run
	4. open the file 'internal/service/xxx_client_test.go' using Goland or VS Code, testing the rpc method.

`)
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
	fields = append(fields, deleteFieldsMark(r, protoShellFile, wellStartMark2, wellEndMark2)...)
	fields = append(fields, deleteFieldsMark(r, appConfigFile, wellStartMark, wellEndMark)...)
	fields = append(fields, replaceFileContentMark(r, readmeFile, "## "+serverName)...)
	fields = append(fields, []replacer.Field{
		{ // replace the contents of the Dockerfile file
			Old: dockerFileMark,
			New: dockerFileHTTPCode,
		},
		{ // replace the contents of the Dockerfile_build file
			Old: dockerFileBuildMark,
			New: dockerFileBuildHTTPCode,
		},
		{ // replace the contents of the docker-compose.yml file
			Old: dockerComposeFileMark,
			New: dockerComposeFileHTTPCode,
		},
		{ // replace the contents of the *-deployment.yml file
			Old: k8sDeploymentFileMark,
			New: k8sDeploymentFileHTTPCode,
		},
		{ // replace the contents of the *-svc.yml file
			Old: k8sServiceFileMark,
			New: k8sServiceFileHTTPCode,
		},
		{ // replace the configuration of the *.yml file
			Old: appConfigFileMark,
			New: rpcServerConfigCode,
		},
		{ // replace the contents of the proto.sh file
			Old: protoShellFileGRPCMark,
			New: protoShellGRPCMark,
		},
		{ // replace the contents of the proto.sh file
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
		// docker image and k8s deployment script replacement
		{
			Old: "server-name-example",
			New: xstrings.ToKebabCase(serverName),
		},
		{
			Old: "projectNameExample",
			New: projectName,
		},
		// docker image and k8s deployment script replacement
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
