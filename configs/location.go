package configs

import (
	"path/filepath"
	"runtime"
)

// used to locate a directory

var basePath string

// nolint
func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

// Path return absolute path
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basePath, rel)
}
