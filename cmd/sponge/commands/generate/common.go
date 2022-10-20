package generate

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/replacer"
)

const (
	// TplNameSponge 模板目录名称
	TplNameSponge = "sponge"
)

var (
	// 指定文件替换标记
	modelFile     = "model/userExample.go"
	modelFileMark = "// todo generate model codes to here"

	daoFile     = "dao/userExample.go"
	daoFileMark = "// todo generate the update fields code to here"
	daoTestFile = "dao/userExample_test.go"

	handlerFile     = "types/userExample_types.go"
	handlerFileMark = "// todo generate the request and response struct to here"
	handlerTestFile = "handler/userExample_test.go"

	mainFile     = "serverNameExample/main.go"
	mainFileMark = "// todo generate the code to register http and grpc services here"

	protoFile     = "v1/userExample.proto"
	protoFileMark = "// todo generate the protobuf code here"

	serviceTestFile   = "service/userExample_test.go"
	serviceClientFile = "service/userExample_client_test.go"
	serviceFileMark   = "// todo generate the service struct code here"

	dockerFile     = "build/Dockerfile"
	dockerFileMark = "# todo generate dockerfile code for http or grpc here"

	dockerFileBuild     = "build/Dockerfile_build"
	dockerFileBuildMark = "# todo generate dockerfile_build code for http or grpc here"

	dockerComposeFile     = "deployments/docker-compose/docker-compose.yml"
	dockerComposeFileMark = "# todo generate docker-compose.yml code for http or grpc here"

	k8sDeploymentFile     = "deployments/kubernetes/serverNameExample-deployment.yml"
	k8sDeploymentFileMark = "# todo generate k8s-deployment.yml code for http or grpc here"

	k8sServiceFile     = "deployments/kubernetes/serverNameExample-svc.yml"
	k8sServiceFileMark = "# todo generate k8s-svc.yml code for http or grpc here"

	imageBuildFile = "scripts/image-build.sh"
	readmeFile     = "sponge/README.md"

	makeFile = "sponge/Makefile"

	// 清除标记的模板代码片段标记
	startMark             = []byte("// delete the templates code start")
	endMark               = []byte("// delete the templates code end")
	wellStartMark         = bytes.ReplaceAll(startMark, []byte("//"), []byte("#"))
	wellEndMark           = bytes.ReplaceAll(endMark, []byte("//"), []byte("#"))
	onlyGrpcStartMark     = []byte("// only grpc use start")
	onlyGrpcEndMark       = []byte("// only grpc use end\n")
	wellOnlyGrpcStartMark = bytes.ReplaceAll(onlyGrpcStartMark, []byte("//"), []byte("#"))
	wellOnlyGrpcEndMark   = bytes.ReplaceAll(onlyGrpcEndMark, []byte("//"), []byte("#"))

	// embed FS模板文件时使用
	selfPackageName = "github.com/zhufuyi/sponge"
)

func adjustmentOfIDType(handlerCodes string) string {
	return idTypeToStr(idTypeFixToUint64(handlerCodes))
}

func idTypeFixToUint64(handlerCodes string) string {
	subStart := "ByIDRequest struct {"
	subEnd := "`" + `json:"id" binding:""` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID uint64 " + subEnd + " // uint64 id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func idTypeToStr(handlerCodes string) string {
	subStart := "ByIDRespond struct {"
	subEnd := "`" + `json:"id"` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID string " + subEnd + " // covert to string id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func deleteFieldsMark(r replacer.Replacer, filename string, startMark []byte, endMark []byte) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		fmt.Printf("read the file '%s' error: %v\n", filename, err)
		return fields
	}
	if subBytes := gofile.FindSubBytes(data, startMark, endMark); len(subBytes) > 0 {
		fields = append(fields,
			replacer.Field{ // 清除标记的模板代码
				Old: string(subBytes),
				New: "",
			},
		)
	}

	return fields
}

func replaceFileContentMark(r replacer.Replacer, filename string, newContent string) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		fmt.Printf("read the file '%s' error: %v\n", filename, err)
		return fields
	}

	fields = append(fields, replacer.Field{
		Old: string(data),
		New: newContent,
	})

	return fields
}

// 解析镜像仓库host和name
func parseImageRepoAddr(addr string) (string, string) {
	splits := strings.Split(addr, "/")

	// 官方仓库地址
	if len(splits) == 1 {
		return "https://index.docker.io/v1", addr
	}

	// 非官方仓库地址
	l := len(splits)
	return strings.Join(splits[:l-1], "/"), splits[l-1]
}
