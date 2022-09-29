package configs

import (
	"path/filepath"
	"runtime"
)

// 用来定位目录

var basePath string

// nolint
func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

// Path 返回绝对路径
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basePath, rel)
}
