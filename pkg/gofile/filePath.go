package gofile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// IsExists 判断文件或文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetRunPath 获取程序执行的绝对路径
func GetRunPath() string {
	dir, err := os.Executable()
	if err != nil {
		fmt.Println("os.Executable error.", err.Error())
		return ""
	}

	return filepath.Dir(dir)
}

// GetFilename 获取文件名
func GetFilename(filePath string) string {
	_, name := filepath.Split(filePath)
	return name
}

// GetPathDelimiter 根据系统类型获取分隔符
func GetPathDelimiter() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}

	return delimiter
}

// ListFiles 遍历指定目录下所有文件，返回文件的绝对路径
func ListFiles(dirPath string, opts ...Option) ([]string, error) {
	files := []string{}
	err := error(nil)

	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		return files, err
	}

	o := defaultOptions()
	o.apply(opts...)

	switch o.filter {
	case prefix:
		return files, walkDirWithFilter(dirPath, &files, matchPrefix(o.name))
	case suffix:
		return files, walkDirWithFilter(dirPath, &files, matchSuffix(o.name))
	case contain:
		return files, walkDirWithFilter(dirPath, &files, matchContain(o.name))
	}

	return files, walkDir(dirPath, &files)
}

// ListDirsAndFiles 遍历指定目录下所有子目录文件，返回文件的绝对路径
func ListDirsAndFiles(dirPath string) (map[string][]string, error) {
	df := make(map[string][]string, 2)

	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return df, err
	}

	dirs := []string{}
	files := []string{}
	err = walkDir2(dirPath, &dirs, &files)
	if err != nil {
		return df, err
	}

	df["dirs"] = dirs
	df["files"] = files

	return df, nil
}

// 带过滤条件通过迭代方式遍历文件
func walkDirWithFilter(dirPath string, allFiles *[]string, filter filterFn) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			err = walkDirWithFilter(deepFile, allFiles, filter)
			if err != nil {
				return err
			}
			continue
		}
		if filter(deepFile) {
			*allFiles = append(*allFiles, deepFile)
		}
	}

	return nil
}

func walkDir2(dirPath string, allDirs *[]string, allFiles *[]string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			*allDirs = append(*allDirs, deepFile)
			err = walkDir2(deepFile, allDirs, allFiles)
			if err != nil {
				return err
			}
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

type filterFn func(string) bool

// 后缀匹配
func matchSuffix(suffixName string) filterFn {
	return func(filename string) bool {
		if suffixName == "" {
			return false
		}

		size := len(filename) - len(suffixName)
		if size >= 0 && filename[size:] == suffixName { // 后缀
			return true
		}
		return false
	}
}

// 前缀匹配
func matchPrefix(prefixName string) filterFn {
	return func(filePath string) bool {
		if prefixName == "" {
			return false
		}
		filename := GetFilename(filePath)
		size := len(filename) - len(prefixName)
		if size >= 0 && filename[:len(prefixName)] == prefixName { // 前缀
			return true
		}
		return false
	}
}

// 包含字符串
func matchContain(containName string) filterFn {
	return func(filePath string) bool {
		if containName == "" {
			return false
		}
		filename := GetFilename(filePath)
		return strings.Contains(filename, containName)
	}
}

// 通过迭代方式遍历文件
func walkDir(dirPath string, allFiles *[]string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			err = walkDir(deepFile, allFiles)
			if err != nil {
				return err
			}
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}
