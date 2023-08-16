// Package replacer is a library of replacement file content, supports replacement of
// files in local directories and embedded directory files via embed.
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

// Replacer interface
type Replacer interface {
	SetReplacementFields(fields []Field)
	SetSubDirsAndFiles(subDirs []string, subFiles ...string)
	SetIgnoreSubDirs(dirs ...string)
	SetIgnoreSubFiles(filenames ...string)
	SetOutputDir(absDir string, name ...string) error
	GetOutputDir() string
	GetSourcePath() string
	SaveFiles() error
	ReadFile(filename string) ([]byte, error)
}

// replacerInfo replacer information
type replacerInfo struct {
	path              string   // template directory or file
	fs                embed.FS // Template directory corresponding to binary objects
	isActual          bool     // true: use os to manipulate files, false: use fs to manipulate files
	files             []string // list of template files
	ignoreFiles       []string // ignore the list of replaced files, e.g. ignore.txt or myDir/ignore.txt
	ignoreDirs        []string // ignore processed subdirectories
	replacementFields []Field  // characters to be replaced when converting from a template file to a new file
	outPath           string   // the directory where the file is saved after replacement
}

// New create replacer with local directory
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

// NewFS create replacer with embed.FS
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

// Field replace field information
type Field struct {
	Old             string // old field
	New             string // new field
	IsCaseSensitive bool   // whether the first letter is case-sensitive
}

// SetReplacementFields set the replacement field, note: old characters should not be included in the relationship,
// if they exist, pay attention to the order of precedence when setting the Field
func (r *replacerInfo) SetReplacementFields(fields []Field) {
	var newFields []Field
	for _, field := range fields {
		if field.IsCaseSensitive && isFirstAlphabet(field.Old) { // splitting the initial case field
			newFields = append(newFields,
				Field{ // convert the first letter to upper case
					Old: strings.ToUpper(field.Old[:1]) + field.Old[1:],
					New: strings.ToUpper(field.New[:1]) + field.New[1:],
				},
				Field{ // convert the first letter to lower case
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

// SetSubDirsAndFiles set up processing of specified subdirectories, files in other directories are ignored
func (r *replacerInfo) SetSubDirsAndFiles(subDirs []string, subFiles ...string) {
	if len(subDirs) == 0 {
		return
	}

	subDirs = r.convertPathsDelimiter(subDirs...)
	subFiles = r.convertPathsDelimiter(subFiles...)

	var files []string
	isExistFile := make(map[string]struct{}) // use map to avoid duplicate files
	for _, file := range r.files {
		for _, dir := range subDirs {
			if isSubPath(file, dir) {
				if _, ok := isExistFile[file]; ok {
					continue
				}
				isExistFile[file] = struct{}{}
				files = append(files, file)
			}
		}
		for _, sf := range subFiles {
			if isMatchFile(file, sf) {
				if _, ok := isExistFile[file]; ok {
					continue
				}
				isExistFile[file] = struct{}{}
				files = append(files, file)
			}
		}
	}

	if len(files) == 0 {
		return
	}
	r.files = files
}

// SetIgnoreSubFiles specify files to be ignored
func (r *replacerInfo) SetIgnoreSubFiles(filenames ...string) {
	r.ignoreFiles = append(r.ignoreFiles, filenames...)
}

// SetIgnoreSubDirs specify subdirectories to be ignored
func (r *replacerInfo) SetIgnoreSubDirs(dirs ...string) {
	dirs = r.convertPathsDelimiter(dirs...)
	r.ignoreDirs = append(r.ignoreDirs, dirs...)
}

// SetOutputDir specify the output directory, preferably using absPath, if absPath is empty,
// the output directory is automatically generated in the current directory according to the name of the parameter
func (r *replacerInfo) SetOutputDir(absPath string, name ...string) error {
	// output to the specified directory
	if absPath != "" {
		abs, err := filepath.Abs(absPath)
		if err != nil {
			return err
		}

		r.outPath = abs
		return nil
	}

	// output to the current directory
	subPath := strings.Join(name, "_")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	r.outPath = pwd + gofile.GetPathDelimiter() + subPath + "_" + time.Now().Format("150405")
	return nil
}

// GetOutputDir get output directory
func (r *replacerInfo) GetOutputDir() string {
	return r.outPath
}

// GetSourcePath get source directory
func (r *replacerInfo) GetSourcePath() string {
	return r.path
}

// ReadFile read file content
func (r *replacerInfo) ReadFile(filename string) ([]byte, error) {
	filename = r.convertPathDelimiter(filename)

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

// SaveFiles save file with setting
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

		var data []byte
		var err error

		if r.isActual {
			data, err = os.ReadFile(file) // read from local files
		} else {
			data, err = r.fs.ReadFile(file) // read from local embed.FS
		}
		if err != nil {
			return err
		}

		// replace text content
		for _, field := range r.replacementFields {
			data = bytes.ReplaceAll(data, []byte(field.Old), []byte(field.New))
		}

		// get new file path
		newFilePath := r.getNewFilePath(file)
		dir, filename := filepath.Split(newFilePath)
		// replace file names and directory names
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
		//nolint
		return fmt.Errorf("existing files detected\n    %s\nCode generation has been cancelled\n",
			strings.Join(existFiles, "\n    "))
	}

	for file, data := range writeData {
		if isForbiddenFile(file, r.path) {
			return fmt.Errorf("disable writing file(%s) to directory(%s), file size=%d", file, r.path, len(data))
		}
	}

	for file, data := range writeData {
		err := saveToNewFile(file, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *replacerInfo) isIgnoreFile(file string) bool {
	isIgnore := false
	for _, v := range r.ignoreFiles {
		if isMatchFile(file, v) {
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

func isForbiddenFile(file string, path string) bool {
	if gofile.IsWindows() {
		path = strings.ReplaceAll(path, "/", "\\")
		file = strings.ReplaceAll(file, "/", "\\")
	}
	return strings.Contains(file, path)
}

func (r *replacerInfo) getNewFilePath(file string) string {
	//var newFilePath string
	//if r.isActual {
	//	newFilePath = r.outPath + strings.Replace(file, r.path, "", 1)
	//} else {
	//	newFilePath = r.outPath + strings.Replace(file, r.path, "", 1)
	//}
	newFilePath := r.outPath + strings.Replace(file, r.path, "", 1)

	if gofile.IsWindows() {
		newFilePath = strings.ReplaceAll(newFilePath, "/", "\\")
	}

	return newFilePath
}

// if windows, convert the path splitter
func (r *replacerInfo) convertPathDelimiter(filePath string) string {
	if r.isActual && gofile.IsWindows() {
		filePath = strings.ReplaceAll(filePath, "/", "\\")
	}
	return filePath
}

// if windows, batch convert path splitters
func (r *replacerInfo) convertPathsDelimiter(filePaths ...string) []string {
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
	// create directory
	dir, _ := filepath.Split(filePath)
	err := os.MkdirAll(dir, 0766)
	if err != nil {
		return err
	}

	// save file
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// iterates over all files in the embedded directory, returning the absolute path to the file
func listFiles(path string, fs embed.FS) ([]string, error) {
	files := []string{}
	err := walkDir(path, &files, fs)
	return files, err
}

// iterating through the embedded catalog
func walkDir(dirPath string, allFiles *[]string, fs embed.FS) error {
	files, err := fs.ReadDir(dirPath)
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

// determine if the first character of a string is a letter
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

func isMatchFile(filePath string, sf string) bool {
	dir1, file1 := filepath.Split(filePath)
	dir2, file2 := filepath.Split(sf)
	if file1 != file2 {
		return false
	}

	if gofile.IsWindows() {
		dir1 = strings.ReplaceAll(dir1, "/", "\\")
		dir2 = strings.ReplaceAll(dir2, "/", "\\")
	} else {
		dir1 = strings.ReplaceAll(dir1, "\\", "/")
		dir2 = strings.ReplaceAll(dir2, "\\", "/")
	}

	l1, l2 := len(dir1), len(dir2)
	if l1 >= l2 && dir1[l1-l2:] == dir2 {
		return true
	}
	return false
}
