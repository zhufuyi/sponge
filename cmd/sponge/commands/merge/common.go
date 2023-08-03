package merge

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
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
	parseCode     func(date []byte, markStr string) []code // parsing code method
	dt            string                                   // character form of date and time
	backupDir     string                                   // backup Code Catalog
}

func newMergeParam(dir string, mark string, isLineCode bool, parseCode func(date []byte, markStr string) []code) *mergeParam {
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
	result, err := gobash.Exec("ls", m.dir+"/"+m.fuzzyFilename)
	if err != nil {
		//fmt.Println("Warring:", err)
		return
	}

	files := strings.Split(string(result), "\n")
	for _, file := range files {
		successFile, err := m.runMergeCode(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if successFile != "" {
			fmt.Printf("merge code to '%s' successfully.\n", successFile)
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
		return "", fmt.Errorf("merge code mark mismatch, please merge codes manually, file = %s", file)
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
		return "", fmt.Errorf("to avoid replacing logical code, please merge codes manually, file = %s", file)
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
	var addCode []byte
	var position []byte

	for _, code1 := range newCode {
		isEqual := false
		for _, code2 := range oldCode {
			if code1.key == code2.key {
				isEqual = true
				position = []byte(code2.value)
				break
			}
		}
		if !isEqual {
			addCode = append(addCode, []byte(code1.value)...)
		}
	}

	return addCode, position
}

func compareCode2(oldCode []code, newCode []code, data []byte) ([]byte, []byte) {
	var addCode []byte
	var position []byte

	for _, code1 := range newCode {
		isEqual := false
		for _, code2 := range oldCode {
			if code1.key == code2.key {
				isEqual = true
				position = []byte(code2.value)
				break
			}
		}
		if !isEqual {
			_, name := getTmplKey(code1.value)
			comment := getComment(name, string(data))
			addCode = append(addCode, []byte("\n\n"+comment+"\n"+code1.value)...)
		}
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

func parseFromECode(date []byte, markStr string) []code {
	var codes []code
	buf := bufio.NewReader(bytes.NewReader(date))
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

func parseFromRouters(date []byte, markStr string) []code {
	var codes []code
	buf := bufio.NewReader(bytes.NewReader(date))
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

func getTmplKey(str string) (string, string) {
	regStr2 := `func \((.*?)\) (.*?)\(`
	reg2 := regexp.MustCompile(regStr2)
	matches := reg2.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return "", ""
	}

	if len(matches[0]) != 3 {
		return "", ""
	}

	return matches[0][0], matches[0][2]
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

func mergeHTTPECode() {
	m := newMergeParam(
		"internal/ecode",
		"errcode.NewError(",
		true,
		parseFromECode,
	)
	m.runMerge()
}

func mergeGRPCECode() {
	m := newMergeParam(
		"internal/ecode",
		"errcode.NewRPCStatus(",
		true,
		parseFromECode,
	)
	m.runMerge()
}

func mergeGinRouters() {
	m := newMergeParam(
		"internal/routers",
		"c.setSinglePath(",
		true,
		parseFromRouters,
	)
	m.runMerge()
}

func mergeHTTPHandlerTmpl() {
	m := newMergeParam(
		"internal/handler",
		`func \(h[\w\W]*?\n}`,
		false,
		parseFromTmplCode,
	)
	m.runMerge()
}

func mergeGRPCServiceClientTmpl() {
	m := newMergeParam(
		"internal/service",
		`func \(c[\w\W]*?\n}`,
		false,
		parseFromTmplCode,
	)
	m.runMerge()
}

func mergeGRPCServiceTmpl() {
	m := newMergeParam(
		"internal/service",
		`func \(s[\w\W]*?\n}`,
		false,
		parseFromTmplCode,
	)
	m.runMerge()
}
