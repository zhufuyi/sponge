package replacer

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

var _ Replacer = (*replacerInfo)(nil)

// Replacer 接口
type Replacer interface {
	SetReplacementFields(fields []Field)
	SetIgnoreFiles(filenames ...string)
	SetIgnoreSubDirs(dirs ...string)
	SetSubDirs(subDirs ...string)
	SetOutputDir(absDir string, name ...string) error
	GetOutputDir() string
	GetSourcePath() string
	SaveFiles() error
	ReadFile(filename string) ([]byte, error)
}

// replacerInfo replacer信息
type replacerInfo struct {
	path              string   // 模板目录或文件
	fs                embed.FS // 模板目录对应二进制对象
	isActual          bool     // fs字段是否来源实际路径，如果为true，使用io操作文件，如果为false使用fs操作文件
	files             []string // 模板文件列表
	ignoreFiles       []string // 忽略替换的文件列表
	ignoreDirs        []string // 忽略处理的子目录
	replacementFields []Field  // 从模板文件转为新文件需要替换的字符
	outPath           string   // 输出替换后文件存放目录路径
}

// New 根据指定路径创建replacer
func New(path string) (Replacer, error) {
	files, err := gofile.ListFiles(path)
	if err != nil {
		return nil, err
	}

	path, _ = filepath.Abs(path)
	return &replacerInfo{
		path:              path,
		isActual:          true,
		files:             files,
		replacementFields: []Field{},
	}, nil
}

// NewFS 根据嵌入的路径创建replacer
func NewFS(path string, fs embed.FS) (Replacer, error) {
	files, err := listFiles(path, fs)
	if err != nil {
		return nil, err
	}

	return &replacerInfo{
		path:              path,
		fs:                fs,
		isActual:          false,
		files:             files,
		replacementFields: []Field{},
	}, nil
}

// Field 替换字段信息
type Field struct {
	Old             string // 模板字段
	New             string // 新字段
	IsCaseSensitive bool   // 第一个字母是否区分大小写
}

// SetReplacementFields 设置替换字段，注：old字符尽量不要存在包含关系，如果存在，在设置Field时注意先后顺序
func (r *replacerInfo) SetReplacementFields(fields []Field) {
	var newFields []Field
	for _, field := range fields {
		if field.IsCaseSensitive && isFirstAlphabet(field.Old) { // 拆分首字母大小写两个字段
			newFields = append(newFields,
				Field{ // 把第一个字母转为大写
					Old: strings.ToUpper(field.Old[:1]) + field.Old[1:],
					New: strings.ToUpper(field.New[:1]) + field.New[1:],
				},
				Field{ // 把第一个字母转为小写
					Old: strings.ToLower(field.Old[:1]) + field.Old[1:],
					New: strings.ToLower(field.New[:1]) + field.New[1:],
				},
			)
		} else {
			newFields = append(newFields, field)
		}
	}
	r.replacementFields = newFields
}

// SetSubDirs 设置处理指定子目录，其他目录下文件忽略处理
func (r *replacerInfo) SetSubDirs(subDirs ...string) {
	if len(subDirs) == 0 {
		return
	}

	subDirs = r.covertPathsDelimiter(subDirs...)

	var files []string
	isExistFile := make(map[string]struct{})
	for _, file := range r.files {
		for _, dir := range subDirs {
			if isSubPath(file, dir) {
				// 避免重复文件
				if _, ok := isExistFile[file]; ok {
					continue
				} else {
					isExistFile[file] = struct{}{}
				}
				files = append(files, file)
			}
		}
	}

	if len(files) == 0 {
		return
	}
	r.files = files
}

// SetIgnoreFiles 设置忽略处理的文件
func (r *replacerInfo) SetIgnoreFiles(filenames ...string) {
	r.ignoreFiles = append(r.ignoreFiles, filenames...)
}

// SetIgnoreSubDirs 设置忽略处理的子目录
func (r *replacerInfo) SetIgnoreSubDirs(dirs ...string) {
	dirs = r.covertPathsDelimiter(dirs...)
	r.ignoreDirs = append(r.ignoreDirs, dirs...)
}

// SetOutputDir 设置输出目录，优先使用absPath，如果absPath为空，自动在当前目录根据参数name名称生成输出目录
func (r *replacerInfo) SetOutputDir(absPath string, name ...string) error {
	// 输出到指定目录
	if absPath != "" {
		abs, err := filepath.Abs(absPath)
		if err != nil {
			return err
		}

		r.outPath = abs
		return nil
	}

	// 输出到当前目录
	subPath := strings.Join(name, "_")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	r.outPath = pwd + gofile.GetPathDelimiter() + subPath + "_" + time.Now().Format("150405")
	return nil
}

// GetOutputDir 获取输出目录
func (r *replacerInfo) GetOutputDir() string {
	return r.outPath
}

// GetSourcePath 获取源路径
func (r *replacerInfo) GetSourcePath() string {
	return r.path
}

// ReadFile 读取文件内容
func (r *replacerInfo) ReadFile(filename string) ([]byte, error) {
	filename = r.covertPathDelimiter(filename)

	foundFile := []string{}
	for _, file := range r.files {
		if strings.Contains(file, filename) && gofile.GetFilename(file) == gofile.GetFilename(filename) {
			foundFile = append(foundFile, file)
		}
	}
	if len(foundFile) != 1 {
		return nil, fmt.Errorf("total %d file named '%s', files=%+v", len(foundFile), filename, foundFile)
	}

	if r.isActual {
		return os.ReadFile(foundFile[0])
	}
	return r.fs.ReadFile(foundFile[0])
}

// SaveFiles 导出文件
func (r *replacerInfo) SaveFiles() error {
	if r.outPath == "" {
		r.outPath = gofile.GetRunPath() + gofile.GetPathDelimiter() + "generate_" + time.Now().Format("150405")
	}

	var existFiles []string
	var writeData = make(map[string][]byte)

	for _, file := range r.files {
		if r.isInIgnoreDir(file) || r.isIgnoreFile(file) {
			continue
		}

		// 从二进制读取模板文件内容使用embed.FS，如果要从指定目录读取使用os.ReadFile
		var data []byte
		var err error
		if r.isActual {
			data, err = os.ReadFile(file)
		} else {
			data, err = r.fs.ReadFile(file)
		}
		if err != nil {
			return err
		}

		// 替换文本内容
		for _, field := range r.replacementFields {
			data = bytes.ReplaceAll(data, []byte(field.Old), []byte(field.New))
		}

		// 获取新文件路径
		newFilePath := r.getNewFilePath(file)
		dir, filename := filepath.Split(newFilePath)
		// 替换文件名和文件夹名
		for _, field := range r.replacementFields {
			if strings.Contains(dir, field.Old) {
				dir = strings.ReplaceAll(dir, field.Old, field.New)
			}
			if strings.Contains(filename, field.Old) {
				filename = strings.ReplaceAll(filename, field.Old, field.New)
			}

			if newFilePath != dir+filename {
				newFilePath = dir + filename
			}
		}

		if gofile.IsExists(newFilePath) {
			existFiles = append(existFiles, newFilePath)
		}
		writeData[newFilePath] = data
	}

	if len(existFiles) > 0 {
		return fmt.Errorf("existing files detected\n    %s\nCode generation has been cancelled\n", strings.Join(existFiles, "\n    "))
	}

	for file, data := range writeData {
		// 保存文件
		err := saveToNewFile(file, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *replacerInfo) isIgnoreFile(file string) bool {
	isIgnore := false
	_, filename := filepath.Split(file)
	for _, v := range r.ignoreFiles {
		if filename == v {
			isIgnore = true
			break
		}
	}
	return isIgnore
}

func (r *replacerInfo) isInIgnoreDir(file string) bool {
	isIgnore := false
	dir, _ := filepath.Split(file)
	for _, v := range r.ignoreDirs {
		if strings.Contains(dir, v) {
			isIgnore = true
			break
		}
	}
	return isIgnore
}

func (r *replacerInfo) getNewFilePath(file string) string {
	var newFilePath string
	if r.isActual {
		newFilePath = r.outPath + strings.Replace(file, r.path, "", 1)
	} else {
		newFilePath = r.outPath + strings.Replace(file, r.path, "", 1)
	}

	if gofile.IsWindows() {
		newFilePath = strings.ReplaceAll(newFilePath, "/", "\\")
	}

	return newFilePath
}

// 如果是windows，转换路径分割符
func (r *replacerInfo) covertPathDelimiter(filePath string) string {
	if r.isActual && gofile.IsWindows() {
		filePath = strings.ReplaceAll(filePath, "/", "\\")
	}
	return filePath
}

// 如果是windows，批量转换路径分割符
func (r *replacerInfo) covertPathsDelimiter(filePaths ...string) []string {
	if r.isActual && gofile.IsWindows() {
		filePathsTmp := []string{}
		for _, dir := range filePaths {
			filePathsTmp = append(filePathsTmp, strings.ReplaceAll(dir, "/", "\\"))
		}
		return filePathsTmp
	}
	return filePaths
}

func saveToNewFile(filePath string, data []byte) error {
	// 创建目录
	dir, _ := filepath.Split(filePath)
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		return err
	}

	// 保存文件
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// 遍历嵌入的目录下所有文件，返回文件的绝对路径
func listFiles(path string, fs embed.FS) ([]string, error) {
	files := []string{}
	err := walkDir(path, &files, fs)
	return files, err
}

// 通过迭代方式遍历嵌入的目录
func walkDir(dirPath string, allFiles *[]string, fs embed.FS) error {
	files, err := fs.ReadDir(dirPath) // 读取目录下文件
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + "/" + file.Name()
		if file.IsDir() {
			_ = walkDir(deepFile, allFiles, fs)
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

// 判断字符串第一个字符是字母
func isFirstAlphabet(str string) bool {
	if len(str) == 0 {
		return false
	}

	if (str[0] >= 'A' && str[0] <= 'Z') || (str[0] >= 'a' && str[0] <= 'z') {
		return true
	}

	return false
}

func isSubPath(filePath string, subPath string) bool {
	dir, _ := filepath.Split(filePath)
	return strings.Contains(dir, subPath)
}
