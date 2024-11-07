package patch

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

// get moduleName and serverName from directory
func getNamesFromOutDir(dir string) (moduleName string, serverName string, suitedMonoRepo bool) {
	if dir == "" {
		return "", "", false
	}
	data, err := os.ReadFile(dir + "/docs/gen.info")
	if err != nil {
		return "", "", false
	}

	ms := strings.Split(string(data), ",")
	if len(ms) == 2 {
		return ms[0], ms[1], false
	} else if len(ms) >= 3 {
		return ms[0], ms[1], ms[2] == "true"
	}

	return "", "", false
}

func cutPath(srcProtoFile string) string {
	dirPath, _ := filepath.Abs("..")
	srcProtoFile = strings.ReplaceAll(srcProtoFile, dirPath, "..")
	return strings.ReplaceAll(srcProtoFile, "\\", "/")
}

func cutPathPrefix(srcProtoFile string) string {
	dirPath, _ := filepath.Abs(".")
	srcProtoFile = strings.ReplaceAll(srcProtoFile, dirPath, ".")
	return strings.ReplaceAll(srcProtoFile, "\\", "/")
}

func listErrCodeFiles(dir string) ([]string, error) {
	files, err := gofile.ListFiles(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("not found files")
	}

	filterFiles := []string{}
	for _, file := range files {
		if strings.Contains(file, "systemCode_http.go") || strings.Contains(file, "systemCode_rpc.go") {
			continue
		}
		if strings.Contains(file, "_http.go") || strings.Contains(file, "_rpc.go") {
			filterFiles = append(filterFiles, file)
		}
	}

	return filterFiles, nil
}

func getSubFiles(selectedFiles map[string][]string) []string {
	subFiles := []string{}
	for dir, files := range selectedFiles {
		for _, file := range files {
			subFiles = append(subFiles, dir+"/"+file)
		}
	}
	return subFiles
}
