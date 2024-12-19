// Package merge is merge the generated code into the template file, you don't worry about it affecting
// the logic code you have already written, in case of accidents, you can find the
// pre-merge code in the directory /tmp/sponge_merge_backup_code
package merge

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/go-dev-frame/sponge/pkg/gobash"
	"github.com/go-dev-frame/sponge/pkg/gofile"
)

var (
	defaultFuzzyFilename = "*.go.gen*"
	defaultSplitLineMark = []byte(`// ---------- Do not delete or move this split line, this is the merge code marker ----------`)
)

type code struct {
	key   string
	value string
}

type mergeParam struct {
	dir           string                                   // specify the folder where the code should be merged
	fuzzyFilename string                                   // fuzzy matching file name
	splitLineMark []byte                                   // file Contents Partition Line Marker
	mark          string                                   // code mark strings or regular expressions
	isLineCode    bool                                     // true:handles line code, false:handles code blocks
	parseCode     func(data []byte, markStr string) []code // parsing code method
	dt            string                                   // character form of date and time
	backupDir     string                                   // backup Code Catalog
}

func newMergeParam(dir string, mark string, isLineCode bool, parseCode func(data []byte, markStr string) []code) *mergeParam {
	return &mergeParam{
		dir:           dir,
		fuzzyFilename: defaultFuzzyFilename,
		splitLineMark: defaultSplitLineMark,
		mark:          mark,
		isLineCode:    isLineCode,
		parseCode:     parseCode,
		dt:            time.Now().Format("20060102T150405"),
		backupDir:     os.TempDir() + gofile.GetPathDelimiter() + "sponge_merge_backup_code",
	}
}

// SetFuzzyFileName setting fuzzy matching file names, use * for fuzzy matching
func (m *mergeParam) SetFuzzyFileName(fuzzyFilename string) {
	m.fuzzyFilename = fuzzyFilename
}

// SetSplitLineMark setting the file split line marker
func (m *mergeParam) SetSplitLineMark(lineMark string) {
	m.splitLineMark = []byte(lineMark)
}

func (m *mergeParam) runMerge() {
	files := gofile.FuzzyMatchFiles(m.dir + "/" + m.fuzzyFilename)
	files = filterAndRemoveOldFiles(files)
	for _, file := range files {
		successFile, err := m.runMergeCode(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if successFile != "" {
			if gofile.IsWindows() {
				ss := strings.Split(successFile, "\\internal\\")
				if len(ss) == 2 {
					successFile = "internal\\" + ss[len(ss)-1]
				}
			} else {
				ss := strings.Split(successFile, "/internal/")
				if len(ss) == 2 {
					successFile = "internal/" + ss[len(ss)-1]
				}
			}
			fmt.Printf("merge code to \"%s\" successfully.\n", successFile)
		}
	}
}

func (m *mergeParam) runMergeCode(file string) (string, error) {
	if file == "" {
		return "", nil
	}

	oldFile := getOldFile(file)

	data1, err := os.ReadFile(oldFile)
	if err != nil {
		return "", err
	}

	data2, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	count1 := bytes.Count(data1, m.splitLineMark)
	count2 := bytes.Count(data2, m.splitLineMark)
	if count1 != count2 {
		return "", errors.New(color.RedString("merge code failed (%s --> %s), manually merge code"+
			" reference document https://github.com/go-dev-frame/sponge/tree/main/cmd/sponge/commands/merge",
			cutPathPrefix(file), getTargetFilename(file)))
	}

	var data []byte

	if count1 == 0 {
		data = m.mergeData(data1, data2)
	} else {
		data1ss := bytes.Split(data1, m.splitLineMark)
		data2ss := bytes.Split(data2, m.splitLineMark)

		for index, subData1 := range data1ss {
			subData := m.mergeData(subData1, data2ss[index])
			data = append(data, subData...)
			if index < len(data1ss)-1 {
				data = append(data, m.splitLineMark...)
			}
		}
	}

	if len(data1) > len(data) {
		return "", errors.New(color.RedString("merge code failed (%s --> %s), to avoid replacing logical code, "+
			"manually merge code reference document https://github.com/go-dev-frame/sponge/tree/main/cmd/sponge/commands/merge",
			cutPathPrefix(file), getTargetFilename(file)))
	}

	if len(data1) == len(data) {
		return "", os.Remove(file)
	}

	err = m.saveFile(oldFile, data)
	if err != nil {
		return "", err
	}

	return oldFile, os.Remove(file)
}

func (m *mergeParam) mergeData(subData1 []byte, subData2 []byte) []byte {
	c1 := m.parseCode(subData1, m.mark)
	c2 := m.parseCode(subData2, m.mark)

	var addCode, position []byte
	if m.isLineCode {
		addCode, position = compareCode(c1, c2)
	} else {
		addCode, position = compareCode2(c1, c2, subData2)
	}
	return mergeCode(subData1, addCode, position)
}

func (m *mergeParam) saveFile(file string, data []byte) error {
	bkDir := m.backupDir + gofile.GetPathDelimiter() + m.dt
	_ = os.MkdirAll(bkDir, 0744)
	_, _ = gobash.Exec("cp", file, bkDir+gofile.GetPathDelimiter()+gofile.GetFilename(file))

	return os.WriteFile(file, data, 0766)
}

func getOldFile(file string) string {
	dir, name := filepath.Split(file)
	return dir + strings.TrimSuffix(name, path.Ext(name))
}

func compareCode(oldCode []code, newCode []code) ([]byte, []byte) {
	var addCode []string
	var position string

	for _, code1 := range newCode {
		isEqual := false
		for _, code2 := range oldCode {
			if code1.key == code2.key {
				isEqual = true
				break
			}
		}
		if !isEqual {
			addCode = append(addCode, code1.value)
		}
	}

	l := len(oldCode)
	if l > 0 {
		position = oldCode[l-1].value // last position
	}

	addData := checkAndAdjustErrorCode(addCode, position, l)

	return addData, []byte(position)
}

func compareCode2(oldCode []code, newCode []code, data []byte) ([]byte, []byte) {
	var addCode []byte
	var position []byte

	for _, code1 := range newCode {
		isEqual := false
		for _, code2 := range oldCode {
			if code1.key == code2.key {
				isEqual = true
				break
			}
		}
		if !isEqual {
			_, name := getTmplKey(code1.value)
			comment := getComment(name, string(data))
			addCode = append(addCode, []byte("\n\n"+comment+"\n"+code1.value)...)
		}
	}

	if len(oldCode) > 0 {
		position = []byte(oldCode[len(oldCode)-1].value) // last position
	}

	return addCode, position
}

func mergeCode(oldCode []byte, addCode []byte, position []byte) []byte {
	if len(addCode) == 0 {
		return oldCode
	}

	var data []byte

	if len(position) == 0 {
		data = append(oldCode, addCode...)
		return data
	}

	ss := bytes.SplitN(oldCode, position, 2)
	if len(ss) != 2 {
		return oldCode
	}
	data = append(ss[0], position...)
	data = append(data, addCode...)
	data = append(data, ss[1]...)

	return data
}

// ------------- parsing the core of data in the internal/ecode directory -------------

func parseFromECode(data []byte, markStr string) []code {
	var codes []code
	buf := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, markStr) {
			name := getECodeMarkName(line, markStr)
			if name != "" {
				codes = append(codes, code{
					key:   name,
					value: line,
				})
			}
		}
	}

	return codes
}

func getECodeMarkName(str string, markStr string) string {
	ss := strings.SplitN(str, markStr, 2)
	name := strings.Replace(ss[0], " ", "", -1)
	name = strings.Replace(name, "=", "", -1)
	return strings.Replace(name, "	", " ", -1)
}

// ------------- parsing the core of data in the internal/routers directory -------------

func parseFromRouters(data []byte, markStr string) []code {
	var codes []code
	buf := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, markStr) {
			name := getRoutersMarkName(line, markStr)
			if name != "" {
				codes = append(codes, code{
					key:   name,
					value: line,
				})
			}
		}
	}

	return codes
}

func getRoutersMarkName(str string, markStr string) string {
	str = strings.Replace(str, " ", "", -1)
	ss := strings.SplitN(str, ",", 3)
	if len(ss) != 3 {
		return ""
	}

	sss := strings.Split(ss[0], markStr)
	if len(sss) != 2 {
		return ""
	}
	method := strings.Replace(sss[1], "(", "", -1)
	method = strings.Replace(method, "\"", "", -1)

	router := strings.Replace(ss[1], "\"", "", -1)

	return method + "-->" + router
}

// ------------- parsing the core of data in the internal/handler or internal/handler directory -------------

func parseFromTmplCode(data []byte, markStr string) []code {
	var codes []code
	str := string(data)
	reg1 := regexp.MustCompile(markStr)
	matches := reg1.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		for _, v := range match {
			key, _ := getTmplKey(v)
			codes = append(codes, code{
				key:   key,
				value: v,
			})
		}
	}
	return codes
}

func getTmplKey(str string) (key string, methodName string) {
	regStr2 := `func \((.*?)\) (.*?)\(`
	reg2 := regexp.MustCompile(regStr2)
	matches := reg2.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return "", ""
	}
	if len(matches[0]) != 3 {
		return "", ""
	}
	key = matches[0][0]
	methodName = matches[0][2]
	return key, methodName
}

func getComment(name string, str string) string {
	regStr := `//( ?)` + name + `[\w\W]*?\nfunc`
	reg := regexp.MustCompile(regStr)
	match := reg.FindAllString(str, -1)
	if len(match) == 0 {
		return ""
	}
	return strings.ReplaceAll(match[0], "\nfunc", "")
}

// ------------------------------------------------------------------------------------------

func adaptDir(dir string) string {
	if dir == "." || dir == "./" || dir == ".\\" {
		return ""
	}
	l := len(dir)
	if dir[l-1] == '/' {
		return dir
	}
	if dir[l-1] == '\\' {
		return dir[:l-1] + "/"
	}
	return dir + "/"
}

func cutPathPrefix(srcFile string) string {
	dirPath, _ := filepath.Abs(".")
	return strings.ReplaceAll(srcFile, dirPath+gofile.GetPathDelimiter(), "")
}

func getTargetFilename(file string) string {
	filename := gofile.GetFilename(file)
	ss := strings.Split(filename, ".go.gen")
	if len(ss) != 2 {
		return file
	}
	return ss[0] + ".go"
}

func filterAndRemoveOldFiles(files []string) []string {
	if len(files) < 2 {
		return files
	}

	var groupFiles = make(map[string][]string)
	for _, file := range files {
		filePrefix := strings.Split(file, ".go.gen")
		if len(filePrefix) != 2 {
			continue
		}
		if _, ok := groupFiles[filePrefix[0]]; !ok {
			groupFiles[filePrefix[0]] = []string{file}
		} else {
			groupFiles[filePrefix[0]] = append(groupFiles[filePrefix[0]], file)
		}
	}

	var newFiles, removeFiles []string
	for _, fs := range groupFiles {
		l := len(fs)
		if l == 1 {
			newFiles = append(newFiles, fs[0])
		} else if l > 1 {
			sort.Strings(fs)
			newFiles = append(newFiles, fs[l-1])
			removeFiles = append(removeFiles, fs[:l-1]...)
		}
	}

	// remove old files
	for _, file := range removeFiles {
		_ = os.Remove(file)
	}

	return newFiles
}

var (
	errCodeStrMark1 = "errcode.NewError("
	errCodeStrMark2 = "errcode.NewRPCStatus("
)

func checkAndGetErrorCodeStr(str string) string {
	if !strings.Contains(str, errCodeStrMark1) && !strings.Contains(str, errCodeStrMark2) {
		return ""
	}

	// match strings between left parentheses and commas using regular expressions
	// string format: ErrLoginUser = errcode.NewError(userBaseCode+2, "failed to Login "+userName)
	pattern := `\(([^)]+?),`
	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(str)
	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func parseErrorCode(str string) (string, int) {
	ss := strings.Split(str, "+")
	if len(ss) != 2 {
		return "", 0
	}
	num, _ := strconv.Atoi(strings.TrimSpace(ss[1]))
	return ss[0], num
}

func checkAndAdjustErrorCode(addCode []string, position string, l int) []byte {
	data := []byte(strings.Join(addCode, ""))

	str := checkAndGetErrorCodeStr(position)
	if str == "" {
		return data
	}
	referenceStr, maxNum := parseErrorCode(str)
	if referenceStr == "" || maxNum == 0 {
		return data
	}
	if maxNum < l {
		maxNum = l
	}

	// adjust error code
	var newCode []byte
	for _, line := range addCode {
		codeStr := checkAndGetErrorCodeStr(line)
		if codeStr == "" {
			return data
		}
		baseStr, num := parseErrorCode(codeStr)
		if baseStr == "" || num == 0 {
			return data
		}
		maxNum++
		newLine := strings.ReplaceAll(line, codeStr, fmt.Sprintf("%s+%d", referenceStr, maxNum))
		newCode = append(newCode, []byte(newLine)...)
	}

	return newCode
}

// ------------------------------------------------------------------------------------------

func mergeHTTPECode(dir string) {
	m := newMergeParam(
		dir+"internal/ecode",
		"errcode.NewError(",
		true,
		parseFromECode,
	)
	m.runMerge()
}

func mergeGRPCECode(dir string) {
	m := newMergeParam(
		dir+"internal/ecode",
		"errcode.NewRPCStatus(",
		true,
		parseFromECode,
	)
	m.runMerge()
}

func mergeGinRouters(dir string) {
	m := newMergeParam(
		dir+"internal/routers",
		"c.setSinglePath(",
		true,
		parseFromRouters,
	)
	m.runMerge()
}

func mergeHTTPHandlerTmpl(dir string) {
	m := newMergeParam(
		dir+"internal/handler",
		`func \(h[\w\W]*?\n}`,
		false,
		parseFromTmplCode,
	)
	m.runMerge()
}

func mergeGRPCServiceClientTmpl(dir string) {
	m := newMergeParam(
		dir+"internal/service",
		`func \(c[\w\W]*?\n}`,
		false,
		parseFromTmplCode,
	)
	m.runMerge()
}

func mergeGRPCServiceTmpl(dir string) {
	m := newMergeParam(
		dir+"internal/service",
		`func \(s[\w\W]*?\n}`,
		false,
		parseFromTmplCode,
	)
	m.runMerge()
}
