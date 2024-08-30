package patch

import (
	"os"
	"path/filepath"
	"strings"
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
